package compose

import (
  "fmt"
  "strings"
  "reflect"
  "regexp"
  "errors"
  slspp "github.com/aboltart/go-support/slice"
  strspp "github.com/aboltart/go-support/string"
  sexec "github.com/docker-support/exec"
  "github.com/docker-support/config"

  "os"

  "github.com/fgrehm/go-dockerpty"
  "github.com/fsouza/go-dockerclient"
)

type CmdParams struct {
  Command string
  Silent  bool
}

var cmd ComposeCommand
type ComposeCommand struct {}

func (cmd *ComposeCommand) Create() { _, _, _ = cmd.standartComposeCommandRun("create", false) }
func (cmd *ComposeCommand) Config() { _, _, _ = composeRun( CmdParams{ Command: "config" } ) }
func (cmd *ComposeCommand) Logs()   { _, _, _ = cmd.standartComposeCommandRun("logs", false) }
func (cmd *ComposeCommand) Ps()     { cmd.ps(false) }
func (cmd *ComposeCommand) Start()  { _, _, _ = cmd.standartComposeCommandRun("start", false) }
func (cmd *ComposeCommand) Stop()   { _, _, _ = cmd.standartComposeCommandRun("stop", false) }

func (cmd *ComposeCommand) Down() {
  fullCmd := "down " + strings.Join(Args(), " ")
  _, _, _ = composeRun( CmdParams{ Command: fullCmd, Silent: false } )
}

func (cmd *ComposeCommand) Help() {
  fullCmd := "help " + strings.Join(Args(), " ")
  _, _, _ = composeRun( CmdParams{ Command: fullCmd, Silent: false } )
}

// Should move out to seperate repo. Probally as docker-support-rails to extend current project
// func (cmd *ComposeCommand) Bundle() {
//   PrependArg("bundle")
//   cmd.host()
// }

// func (cmd *ComposeCommand) Rails() {
//   PrependArg("rails")
//   cmd.host()
// }

func (cmd *ComposeCommand) Host() {
  cmd.host()
}

func (cmd *ComposeCommand) Remove() {
  fullCmd := "rm -f " + strings.Join(Args(), " ")
  _, _, _ = composeRun( CmdParams{ Command: fullCmd, Silent: false } )
}

func (cmd *ComposeCommand) Up() {
  fullCmd := "up -d " + strings.Join(Args(), " ")
  _, _, _ = composeRun( CmdParams{ Command: fullCmd, Silent: false } )
}

func (cmd *ComposeCommand) Services() {
  // If was not passed docker service name, then printout Service Definitions
  if Service() == "" {
    ServiceDefinitionNames()
  } else {
    _ = cmd.serviceNames(false)
  }
}

func (cmd *ComposeCommand) ContainerNames() {
  for _, name := range cmd.containerNames() {
    fmt.Println(name)
  }
}

func (cmd *ComposeCommand) Stats() {
  containerNames := cmd.containerNames()

  fullCmd := "docker stats " + strings.Join(containerNames, " ")
  _, _, _  = sexec.Exec( sexec.ExecParams{ Command: fullCmd } )
}

func (cmd *ComposeCommand) Build() {

  fullCmd := "build "
  args    := Args()

  // Should construct
  // docker-compose build --no-cache service1 service2
  if slspp.StringInSlice("--no-cache", args) {
    fullCmd = fullCmd + "--no-cache "
    args    = slspp.RemoveFromSlice("--no-cache", args)
  }


  fullCmd = fullCmd + SubService() + " " + strings.Join(args, " ")
  _, _, _ = composeRun( CmdParams{ Command: fullCmd } )
}

// Private (Struct) section
func (cmd *ComposeCommand) standartComposeCommandRun(command string, silent bool) (string, string, error) {
  fullCmd := command + " " + SubService() + " " + strings.Join(Args(), " ")
  stdout, stderr, err := composeRun( CmdParams{ Command: fullCmd, Silent: silent } )

  return stdout, stderr, err
}

func (cmd *ComposeCommand) ps(silent bool) string {
  stdout, _, _ := cmd.standartComposeCommandRun("ps", silent)
  return stdout
}

