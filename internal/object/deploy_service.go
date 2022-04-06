package object

import (
	"gopkg.in/yaml.v3"
)

type DeployService struct {
	GenericService

	Deploy struct {
		Releases []*ServiceEntry
	}
}

var _ Object = DeployService{}

func NewDeployService(filename string, data *yaml.Node) (service DeployService, err error) {
	service.GenericObject.metadata = NewMetadata(filename, data)
	err = decode(data, &service)

	service.entries = map[string][]Entry{}
	service.entries[DeployEntryType] = make([]Entry, len(service.Deploy.Releases))
	for i, entry := range service.Deploy.Releases {
		entry.executorKind = DeployerKind
		service.entries[DeployEntryType][i] = entry
	}

	return
}
