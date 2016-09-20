package argument

import (
  "os"
)

var a *Arg

// Define Arguments instance
type Arg struct {
  Args        []string
  parsed      bool
}

func init() {
  a = new(Arg)
  a.parsed = false
}

func Arguments() []string { return a.Arguments() }
func (a *Arg) Arguments() []string {
  a.parseArgs()
  return a.Args
}

// Private function section
func (a *Arg) parseArgs() {

  if ! a.parsed {
    a.Args = os.Args[1:] // args Without Programm
    a.parsed = true
  }
}
