package object

import (
	"fmt"

	"github.com/g2a-com/cicd/internal/schema"
	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v3"
)

type fakeObject struct {
	metadata          metadata
	kind              Kind
	name              string
	directory         string
	displayName       string
	schema            string
	script            string
	entryTypes        []string
	entries           []Entry
	placeholderValues map[string]interface{}
}

var _ Object = fakeObject{}

func (o fakeObject) Name() string {
	return o.name
}

func (o fakeObject) Kind() Kind {
	return o.kind
}

func (o fakeObject) Metadata() Metadata {
	return o.metadata
}

func (o fakeObject) Directory() string {
	return o.directory
}

func (o fakeObject) DisplayName() string {
	return o.displayName
}

func (o fakeObject) Validate(ObjectCollection) error {
	return nil
}

func (o fakeObject) Schema() *jsonschema.Schema {
	if o.schema == "" {
		return jsonschema.Must("{}")
	}
	return jsonschema.Must(o.schema)
}

func (o fakeObject) Script() string {
	return o.script
}

func (o fakeObject) EntryTypes() []string {
	return o.entryTypes
}

func (o fakeObject) Entries(string) []Entry {
	return o.entries
}

func (o fakeObject) PlaceholderValues() map[string]interface{} {
	return o.placeholderValues
}

// testInput validates input against schema and returns it back. Use only in
// tests.
func prepareTestInput(input string) *yaml.Node {
	_, err := schema.Validate([]byte(input))
	if err != nil {
		panic(err)
	}
	result := &yaml.Node{}
	err = yaml.Unmarshal([]byte(input), result)
	if err != nil {
		panic(err)
	}
	return result
}

type fakeCollection []Object

func (c fakeCollection) GetObject(kind Kind, name string) Object {
	for _, o := range c {
		if o.Kind() == kind && o.Name() == name {
			return o
		}
	}
	return nil
}

func (c fakeCollection) GetUniqueObject(kind Kind) Object {
	var result Object
	for _, o := range c {
		if o.Kind() == kind {
			if result == nil {
				result = o
			} else {
				panic(fmt.Errorf("duplicated object of %s kind", o.Kind()))
			}
		}
	}
	return result
}

func (c fakeCollection) GetObjectsByKind(kind Kind) []Object {
	result := []Object{}
	for _, o := range c {
		if o.Kind() == kind {
			result = append(result, o)
		}
	}
	return result
}
