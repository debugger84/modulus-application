package application

import (
	"go.uber.org/dig"
	"os"
	"strconv"
	"strings"
)

// ServiceProvider describes all services of a module in the dependency injection container
// it is the first step of the application running
type ServiceProvider interface {
	// ProvidedServices returns a list of constructors that presents all services of the module.
	// All of them will be placed in the dependency injection container
	ProvidedServices() []interface{}
}

// ContainerHolder allows module config to have a link to the dependency injection container
type ContainerHolder interface {
	// SetContainer receives the dependency injection container from application
	SetContainer(*dig.Container)
}

// ConfigInitializer if service provider implements this method it will be called after
// providing dependencies of the module
type ConfigInitializer interface {
	// InitConfig is called for each module to initialize module's variables
	InitConfig(config Config) error
}

// HttpRoutesInitializer if service provider implements this method it will be called after
// initializing the configuration and its result will be added to the http routes
// listened by the application router
type HttpRoutesInitializer interface {
	// HttpRoutesInitializer Returns a set of http routes processed by the module
	ModuleRoutes() []RouteInfo
}

// StartApplicationListener if service provider implements this method it will be called after
// initializing the routes
type StartApplicationListener interface {
	// OnStart Starts module's application such as a web-server
	OnStart() error
}

// CloseApplicationListener if service provider implements this method it will be called after
// stopping the application
type CloseApplicationListener interface {
	// OnClose may close some resources of a module, for example a db connection
	OnClose() error
}

type Config struct {
	appEnv string
}

const (
	TestEnv = "test"
	DevEnv  = "dev"
	ProdEnv = "prod"

	defaultEnv = TestEnv
)

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ProvidedServices() []interface{} {
	return []interface{}{
		func() *Config { return c },
	}
}

func (c *Config) AppEnv() string {
	if c.appEnv == "" {
		c.appEnv = c.GetEnv("APP_ENV")
	}
	return c.appEnv
}

func (c *Config) AppEnvIsProd() bool {
	return c.AppEnv() == ProdEnv
}

func (c *Config) GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	panic("The key " + key + " is not exists in the .env file")
}

func (c *Config) GetEnvAsInt(name string) int {
	valueStr := c.GetEnv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	panic("The value of the key " + name + " in the .env file should be Integer")
}

func (c *Config) GetEnvAsBool(name string) bool {
	valStr := c.GetEnv(name)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	panic("The value of the key " + name + " in the .env file should be Boolean")
}

func (c *Config) GetEnvAsSlice(name string, sep string) []string {
	valStr := c.GetEnv(name)

	val := strings.Split(valStr, sep)

	return val
}
