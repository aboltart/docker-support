package echo

import (
  "fmt"
)

const (
  escape = "\x1b"

  yellow = "93m"
  red    = "31m"
  green  = "32m"
)

func Green(text string) {
  fmt.Println(escape + "[" + green + text + escape + "[0m")
}

func Yellow(text string) {
  fmt.Println(escape + "[" + yellow + text + escape + "[0m")
}

func Red(text string) {
  fmt.Println(escape + "[" + red + text + escape + "[0m")
}
