package object

import (
	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v3"
)

type Executor interface {
	Object
	Schema() *jsonschema.Schema
	Script() string
}

type executor struct {
	GenericObject

	Data struct {
		Script string
	} `mapstructure:",squash"`
	schema jsonschema.Schema
}

var _ Executor = executor{}

func NewExecutor(filename string, data *yaml.Node) (Executor, error) {
	e := executor{}
	e.GenericObject.metadata = NewMetadata(filename, data)
	err := decode(data, &e)
	return e, err
}

func (e executor) Schema() *jsonschema.Schema {
	return &e.schema
}

func (e executor) Script() string {
	return e.Data.Script
}
