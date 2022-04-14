package object

import (
	"sort"

	"github.com/hashicorp/go-multierror"
)

type GenericService struct {
	GenericObject
	entries map[string][]Entry
}

var _ Object = GenericService{}

func (s GenericService) Validate(c ObjectCollection) (err error) {
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

func (s GenericService) Entries(entryType string) []Entry {
	result := make([]Entry, len(s.entries[entryType]))
	copy(result, s.entries[entryType])
	return result
}

func (s GenericService) PlaceholderValues() map[string]interface{} {
	return map[string]interface{}{
		"Service.Name": s.Name(),
		"Service.Dir":  s.Directory(),
	}
}
