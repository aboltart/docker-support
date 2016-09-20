package exec

import (
  "fmt"
  "os/exec"
  "bufio"
  "github.com/docker-support/echo"
)

type ExecParams struct {
  Command string
  Silent  bool
}

// Some code from https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
func Exec(params ExecParams) (string, string, error) {

  cmd := exec.Command("sh", "-c", params.Command)

  cmdOutReader, err := cmd.StdoutPipe()
  if err != nil {
    fmt.Printf("Error creating StdOutPipe for Cmd", err)
  }

  cmdErrReader, err := cmd.StderrPipe()
  if err != nil {
    fmt.Printf("Error creating StdErrPipe for Cmd", err)
  }

  stdOutText     := ""
  fullStdOutText := ""
  outScanner     := bufio.NewScanner(cmdOutReader)
  go func() {
    for outScanner.Scan() {
      stdOutText = outScanner.Text()

      if !params.Silent { fmt.Println(stdOutText) }
      fullStdOutText = fullStdOutText + stdOutText + "\n"
    }
  }()

  stdErrText     := ""
  fullStdErrText := ""
  errScanner     := bufio.NewScanner(cmdErrReader)
  go func() {
    for errScanner.Scan() {
      stdErrText = errScanner.Text()

      if !params.Silent { fmt.Println(stdErrText) }
      fullStdErrText = fullStdErrText + stdErrText + "\n"
    }
  }()

  //Print runing command
  if !params.Silent { echo.Yellow(params.Command) }

  cmd_err := cmd.Run()

  return fullStdOutText, fullStdErrText, cmd_err
}
