package object

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v3"
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

type Service struct {
	GenericObject

	Build struct {
		Tags      []*ServiceEntry
		Artifacts struct {
			ToBuild []*ServiceEntry
			ToPush  []*ServiceEntry
		}
	}
	Deploy struct {
		Releases []*ServiceEntry
	}
	Run struct {
		Tasks map[string][]*ServiceEntry
	}
	entries map[string][]Entry
}

var _ Object = Service{}

func NewService(filename string, data *yaml.Node) (service Service, err error) {
	service.GenericObject.metadata = NewMetadata(filename, data)
	err = decode(data, &service)

	service.entries = map[string][]Entry{}
	service.entries[TagEntryType] = make([]Entry, len(service.Build.Tags))
	for i, entry := range service.Build.Tags {
		entry.executorKind = TaggerKind
		service.entries[TagEntryType][i] = entry
	}
	service.entries[BuildEntryType] = make([]Entry, len(service.Build.Artifacts.ToBuild))
	for i, entry := range service.Build.Artifacts.ToBuild {
		entry.executorKind = BuilderKind
		service.entries[BuildEntryType][i] = entry
	}
	service.entries[PushEntryType] = make([]Entry, len(service.Build.Artifacts.ToPush))
	for i, entry := range service.Build.Artifacts.ToPush {
		entry.executorKind = PusherKind
		service.entries[PushEntryType][i] = entry
	}
	service.entries[DeployEntryType] = make([]Entry, len(service.Deploy.Releases))
	for i, entry := range service.Deploy.Releases {
		entry.executorKind = DeployerKind
		service.entries[DeployEntryType][i] = entry
	}

	return
}

func (s Service) Validate(c ObjectCollection) (err error) {
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

func (s Service) EntryTypes() []string {
	// FIXME: return types for tasks
	return []string{TagEntryType, BuildEntryType, PushEntryType, DeployEntryType}
}

func (s Service) Entries(entryType string) []Entry {
	result := make([]Entry, len(s.entries[entryType]))
	copy(result, s.entries[entryType])
	return result
}
