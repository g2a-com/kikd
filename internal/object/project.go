package object

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type Project struct {
	GenericObject

	Files     []string
	Variables map[string]string
}

var _ Object = Project{}

func NewProject(filename string, data *yaml.Node) (project Project, err error) {
	project.GenericObject.metadata = NewMetadata(filename, data)
	err = decode(data, &project)
	return
}

func (p Project) Validate(objects ObjectCollection) (err error) {
	for _, project := range objects.GetObjectsByKind(ProjectKind) {
		if project.Metadata() != p.Metadata() {
			err = multierror.Append(err, fmt.Errorf("project is duplicated, it's defined in:\n\t* %s\n\t* %s", project.Metadata(), p.Metadata()))
		}
	}

	return
}
