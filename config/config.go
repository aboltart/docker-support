package config

import (
  "fmt"
  "os"
  "strings"
  "github.com/spf13/viper"
)

var c *Config
const multipleComposeServicesKey = "multiple_compose_services"

// Define Config instance
type Config struct {
  conf          map[string]interface{} //Like a JSON object.  interface{} mean that value can be any type of value.
  global        map[string]interface{}
  environments  map[string]interface{}
  environment string

  loaded        bool
}

func init() {
  c = New()
}

// Returns an initialized Config instance.
func New() *Config {
  c := new(Config)
  c.conf          = make(map[string]interface{})
  c.global        = make(map[string]interface{})
  c.environments  = make(map[string]interface{})

  c.loaded        = false
  return c
}

func LoadEnvShellVariables() { c.LoadEnvShellVariables() }
func (c *Config) LoadEnvShellVariables() {
  environment       := Environment()
  envVariables      := EnvVariables(environment)
  envShellVariables := EnvShellVariables(envVariables)

  BuildEnvShellVariables(envShellVariables)
}

func BuildEnvShellVariables(envShellVariables map[string]interface{}) { c.BuildEnvShellVariables(envShellVariables) }
func (c *Config) BuildEnvShellVariables(envShellVariables map[string]interface{}) {

  for key, value := range envShellVariables {

    // If environment variable is not set
    if os.Getenv(key) == "" {
      os.Setenv(key, fmt.Sprint(value))
    }
  }
}

func EnvShellVariables(envVariables map[string]interface{}) map[string]interface{} { return c.EnvShellVariables(envVariables) }
func (c *Config) EnvShellVariables(envVariables map[string]interface{}) map[string]interface{} {

  variables := make(map[string]interface{})

  for key, value := range envVariables {
    if key != multipleComposeServicesKey {
      variables[strings.ToUpper(key)] = value
    }
  }

  return variables
}

func MultipleComposeServices() map[string]interface{} { return c.MultipleComposeServices() }
func (c *Config) MultipleComposeServices() map[string]interface{} {
  environment       := Environment()
  envVariables      := EnvVariables(environment)

  return GetMultipleComposeServices(envVariables)
}

func GetMultipleComposeServices(envVariables map[string]interface{}) map[string]interface{} { return c.GetMultipleComposeServices(envVariables) }
func (c *Config) GetMultipleComposeServices(envVariables map[string]interface{}) map[string]interface{} {

  var multipleComposeServiceValue map[string]interface{}

  if value, ok := envVariables[multipleComposeServicesKey]; ok {
    multipleComposeServiceValue = value.(map[string]interface{})
  }

  return multipleComposeServiceValue
}

func EnvVariables(env string) map[string]interface{} { return c.EnvVariables(env) }
func (c *Config) EnvVariables(env string) map[string]interface{} {

  envVariables := make(map[string]interface{})

  // fmt.Println( "------------------------------------------")
  // fmt.Println( "-->>> conf: ", c.conf )
  // fmt.Println( "-->>> global: ", c.global )
  // fmt.Println( "-->>> environments: ", c.environments )
  // fmt.Println( "------------------------------------------")

  envVariables["environment"] = env

  for k, v := range c.global {
    envVariables[k] = v
  }

  if _, ok :=  c.environments[env]; ok {
    for k, v := range c.environments[env].(map[string]interface{}) {
      envVariables[k] = v
    }
  }

  return envVariables
}

func Environment() string { return c.Environment() }
func (c *Config) Environment() string {
  collectConfig()
  return c.environment
}

// Private functions
func read() { c.read() }
func (c *Config) read() {
  for _, key := range keys() {
    c.conf[key] = viper.GetStringMap(key)
  }
}

func collectConfig() { c.collectConfig() }
func (c *Config) collectConfig() {

  if ! c.loaded {
    read()

    for key, value := range c.conf {
      // For len cannot use interface map
      // To acces we use a type assertion to access `value`'s underlying map[string]interface{}:
      if len(value.(map[string]interface{})) == 0 || key == multipleComposeServicesKey {
        c.global[key] = viper.Get(key)
      } else {

        // Build Env specific variable for collecting its variables
        env := make(map[string]interface{})
        for k, v := range value.(map[string]interface{}) {
          env[k] = v
        }
        // Append to Config struct object
        c.environments[key] = env
      }
    }

    // Detect environment
    // First try to get if exists set environment OS variable
    environment := os.Getenv("ENVIRONMENT")

    // Try to get environment from global config
    if environment == "" {
      // Convert value from interface (any type) to string. Nil cannot be converted,
      // In such cases error
      envValue, ok := c.global["environment"].(string); if !ok {
        fmt.Println("--->>> No environment configured....")
      } else {
        environment = envValue
      }
    }

    // If no environment set default
    if environment == "" {
      environment = "development"
      fmt.Println("      .... Set to", environment)
    }

    c.environment = environment

    c.loaded = true

    // fmt.Println("--->>> c.global ",c.global)
    // fmt.Println("--->>> c.environments ",c.environments)
  }

}

func viper_config() {
  viper.SetConfigName("abc")
  viper.AddConfigPath(".")
}

func exists_file() bool {
  viper_config()

  err := viper.ReadInConfig()
  if err != nil {
    fmt.Errorf("Fatal error config file: %s \n", err)
    return false
  }

  return true
}

func keys() []string {
  if exists_file() {
    return viper.AllKeys()
  } else {
    return []string{}
  }
}
