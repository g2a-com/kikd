package placeholders

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var nameRegexp *regexp.Regexp = regexp.MustCompile(`^(\.[A-Za-z0-9_]+)+$`)

func normalizeValues(values map[string]interface{}) (map[string]string, error) {
	maps := []map[string]interface{}{values}
	prefixes := []string{""}
	names := map[string]string{}
	result := map[string]string{}

	for i := 0; i < len(maps); i++ {
		for k, v := range maps[i] {
			name := prefixes[i] + "." + k

			switch val := v.(type) {
			case string:
				id := strings.ToLower(name)
				if duplicate, ok := names[id]; ok {
					names := sort.StringSlice([]string{duplicate, name})
					return nil, &DuplicatedPlaceholderError{names[0], names[1]}
				}
				result[name[1:]] = val
				names[id] = name

			case map[string]string:
				prefixes = append(prefixes, name)
				maps = append(maps, toMapStringInterface(val))

			case map[string]interface{}:
				prefixes = append(prefixes, name)
				maps = append(maps, val)

			default:
				return nil, fmt.Errorf("invalid params type: %T", v)
			}
		}
	}

	return result, nil
}

func MergeValues(valuesMaps ...map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for _, values := range valuesMaps {
		normalized, err := normalizeValues(values)
		if err != nil {
			return nil, err
		}
		for k, v := range normalized {
			result[k] = v
		}
	}

	return result, nil
}

type valuesCollection struct {
	ids    []string
	values map[string]string
	names  map[string]string
}

func newValuesCollection(values map[string]interface{}) (*valuesCollection, error) {
	collection := &valuesCollection{
		ids:    []string{},
		values: map[string]string{},
		names:  map[string]string{},
	}

	// Flatten values, it also checks for duplicates
	normalized, err := normalizeValues(values)
	if err != nil {
		return nil, err
	}

	// Convert normalized values (map of strings) to internal representation
	for k, v := range normalized {
		name := "." + k
		id := strings.ToLower(name)
		collection.ids = append(collection.ids, id)
		collection.values[id] = v
		collection.names[id] = name
	}

	// Sort IDs to ensure errors are always the same for given values.
	collection.ids = sort.StringSlice(collection.ids)

	// Validate names.
	for _, id := range collection.ids {
		if !nameRegexp.MatchString(collection.names[id]) {
			return collection, &InvalidPlaceholderNameError{collection.names[id]}
		}
	}

	// Expland placeholders.
	err = collection.expandPlaceholders()
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// Get returns value for given placeholder name.
func (v *valuesCollection) Get(name string) (string, error) {
	id := strings.ToLower(name)
	value, ok := v.values[id]
	if !ok {
		validNames := make([]string, 0, len(v.names))
		for _, n := range v.names {
			validNames = append(validNames, n)
		}
		validNames = sort.StringSlice(validNames)
		return "", &MissingPlaceholderError{name, validNames}
	}
	return value, nil
}

// expandPlaceholders replaces placeholders within values and checks for cycles.
func (v *valuesCollection) expandPlaceholders() (err error) {
	var replace ReplaceFunc

	stack := []string{}
	replace = func(name string) (value string, err error) {
		value, err = v.Get(name)
		if err != nil {
			return
		}

		for i, parent := range stack {
			if strings.EqualFold(parent, name) {
				err = &CyclicPlaceholderError{append(stack[i:], name)}
				return
			}
		}

		if containsMarkers(value) {
			stack = append(stack, name)
			value, err = replaceMarkers(value, replace)
			stack = stack[0 : len(stack)-1]
		}

		return
	}

	for _, id := range v.ids {
		v.values[id], err = replaceMarkers(v.values[id], replace)
		if err != nil {
			return
		}
	}

	return
}

// toMapStringInterface converts map of strings to map of interfaces
func toMapStringInterface(v map[string]string) map[string]interface{} {
	result := make(map[string]interface{}, len(v))
	for k, v := range v {
		result[k] = v
	}
	return result
}
