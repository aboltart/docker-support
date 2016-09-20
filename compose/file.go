package compose

import (
  "fmt"
  osspp "github.com/aboltart/go-support/os"
  "github.com/docker-support/config"
)


func ComposeFiles() []string {
  serviceName           := Service()
  serviceDefinitionPath := ServiceDefinitionPath()
  composeFiles          := []string{}
  multiComposeServices  := config.MultipleComposeServices()

  if serviceName == "" {
    fmt.Println("No service.")
    return composeFiles
  }

  if serviceDefinitionPath == "" {
    fmt.Println("No service path.")
    return composeFiles
  }

  // #################################################################
  // # Build service from multiple compose files
  // #   to fix services with 'volumes_from' cannot be extended case
  // #################################################################
  composeFiles = append(composeFiles, composeFileForService(serviceName, serviceDefinitionPath))

  subServices := multiComposeServices[serviceName]

  if subServices != nil {
    for _, subService := range subServices.([]interface{}) {
      composeFiles = append(composeFiles, composeFileForService(subService.(string), serviceDefinitionPath))
    }
  }

  return composeFiles
}

// Private Function section
func composeFileForService(serviceName string, serviceDefinitionPath string) string {
  environment := config.Environment()
  serviceDir  := serviceDefinitionPath + "/" + serviceName

   //Check if exists specific compose file for environment
  if environment != "" && osspp.IsFile(serviceDir + "/compose." + environment + ".yml") {
    return serviceDir + "/compose." + environment + ".yml"
  } else if osspp.IsFile(serviceDir + "/compose.yml") {
    return serviceDir + "/compose.yml"
  } else {
    fmt.Println("Not found compose file for service: ", serviceName, " in directory ", serviceDir)
    return ""
  }
}

