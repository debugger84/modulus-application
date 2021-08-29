package application

import (
	"context"
	"go.uber.org/dig"
)

type ServiceProviderInterface interface {
	Provide() []interface{}
	OnStart() []interface{}
	OnClose() []interface{}
	setContainer(*dig.Container)
}

type LoggerInterface interface {
	Warn(ctx context.Context, s string, i ...interface{})
	Info(ctx context.Context, s string, i ...interface{})
	Error(ctx context.Context, s string, i ...interface{})
	Debug(ctx context.Context, s string, i ...interface{})
}

type ServiceProvider struct {
	container *dig.Container
}

func (sp *ServiceProvider) setContainer(container *dig.Container) {
	sp.container = container
}

func (sp *ServiceProvider) Provide() []interface{} {
	var services []interface{}
	return services
}

func (sp *ServiceProvider) OnStart() []interface{} {
	return []interface{}{}
}

func (sp *ServiceProvider) OnClose() []interface{} {
	return []interface{}{}
}
