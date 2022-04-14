package object

import (
	"context"
	"fmt"
	"strings"

	"github.com/g2a-com/cicd/internal/placeholders"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type buildService struct {
	GenericService

	Build struct {
		Tags      []*buildServiceEntry
		Artifacts struct {
			ToBuild []*buildServiceEntry
			ToPush  []*buildServiceEntry
		}
	}
}

var _ Object = buildService{}

func NewBuildService(filename string, data *yaml.Node) (Object, error) {
	service := buildService{}
	service.GenericObject.metadata = NewMetadata(filename, data)
	err := decode(data, &service)

	service.entries = map[string][]Entry{}
	service.entries[TagEntryType] = make([]Entry, len(service.Build.Tags))
	for i, entry := range service.Build.Tags {
		entry.service = service
		entry.executorKind = TaggerKind
		service.entries[TagEntryType][i] = entry
	}
	service.entries[BuildEntryType] = make([]Entry, len(service.Build.Artifacts.ToBuild))
	for i, entry := range service.Build.Artifacts.ToBuild {
		entry.service = service
		entry.executorKind = BuilderKind
		service.entries[BuildEntryType][i] = entry
	}
	service.entries[PushEntryType] = make([]Entry, len(service.Build.Artifacts.ToPush))
	for i, entry := range service.Build.Artifacts.ToPush {
		entry.service = service
		entry.executorKind = PusherKind
		service.entries[PushEntryType][i] = entry
	}

	return service, err
}

type buildServiceEntry struct {
	executorKind Kind
	service      Object
	Data         struct {
		Index int
		Type  string
		Spec  interface{}
	} `mapstructure:",squash"`
}

func (e *buildServiceEntry) Index() int {
	return e.Data.Index
}

func (e *buildServiceEntry) ExecutorKind() Kind {
	return e.executorKind
}

func (e *buildServiceEntry) ExecutorName() string {
	return e.Data.Type
}

func (e *buildServiceEntry) Validate(objects ObjectCollection) error {
	spec, err := e.spec(objects)
	if err != nil {
		return err
	}

	obj := objects.GetObject(e.ExecutorKind(), e.ExecutorName())

	if obj == nil {
		return fmt.Errorf(
			"missing %s %q used by service %q defined in the file:\n\t  %s",
			strings.ToLower(string(e.ExecutorKind())), e.ExecutorName(), e.service.Name(), e.service.Metadata(),
		)
	}

	executor, ok := obj.(Executor)
	if !ok {
		panic("not an executor")
	}

	schema := executor.Schema()
	result := schema.Validate(context.Background(), spec)

	if len(*result.Errs) > 0 {
		var err error
		for _, msg := range *result.Errs {
			err = multierror.Append(err, fmt.Errorf(
				"%s contains invalid configuration for %s:\n\t  %s\n\t  Definition files:\n\t    %s\n\t    %s",
				e.service.DisplayName(), executor.DisplayName(), msg, e.service.Metadata(), executor.Metadata(),
			))
		}
		return err
	}

	return nil
}

func (e *buildServiceEntry) Spec(objects ObjectCollection) interface{} {
	spec, err := e.spec(objects)
	if err != nil {
		// Errors should have been handled during validation phase.
		panic(err)
	}
	return spec
}

func (e *buildServiceEntry) spec(b ObjectCollection) (interface{}, error) {
	project := b.GetUniqueObject(ProjectKind)
	if project == nil {
		return nil, fmt.Errorf("cannot find project")
	}
	options := b.GetUniqueObject(OptionsKind)
	if options == nil {
		return nil, fmt.Errorf("cannot find options")
	}

	values, err := placeholders.MergeValues(
		project.PlaceholderValues(),
		options.PlaceholderValues(),
		e.service.PlaceholderValues(),
	)
	if err != nil {
		return nil, err
	}

	return placeholders.ReplaceWithValues(e.Data.Spec, values)
}
