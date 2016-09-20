package compose

import (
  // "fmt"
  slspp "github.com/aboltart/go-support/slice"
  // "github.com/docker-support/config"
  "github.com/docker-support/argument"
  // "reflect"
)

var a *Arg

type Arg struct {
  argument.Arg //Extend struct from Common Arguments

  // Compose attributes
  service     string
  subservice  string
  command     string

  parsed      bool
}

func init() {
  a = new(Arg)
  a.parsed = false
}

func Service() string { return a.Service() }
func (a *Arg) Service() string {
  a.parseArgs()
  return a.service
}

func SubService() string { return a.SubService() }
func (a *Arg) SubService() string {
  a.parseArgs()
  return a.subservice
}

func Command() string { return a.Command() }
func (a *Arg) Command() string {
  a.parseArgs()
  return a.command
}

func Args() []string { return a.Args() }
func (a *Arg) Args() []string {
  a.parseArgs()
  return a.Arg.Args
}

func PrependArg(arg string) { a.PrependArgs(arg) }
func (a *Arg) PrependArgs(arg string) {
  a.parseArgs()
  a.Arg.Args = slspp.PrependForSlice(arg, a.Arg.Args)
}

func parseArgs() { a.parseArgs()}
func (a *Arg) parseArgs() {
  if ! a.parsed {
    a.parsed = true

    definedServices := ServiceDefinitions()
    args            := a.Arguments()

    // Find from Arguments Service Name
    if len(args) > 0 && slspp.StringInSlice(args[0], definedServices) {
      a.service  = args[0]
      a.Arg.Args = args[1:]
    }

    subServices := cmd.serviceNames(true)
    args        = a.Arguments()

    // If service present AND subservices found AND arguments present
    // Find subservice name
    if a.service != "" && len(subServices) > 0 && len(args) > 0  {
      firstArgument := args[0]

      if slspp.StringInSlice(firstArgument, subServices) {
        a.subservice = firstArgument
        a.Arg.Args   = args[1:]
      }
    }

    commands := AvailableComposeCommands()

    args = a.Arguments()

    // If service present AND commands defined AND arguments present
    // Find command for execute
    if a.service != "" && len(commands) > 0 && len(args) > 0 {
      firstArgument := args[0]

      if slspp.StringInSlice(firstArgument, commands) {
        a.command = firstArgument
        a.Arg.Args   = args[1:]
      }
    }

    if a.service == "" && len(args) > 0 {
      firstArgument := args[0]

      if firstArgument == "services" {
        a.command  = firstArgument
        a.Arg.Args = args[1:]
      }
    }


  }
}
