package object

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type environment struct {
	GenericObject

	DeployServices []string
	Variables      map[string]string
}

var _ Object = environment{}

func NewEnvironment(filename string, data *yaml.Node) (Object, error) {
	e := environment{}
	e.GenericObject.metadata = NewMetadata(filename, data)
	err := decode(data, &e)
	return e, err
}

func (e environment) Validate(c ObjectCollection) (err error) {
	for _, name := range e.DeployServices {
		if c.GetObject(ServiceKind, name) == nil {
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
