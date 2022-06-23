package service

import (
	"github.com/g2a-com/cicd/internal/component"
	"sort"

	"github.com/hashicorp/go-multierror"
)

type GenericService struct {
	component.Backbone
	entries map[string][]component.Entry
}

var _ component.Component = GenericService{}

func (s GenericService) Validate(c component.ObjectCollection) (err error) {
	for _, entryType := range s.EntryTypes() {
		for _, entry := range s.entries[entryType] {
			e := entry.Validate(c)
			if e != nil {
				err = multierror.Append(err, e)
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

func (s GenericService) Entries(entryType string) []component.Entry {
	result := make([]component.Entry, len(s.entries[entryType]))
	copy(result, s.entries[entryType])
	return result
}

func (s GenericService) PlaceholderValues() map[string]interface{} {
	return map[string]interface{}{
		"Service.Name": s.Name(),
		"Service.Dir":  s.Directory(),
	}
}
