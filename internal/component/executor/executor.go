package executor

import (
	"github.com/g2a-com/cicd/internal/component"
	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v3"
)

type Executor interface {
	component.Component
	Schema() *jsonschema.Schema
	Script() string
}

type executor struct {
	component.Backbone

	Data struct {
		Script string
	} `mapstructure:",squash"`
	schema jsonschema.Schema
}

var _ Executor = executor{}

func NewExecutor(filename string, data *yaml.Node) (Executor, error) {
	e := executor{}
	e.Backbone.SetMetadata(component.NewMetadata(filename, data))
	err := component.Decode(data, &e)
	return e, err
}

func (e executor) Schema() *jsonschema.Schema {
	return &e.schema
}

func (e executor) Script() string {
	return e.Data.Script
}