func (cmd *ComposeCommand) containerNames() []string {

  stdout  := cmd.ps(true)
  results := slspp.CompactStringSlice(strings.Split(stdout, "\n"))

  // Compile the delimiter as a regular expression.
  containerNames  := []string{}
  re              := regexp.MustCompile(`\s*Name\s*Command\s*State\s*Ports\s*?`)
  foundHeaderLine := false
  headerLineNr    := 0

  for i, line := range results {
    if !foundHeaderLine {
      // Find ps result header row number
      if re.MatchString(line) {
        foundHeaderLine = true
        headerLineNr    = i + 2
      }
    } else {
      if i >= headerLineNr {
        // Execute as extra bash command to get first column as result
        fullCommand  := "result=$(echo \"$(echo \"" + line + "\" | awk '{print $1}')\"); echo $result"
        name, _, _ := sexec.Exec( sexec.ExecParams{ Command: fullCommand, Silent: true } )

        containerNames = append(containerNames, strings.TrimSpace(name))
      }
    }
  }

  return containerNames
}

func (cmd *ComposeCommand) serviceNames(silent bool) []string {
  names := []string{}

  stdout, _, _ := composeRun( CmdParams{ Command: "config --services", Silent: silent } )

  // TODO: Some error checking
  names = slspp.CompactStringSlice(strings.Split(stdout, "\n"))

  return names
}

// END ComposeCommand struct

func AvailableComposeCommands() []string {
  availableCommands := []string{}
  cmdType           := reflect.TypeOf(&cmd)

  for i := 0; i < cmdType.NumMethod(); i++ {
    method := cmdType.Method(i)

    // Select only exposed (Public methodes)
    if strspp.IsUpper(method.Name) {
      //Transform to lover for easily mapping
      availableCommands = append(availableCommands, strspp.CamelToSnake(method.Name))
    }
  }

  return availableCommands
}

func ComposePerform() {
  if slspp.StringInSlice(Command(), AvailableComposeCommands()) {
    // Convert command to Capetalized, to be as exposable
    cmdForExec := strspp.SnakeToCamel(Command())

    // pointer to struct - addressable
    reflect.ValueOf(&cmd).MethodByName(cmdForExec).Call([]reflect.Value{})
  } else {
    fmt.Println("Unknown Docker Compose command: ", Command())
  }

}

// Private section
func composeRun(params CmdParams) (string, string, error) {
  composeFiles := ComposeFiles()
  filesArg     := ""

  config.LoadEnvShellVariables()

  for _, file := range composeFiles {
    filesArg = filesArg + " -f " + file
  }

  fullCommand := "docker-compose " + filesArg + " " + params.Command
  stdout, stderr, err := sexec.Exec( sexec.ExecParams{ Command: fullCommand, Silent: params.Silent } )

  return stdout, stderr, err
}

func (cmd *ComposeCommand) host() {
  args          := Args()
  fullCmd       := "run "
  containerName := ""
  tty           := false

  // If there is performed command where is need TTY
  // Start container in deatch mode for attaching leater
  if len(args) == 0 || len(args) > 0 && slspp.StringInSlice(args[0], []string{"/bin/bash", "bin/bash", "bash"}) {
    fullCmd += "-d " + SubService() + " "
    fullCmd += " tail -f /dev/null "
    tty = true
  // Else run command on exit rm container
  } else {
    fullCmd += "--rm " + SubService() + " "
  }

  fullCmd += strings.Join(Args(), " ")

  stdout, _, err := composeRun( CmdParams{ Command: fullCmd, Silent: false } )
  if err != nil {
  } else {
    // Last line should be started container name
    data := slspp.CompactStringSlice(strings.Split(stdout, "\n"))
    containerName = data[len(data)-1]
  }

  // If TTY then attach to it
  if tty {
    _ = cmd.attach(containerName, "/bin/bash")
  }
}

func (cmd *ComposeCommand) attach(containerName string, command string) (error) {
  if containerName == "" {
    return errors.New("No Container")
  }

  if command == "" {
    command = "/bin/bash"
  }


  containerID, _, _ := sexec.Exec( sexec.ExecParams{ Command: "docker ps -aqf \"name=" + containerName + "\"", Silent: true } )
  client, _         := docker.NewClientFromEnv()


  exec, err := client.CreateExec(docker.CreateExecOptions{
    Container:    containerName,
    AttachStdin:  true,
    AttachStdout: true,
    AttachStderr: true,
    Tty:          true,
    Cmd:          []string{command},
  })

if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // Remove container
  defer func() {
    _, _, _ = sexec.Exec( sexec.ExecParams{ Command: "docker rm -f " + containerID } )
  }()

  // Fire up the console
  if err = dockerpty.StartExec(client, exec); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  return nil
}
