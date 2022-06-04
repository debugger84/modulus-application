# Modulus Framework
Modulus is a framework for the web development. It allows a developer to create 
a modular monolithic application.

A modular monolith is an approach where you build and deploy a single application, 
but you build application in a way that breaks up the code into independent modules 
for each of the features needed in your application.

Our way to create a modular monolith is having a core with the base interfaces and a system of adding 
as many configurable parts of code as necessary to it. Also, interaction between these parts 
should be organized via this core (framework).


# Module structure
In fact, there is no strong module structure. You are free to implement a module in a way you want. 
Only one restriction to be a module is having a structure that implements ServiceProvider interface.
It will be an entrypoint to your module. We propose to call it ModuleConfig, but actually it is not necessary.

ModuleConfig's ProvidedServices method should return all constructors of structures that you want to mark 
as visible for other parts of the application. All these constructors will be added to the Dependency Injection Container (DIC).
We are using the Uber's container https://github.com/uber-go/dig as DIC in our project.
For example:
```go
//internal/my_module/service/registration.go
type Registration struct {
    logger application.Logger
}

func NewRegistration(logger application.Logger) *Registration {
    return &Registration{logger: logger}
}

//internal/my_module/service/config.go
type ModuleConfig struct {
}

func (s *ModuleConfig) ProvidedServices() []interface{} {
	return []interface{}{		
		service.NewRegistration,
	}
}
```
By the way, all dependencies in the constructor will be resolved automatically, 
if their constructors are returned from any ProvidedServices method in any module.  

# Getting dependencies inside config
Sometimes it is necessary to get some dependencies in module's entrypoint. In this case your ModuleConfig 
should implement ContainerHolder interface.
For example:
```go
type ModuleConfig struct {
	container *dig.Container
}

func (s *ModuleConfig) SetContainer(container *dig.Container) {
	s.container = container
}

func (s *ModuleConfig) ModuleRoutes() []application.RouteInfo {
	var moduleActions *ModuleActions
	err := s.container.Invoke(func(dep *ModuleActions) {
		moduleActions = dep
	})
	if err != nil {
		panic("cannot instantiate module dependencies" + err.Error())
	}
	return moduleActions.Routes()
}
```

# Parameters injecting 
All entrypoint of modules are configurations for modules, and can hold some values.
We propose to set values from env variables. Also, it is nice to have description of these variables
with default values in the .env.dist file in the root of a module.

Let your ModuleConfig implement the ConfigInitializer interface and fill all values there: 
For example:
```go
type ModuleConfig struct {
	apiUrl          string
}

func (s *ModuleConfig) InitConfig(config application.Config) error {	
	if s.apiUrl == "" {
		s.apiUrl = config.GetEnv("MODULE_NAME_API_URL")
	}

	return nil
}
```
We propose to prefix your env variables with a module name to prevent 
names intersection of modules variables.

# Routes description
If your module processes some http routes it is necessary to implement 
the HttpRoutesInitializer interface to return all supported routes.
For example:
```go
//internal/my_module/actions.go
type ModuleActions struct {
    routes *application.Routes
}

func NewModuleActions(
    registerAction *action2.RegisterAction,
) *ModuleActions {
    routes := application.NewRoutes()
    routes.Post(
        "/users",
        registerAction.Handle,
    )
    
    return &ModuleActions{
        routes: routes,
    }
}

func (a *ModuleActions) Routes() []application.RouteInfo {
    return a.routes.GetRoutesInfo()
}

//internal/my_module/config.go

func (s *ModuleConfig) ModuleRoutes() []application.RouteInfo {
	var moduleActions *ModuleActions
	err := s.container.Invoke(func(dep *ModuleActions) {
		moduleActions = dep
	})
	if err != nil {
		panic("cannot instantiate module dependencies" + err.Error())
	}
	return moduleActions.Routes()
}
```

# Application lifecycle events
Any application has own lifecycle, divided to 5 steps: 
#Gather all dependencies from modules
#Config initialization
#Routes initialization
#Running the application
#Closing the application

Configuration module reactions on the first 3 steps has been described previously in the document.
If you want to start for example a server in your application, and, for example, 
release some resources in the end of the application running, than implement interfaces
StartApplicationListener and CloseApplicationListener
For example:
```go

func (s *ServiceProvider) OnStart() error {
	var router *Router
	err := s.container.Invoke(func(dep *Router) error {
		router = dep
		return nil
	})
	if err != nil {
		return err
	}

	return router.Run()
}

func (s *ServiceProvider) OnClose() error {
    var db *Db
    err := s.container.Invoke(func(dep *Db) error {
        router = dep
        return nil
    })
    if err != nil {
        return err
    }
	db.Close()
	return nil
}
```
