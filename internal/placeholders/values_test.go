package placeholders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_merging_values_works(t *testing.T) {
	input1 := map[string]interface{}{"a": "1"}
	input2 := map[string]interface{}{"b": "2"}

	result, err := MergeValues(input1, input2)

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": "1", "b": "2"}, result)
}

func Test_merging_values_results_in_flattened_map(t *testing.T) {
	input := map[string]interface{}{
		"a": map[string]interface{}{
			"b.c": "d",
		},
	}

	result, err := MergeValues(input)

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a.b.c": "d"}, result)
}

func Test_merging_values_fails_if_there_are_name_conflicts_within_single_argument(t *testing.T) {
	input := map[string]interface{}{
		"Foo": map[string]interface{}{
			"Bar": "",
		},
		"foo.bar": "",
	}

	result, err := MergeValues(input)

	assert.Equal(t, err, &DuplicatedPlaceholderError{Name1: ".foo.bar", Name2: ".Foo.Bar"})
	assert.Nil(t, result)
}

func Test_merging_values_overrides_values_if_there_are_name_conflicts_between_maps_provided_in_separate_arguments(t *testing.T) {
	input1 := map[string]interface{}{"a": "1"}
	input2 := map[string]interface{}{"a": "2"}

	result, err := MergeValues(input1, input2)

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": "2"}, result)
}
