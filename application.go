package application

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"log"
	"os"
	"reflect"
)

type Application struct {
	container     *dig.Container
	moduleConfigs []interface{}
}

func (a *Application) Container() *dig.Container {
	return a.container
}

func New(moduleConfigs []interface{}) *Application {
	container := dig.New()
	app := &Application{
		container: container,
	}
	app.readEnv()

	applicationConfig := NewConfig()

	app.moduleConfigs = append(moduleConfigs, applicationConfig)
	app.fillProvidedServices()

	return app
}

func (a *Application) Run() error {
	a.initConfig()

	a.setDefaultLogger()
	a.setDefaultJsonResponseWriter()

	a.initHttpRoutes()

	a.onStart()
	defer a.onClose()
	return nil
}

func (a *Application) fillProvidedServices() {
	for _, moduleConfig := range a.moduleConfigs {
		if containerHolder, ok := moduleConfig.(ContainerHolder); ok {
			containerHolder.SetContainer(a.container)
		}
		if serviceProvider, ok := moduleConfig.(ServiceProvider); ok {
			if services := serviceProvider.ProvidedServices(); services != nil {
				for _, service := range services {
					err := a.container.Provide(service)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

func (a *Application) initConfig() {
	config := NewConfig()
	for _, serviceProvider := range a.moduleConfigs {
		if routesContainer, ok := serviceProvider.(ConfigInitializer); ok {
			err := routesContainer.InitConfig(*config)
			if err != nil {
				logger := a.getLogger()
				logger.Panic(context.Background(), "Init config error: "+err.Error())
			}
		}
	}
}

func (a *Application) readEnv() {
	filename := ".env"
	if value, exists := os.LookupEnv("APP_ENV"); exists {
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

func (a *Application) initHttpRoutes() {
	moduleName := ""
	defer func() {
		if err := recover(); err != nil {
			logger := a.getLogger()
			logger.Panic(
				context.Background(),
				fmt.Sprint(moduleName, ": routes initialisation failed: ", err),
			)
		}
	}()
	router := a.getRouter()

	if router == nil {
		return
	}
	for _, serviceProvider := range a.moduleConfigs {
		if t := reflect.TypeOf(serviceProvider); t.Kind() == reflect.Ptr {
			moduleName = t.Elem().PkgPath()
		}
		if routesContainer, ok := serviceProvider.(HttpRoutesInitializer); ok {
			router.AddRoutes(routesContainer.ModuleRoutes())
		}
	}
}

func (a *Application) onStart() {
	for _, serviceProvider := range a.moduleConfigs {
		if appListener, ok := serviceProvider.(StartApplicationListener); ok {
			err := appListener.OnStart()
			if err != nil {
				logger := a.getLogger()
				logger.Panic(context.Background(), "Start application error: "+err.Error())
			}
		}
	}
}

func (a *Application) onClose() {
	for _, serviceProvider := range a.moduleConfigs {
		if appListener, ok := serviceProvider.(CloseApplicationListener); ok {
			err := appListener.OnClose()
			if err != nil {
				logger := a.getLogger()
				logger.Panic(context.Background(), "Close application error: "+err.Error())
			}
		}
	}
}

func (a *Application) setDefaultLogger() {
	var logger Logger
	err := a.container.Invoke(func(dep Logger) error {
		logger = dep
		return nil
	})
	if err != nil || logger == nil {
		err := a.container.Provide(NewDefaultLogger)
		if err != nil {
			panic("Default logger cannot be setup")
		}
	}
}

func (a *Application) setDefaultJsonResponseWriter() {
	var logger JsonResponseWriter
	err := a.container.Invoke(func(dep JsonResponseWriter) error {
		logger = dep
		return nil
	})
	if err != nil || logger == nil {
		err := a.container.Provide(NewJsonResponse)
		if err != nil {
			panic("Default json response writer cannot be setup")
		}
	}
}

func (a *Application) getLogger() Logger {
	var logger Logger
	err := a.container.Invoke(func(dep Logger) error {
		logger = dep
		return nil
	})
	if err != nil {
		return NewDefaultLogger()
	}

	return logger
}

func (a *Application) getRouter() Router {
	var router Router
	err := a.container.Invoke(func(dep Router) error {
		router = dep
		return nil
	})
	if err != nil {
		return nil
	}

	return router
}
