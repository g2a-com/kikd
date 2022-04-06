package object

import (
	"gopkg.in/yaml.v3"
)

type BuildService struct {
	GenericService

	Build struct {
		Tags      []*ServiceEntry
		Artifacts struct {
			ToBuild []*ServiceEntry
			ToPush  []*ServiceEntry
		}
	}
}

var _ Object = BuildService{}

func NewBuildService(filename string, data *yaml.Node) (service BuildService, err error) {
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

	return
}
