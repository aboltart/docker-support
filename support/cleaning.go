package support

import (
  "fmt"
  "strconv"
  "github.com/docker-support/exec"
)


func RemoveUntagedImages() {
  stdout, stderr, err := exec.Exec( exec.ExecParams{ Command: "docker images | grep \"^<none>\" | wc -l | sed -e 's/^[ \t]*//'" } )

  if err != nil {
    // fmt.Println("Error: ",err)
    fmt.Println(stderr)
  } else {
    count, _ := strconv.Atoi(stdout)

    if count > 0 {

      stdout, stderr, err := exec.Exec( exec.ExecParams{ Command: "docker rmi $(docker images | grep -w '<none>' | awk '{print $3}')" } )

      if err != nil {
        // fmt.Println("Error: ",err)
        fmt.Println(stderr)
      } else {
        fmt.Println(stdout) // out is the standart output. It's in the format []byte
      }
    } else {
      fmt.Println("No untaged images for cleaning")
    }
  }
}

// https://github.com/chadoe/docker-cleanup-volumes
func CleaningOrphanedVolumes() {
  stdout, stderr, err := exec.Exec( exec.ExecParams{ Command:"docker volume ls -qf dangling=true | wc -l | sed -e 's/^[ \t]*//'" } )

  if err != nil {
    // fmt.Println("Error: ",err)
    fmt.Println(stderr)
  } else {
    count, _ := strconv.Atoi(stdout)

    if count > 0 {
      stdout, stderr, err := exec.Exec( exec.ExecParams{ Command:"docker volume rm $(docker volume ls -qf dangling=true)" } )

      if err != nil {
        // fmt.Println("Error: ",err)
        fmt.Println(stderr)
      } else {
        fmt.Println(stdout) // out is the standart output. It's in the format []byte
      }
    } else {
      fmt.Println("No orphaned volumes for cleaning")
    }
  }

}
