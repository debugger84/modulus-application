package application

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"testing"
)

func TestNewApplication(t *testing.T) {
	sp := &TestSp{}
	app := New([]ServiceProviderInterface{sp})
	var dp *TestDependency
	err := app.Container().Invoke(func(dep *TestDependency) error {
		dp = dep
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "test", dp.TestData)
}

func TestRunApplication(t *testing.T) {
	sp := &TestSp{}
	app := New([]ServiceProviderInterface{sp})
	err := app.Run()
	assert.Nil(t, err)
	var dp *TestDependency
	err = app.Container().Invoke(func(dep *TestDependency) error {
		dp = dep
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "test", dp.TestData)
}

type TestDependency struct {
	TestData string
}

func NewTestDependency() *TestDependency {
	return &TestDependency{TestData: "test"}
}

type TestSp struct {
	container *dig.Container
}

func (t TestSp) Provide() []interface{} {
	return []interface{}{
		NewTestDependency,
	}
}

func (t TestSp) OnStart() []interface{} {
	return nil
}

func (t TestSp) OnClose() []interface{} {
	return nil
}

func (t TestSp) setContainer(container *dig.Container) {
	t.container = container
}
