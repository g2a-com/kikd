package environment

import (
	"fmt"
	"github.com/g2a-com/cicd/internal/component"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type environment struct {
	component.Backbone

	DeployServices []string
	Variables      map[string]string
}

var _ component.Component = environment{}

func NewEnvironment(filename string, data *yaml.Node) (component.Component, error) {
	e := environment{}
	e.Backbone.SetMetadata(component.NewMetadata(filename, data))
	err := component.Decode(data, &e)
	return e, err
}

func (e environment) Validate(c component.ObjectCollection) (err error) {
	for _, name := range e.DeployServices {
		if c.GetObject(component.ServiceKind, name) == nil {
			err = multierror.Append(err, fmt.Errorf("missing service %q deployed to environment %q defined in the file:\n\t  %s", name, e.Name(), e.Metadata().Filename()))
		}
	}
	return
}

func (e environment) PlaceholderValues() map[string]interface{} {
	return map[string]interface{}{
		"Environment.Name": e.Name(),
		"Environment.Dir":  e.Directory(),
		"Environment.Vars": e.Variables,
	}
}
