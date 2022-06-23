package component

import (
	"fmt"
	"github.com/g2a-com/cicd/internal/schema"
	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v3"
)

type FakeComponent struct {
	FakeMetadata          metadata
	FakeKind              Kind
	FakeName              string
	FakeDirectory         string
	displayName           string
	FakeSchema            string
	script                string
	entryTypes            []string
	entries               []Entry
	FakePlaceholderValues map[string]interface{}
}

var _ Component = FakeComponent{}

func (o FakeComponent) Name() string {
	return o.FakeName
}

func (o FakeComponent) Kind() Kind {
	return o.FakeKind
}

func (o FakeComponent) Metadata() Metadata {
	return o.FakeMetadata
}

func (o FakeComponent) Directory() string {
	return o.FakeDirectory
}

func (o FakeComponent) DisplayName() string {
	return o.displayName
}

func (o FakeComponent) Validate(ObjectCollection) error {
	return nil
}

func (o FakeComponent) Schema() *jsonschema.Schema {
	if o.FakeSchema == "" {
		return jsonschema.Must("{}")
	}
	return jsonschema.Must(o.FakeSchema)
}

func (o FakeComponent) Script() string {
	return o.script
}

func (o FakeComponent) EntryTypes() []string {
	return o.entryTypes
}

func (o FakeComponent) Entries(string) []Entry {
	return o.entries
}

func (o FakeComponent) PlaceholderValues() map[string]interface{} {
	return o.FakePlaceholderValues
}

// testInput validates input against schema and returns it back. Use only in
// tests.
func PrepareTestInput(input string) *yaml.Node {
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

type FakeCollection []Component

func (c FakeCollection) GetObject(kind Kind, name string) Component {
	for _, o := range c {
		if o.Kind() == kind && o.Name() == name {
			return o
		}
	}
	return nil
}

func (c FakeCollection) GetUniqueObject(kind Kind) Component {
	var result Component
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

func (c FakeCollection) GetObjectsByKind(kind Kind) []Component {
	result := []Component{}
	for _, o := range c {
		if o.Kind() == kind {
			result = append(result, o)
		}
	}
	return result
}
