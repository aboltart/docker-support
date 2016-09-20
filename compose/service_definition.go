package compose

import (
  "fmt"
  "os"
  "io/ioutil"
  osspp "github.com/aboltart/go-support/os"
)

var sd *ServiceDefinition

type ServiceDefinition struct {
  path      string
  names     []string

  loaded    bool
}

func init() {
  sd = new(ServiceDefinition)
  sd.loaded = false
}

// Package body section
func ServiceDefinitionPath() string { return sd.ServiceDefinitionPath() }
func (sd *ServiceDefinition) ServiceDefinitionPath() string {
  sd.loadServiceDefinition()
  return sd.path
}

func ServiceDefinitions() []string { return sd.ServiceDefinitions() }
func (sd *ServiceDefinition) ServiceDefinitions() []string {
  sd.loadServiceDefinition()
  return sd.names
}

func ServiceDefinitionNames() { sd.ServiceDefinitionNames() }
func (sd *ServiceDefinition) ServiceDefinitionNames() {
  sd.loadServiceDefinition()

  for _, name := range sd.ServiceDefinitions() {
    fmt.Println( name )
  }
}

// Private function section
func (sd *ServiceDefinition) loadServiceDefinition() {
  if ! sd.loaded {
    sd.getServiceDefinitionPath()
    sd.getServiceDefinitionNames()
    sd.loaded = true
  }
}

func (sd *ServiceDefinition) getServiceDefinitionPath() {

  if exists_config_file() {
    fmt.Println("Will check what to do. We have config file")
  } else {

    service_path := defaultServiceDefinitionPath()
    if service_path == "" {
      fmt.Println( "Not found service definition path" )
    } else {
      sd.path = service_path
    }
  }
}

func (sd *ServiceDefinition) getServiceDefinitionNames() {
  services := []string{}

  if sd.path != "" {
    files, _ := ioutil.ReadDir(sd.path)

    for _, f := range files {
      services = append(services, f.Name())
    }

    sd.names = services
  }
}

func exists_config_file() bool {
  pwd, _      := os.Getwd()
  config_file := pwd + "/.docker_support"

  return osspp.IsFile(config_file)
}

func defaultServiceDefinitionPath() string {
  pwd, _  := os.Getwd()
  path := pwd + "/compose/services"

  if osspp.IsDirectory(path) {
    return path
  } else {
    return ""
  }
}
