package application

import (
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"log"
	"os"
)

func init() {
	filename := ".env"
	if value, exists := os.LookupEnv("PROFILER"); exists {
		filename = ".env." + value
		if _, err := os.Stat(filename); err == nil {
			if err := godotenv.Load(filename); err != nil {
				log.Print("No " + filename + " file found")
			}
		}
	}
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Application struct {
	container        *dig.Container
	serviceProviders []ServiceProviderInterface
}

func (a *Application) Container() *dig.Container {
	return a.container
}

func New(serviceProviders []ServiceProviderInterface) *Application {

	container := dig.New()

	for _, serviceProvider := range serviceProviders {
		serviceProvider.setContainer(container)
		if services := serviceProvider.Provide(); services != nil {
			for _, service := range services {
				err := container.Provide(service)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	app := &Application{
		container:        container,
		serviceProviders: serviceProviders,
	}

	app.serviceProviders = serviceProviders

	return app
}

func (a *Application) Run() error {
	a.onStart()
	defer a.onClose()
	return nil
}

func (a *Application) onStart() {
	for _, serviceProvider := range a.serviceProviders {
		if functions := serviceProvider.OnStart(); functions != nil {
			for _, function := range functions {
				err := a.container.Invoke(function)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (a *Application) onClose() {
	for _, serviceProvider := range a.serviceProviders {
		if functions := serviceProvider.OnClose(); functions != nil {
			for _, function := range functions {
				err := a.container.Invoke(function)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
