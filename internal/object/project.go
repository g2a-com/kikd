package object

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type Project interface {
	Object

	Files() []string
}

type project struct {
	GenericObject

	Data struct {
		Files     []string
		Variables map[string]string
	} `mapstructure:",squash"`
}

var _ Project = project{}

func NewProject(filename string, data *yaml.Node) (Project, error) {
	p := project{}
	p.GenericObject.metadata = NewMetadata(filename, data)
	err := decode(data, &p)
	return p, err
}

func (p project) Validate(objects ObjectCollection) (err error) {
	for _, project := range objects.GetObjectsByKind(ProjectKind) {
		if project.Metadata() != p.Metadata() {
			err = multierror.Append(err, fmt.Errorf("project is duplicated, it's defined in:\n\t* %s\n\t* %s", project.Metadata(), p.Metadata()))
		}
	}

	return
}

func (p project) PlaceholderValues() map[string]interface{} {
	return map[string]interface{}{
		"Project.Name": p.Name(),
		"Project.Dir":  p.Directory(),
		"Project.Vars": p.Data.Variables,
	}
}

func (p project) Files() []string {
	return p.Data.Files
}
