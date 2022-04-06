package object

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/qri-io/jsonschema"
)

type ServiceEntry struct {
	executorKind Kind
	Data         struct {
		Index int
		Type  string
		Spec  interface{}
	} `mapstructure:",squash"`
}

func (e *ServiceEntry) Index() int {
	return e.Data.Index
}

func (e *ServiceEntry) ExecutorKind() Kind {
	return e.executorKind
}

func (e *ServiceEntry) ExecutorName() string {
	return e.Data.Type
}

func (e *ServiceEntry) Spec(ObjectCollection) interface{} {
	return e.Data.Spec
}

type GenericService struct {
	GenericObject
	entries map[string][]Entry
}

var _ Object = GenericService{}

func (s GenericService) Validate(c ObjectCollection) (err error) {
	for _, entryType := range s.EntryTypes() {
		for _, entry := range s.entries[entryType] {
			obj := c.GetObject(entry.ExecutorKind(), entry.ExecutorName())

			if obj == nil {
				err = multierror.Append(err, fmt.Errorf(
					"missing %s %q used by service %q defined in the file:\n\t  %s",
					strings.ToLower(string(entry.ExecutorKind())), entry.ExecutorName(), s.Name(), s.Metadata(),
				))
				continue
			}

			executor := obj.(interface {
				Object
				Schema() *jsonschema.Schema
			})

			schema := executor.Schema()
			result := schema.Validate(context.Background(), entry.Spec(c))

			if len(*result.Errs) > 0 {
				for _, e := range *result.Errs {
					err = multierror.Append(err, fmt.Errorf(
						"%s contains invalid configuration for %s:\n\t  %s\n\t  Definition files:\n\t    %s\n\t    %s",
						s.DisplayName(), executor.DisplayName(), e, s.Metadata(), executor.Metadata(),
					))
				}
			}
		}
	}

	return
}

func (s GenericService) EntryTypes() []string {
	result := make([]string, 0, len(s.entries))
	for key := range s.entries {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func (s GenericService) Entries(entryType string) []Entry {
	result := make([]Entry, len(s.entries[entryType]))
	copy(result, s.entries[entryType])
	return result
}
