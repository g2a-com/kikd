package environment

import (
	"testing"

	"github.com/g2a-com/cicd/internal/component"
	"github.com/stretchr/testify/assert"
)

func Test_unmarshalling_empty_environment(t *testing.T) {
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Environment,
		name: test,
	}`)

	result, err := NewEnvironment("dir/file.yaml", input)

	assert.NoError(t, err)
	assert.Equal(t, "dir/file.yaml", result.Metadata().Filename())
	assert.Equal(t, component.EnvironmentKind, result.Kind())
	assert.Equal(t, "test", result.Name())
	assert.Equal(t, "dir", result.Directory())
	assert.Equal(t, `environment "test"`, result.DisplayName())
}

func Test_validating_empty_environment_passes(t *testing.T) {
	collection := component.FakeCollection{}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Environment,
		name: test,
	}`)

	environment, _ := NewEnvironment("dir/file.yaml", input)
	err := environment.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_environment_with_deploy_services_containing_known_services_passes(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ServiceKind, FakeName: "known"},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Environment,
		name: test,
		deployServices: [ known ],
	}`)

	environment, _ := NewEnvironment("dir/file.yaml", input)
	err := environment.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_environment_with_deploy_services_containing_unknown_services_fails(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ServiceKind, FakeName: "known"},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Environment,
		name: test,
		deployServices: [ unknown ],
	}`)

	environment, _ := NewEnvironment("dir/file.yaml", input)
	err := environment.Validate(collection)

	assert.Error(t, err)
}
