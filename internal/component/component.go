package component

import (
	"encoding/json"
	"fmt"
	"github.com/g2a-com/cicd/internal/component/scheme"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Component interface {
	Metadata() Metadata
	Name() string
	Kind() Kind
	DisplayName() string
	Directory() string
	Validate(ObjectCollection) error
	EntryTypes() []string
	Entries(string) []Entry
	PlaceholderValues() map[string]interface{}
}

type ObjectCollection interface {
	GetObject(kind Kind, name string) Component
	GetUniqueObject(kind Kind) Component
	GetObjectsByKind(kind Kind) []Component
}

type Backbone struct {
	metadata Metadata
	Data     struct {
		Kind Kind
		Name string
	} `mapstructure:",squash" yaml:",inline"`
}

func (o Backbone) SetMetadata(metadata Metadata) {
	o.metadata = metadata
}

func (o Backbone) Name() string {
	return o.Data.Name
}

func (o Backbone) Kind() Kind {
	return o.Data.Kind
}

func (o Backbone) Metadata() Metadata {
	return o.metadata
}

func (o Backbone) Directory() string {
	return filepath.Dir(o.metadata.Filename())
}

func (o Backbone) DisplayName() string {
	return fmt.Sprintf("%s %q", strings.ToLower(string(o.Kind())), o.Name())
}

func (o Backbone) Validate(ObjectCollection) error {
	panic("Function not implemented")
	return nil
}

func (o Backbone) EntryTypes() []string {
	panic("Function not implemented")
	return []string{}
}

func (o Backbone) Entries(_ string) []Entry {
	panic("Function not implemented")
	return []Entry{}
}

func (o Backbone) PlaceholderValues() map[string]interface{} {
	panic("Function not implemented")
	return map[string]interface{}{}
}

func Decode(data *yaml.Node, result interface{}) (err error) {
	var aux interface{}

	err = data.Decode(&aux)
	if err != nil {
		return
	}

	aux, err = scheme.ToInternal(aux)
	if err != nil {
		return
	}

	decoderConfig := &mapstructure.DecoderConfig{
		ErrorUnused: false,
		Squash:      true,
		Result:      result,
		DecodeHook: func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			result := reflect.New(t).Interface()
			unmarshaller, ok := result.(json.Unmarshaler)
			if !ok {
				return data, nil
			}
			if err := unmarshaller.UnmarshalJSON([]byte(data.(string))); err != nil {
				return nil, err
			}
			return result, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}

	return decoder.Decode(aux)
}
