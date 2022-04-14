package script

import (
	"fmt"
	"testing"

	"github.com/g2a-com/cicd/internal/object"
	fakelogger "github.com/g2a-com/cicd/internal/utils/fake_logger"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_input_is_passed_down_to_script(t *testing.T) {
	log := fakelogger.New()
	executor := newExecutor(`import("log").print(input)`)
	script := New(executor)
	script.Logger = log

	_, err := script.Run(map[string]interface{}{
		"foo": "bar",
	})

	assert.NoError(t, err)
	assert.Contains(t, log.Messages, fakelogger.Message{Level: "info", Method: "Print", Args: []interface{}{`{foo: "bar"}`}})
}

func Test_returns_results_added_by_script(t *testing.T) {
	executor := newExecutor(`addResult("a", "b"); addResult("c")`)
	script := New(executor)
	script.Logger = fakelogger.New()

	result, err := script.Run(nil)

	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func Test_returns_error_when_script_has_invalid_syntax(t *testing.T) {
	executor := newExecutor(`if`)
	script := New(executor)
	script.Logger = fakelogger.New()

	_, err := script.Run(nil)

	assert.Error(t, err)
}

func Test_returns_error_when_script_is_aborted(t *testing.T) {
	executor := newExecutor(`abort("error")`)
	script := New(executor)
	script.Logger = fakelogger.New()

	_, err := script.Run(nil)

	assert.Error(t, err)
}

// TODO: use object.fakeObject instead (needs to be exported first)
func newExecutor(script string) object.Executor {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(fmt.Sprintf(`{
		kind: Builder,
		name: test,
		schema: {},
		script: %q,
	}`, script)), &node)
	if err != nil {
		panic(err)
	}

	obj, err := object.NewExecutor("file.yaml", &node)
	if err != nil {
		panic(err)
	}

	return obj
}
