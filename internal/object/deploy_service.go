package object

import (
	"context"
	"fmt"
	"strings"

	"github.com/g2a-com/cicd/internal/placeholders"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type deployService struct {
	GenericService

	Deploy struct {
		Releases []*deployServiceEntry
	}
}

var _ Object = deployService{}

func NewDeployService(filename string, data *yaml.Node) (Object, error) {
	service := deployService{}
	service.GenericObject.metadata = NewMetadata(filename, data)
	err := decode(data, &service)

	service.entries = map[string][]Entry{}
	service.entries[DeployEntryType] = make([]Entry, len(service.Deploy.Releases))
	for i, entry := range service.Deploy.Releases {
		entry.service = service
		entry.executorKind = DeployerKind
		service.entries[DeployEntryType][i] = entry
	}

	return service, err
}

type deployServiceEntry struct {
	executorKind Kind
	service      Object
	Data         struct {
		Index int
		Type  string
		Spec  interface{}
	} `mapstructure:",squash"`
}

func (e *deployServiceEntry) Index() int {
	return e.Data.Index
}

func (e *deployServiceEntry) ExecutorKind() Kind {
	return e.executorKind
}

func (e *deployServiceEntry) ExecutorName() string {
	return e.Data.Type
}

func (e *deployServiceEntry) Validate(objects ObjectCollection) error {
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

func (e *deployServiceEntry) Spec(objects ObjectCollection) interface{} {
	spec, err := e.spec(objects)
	if err != nil {
		// Errors should have been handled during validation phase.
		panic(err)
	}
	return spec
}

func (e *deployServiceEntry) spec(b ObjectCollection) (interface{}, error) {
	project := b.GetUniqueObject(ProjectKind)
	if project == nil {
		return nil, fmt.Errorf("cannot find project")
	}
	environment := b.GetUniqueObject(EnvironmentKind)
	if environment == nil {
		return nil, fmt.Errorf("cannot find environment")
	}
	options := b.GetUniqueObject(OptionsKind)
	if options == nil {
		return nil, fmt.Errorf("cannot find options")
	}

	values, err := placeholders.MergeValues(
		project.PlaceholderValues(),
		environment.PlaceholderValues(),
		options.PlaceholderValues(),
		e.service.PlaceholderValues(),
	)
	if err != nil {
		return nil, err
	}

	return placeholders.ReplaceWithValues(e.Data.Spec, values)
}
