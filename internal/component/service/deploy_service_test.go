package service

import (
	"testing"

	"github.com/g2a-com/cicd/internal/component"
	"github.com/stretchr/testify/assert"
)

func Test_unmarshalling_empty_deploy_service(t *testing.T) {
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Service,
		name: test,
	}`)

	result, err := NewDeployService("dir/file.yaml", input)

	assert.NoError(t, err)
	assert.Equal(t, "dir/file.yaml", result.Metadata().Filename())
	assert.Equal(t, component.ServiceKind, result.Kind())
	assert.Equal(t, "test", result.Name())
	assert.Equal(t, "dir", result.Directory())
	assert.Equal(t, `service "test"`, result.DisplayName())
}

func Test_validating_empty_deploy_service_passes(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Service,
		name: test,
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_deploy_service_using_unknown_deployer_fails(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
		component.FakeComponent{FakeKind: component.DeployerKind, FakeName: "known", FakeSchema: "{}"},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [ { unknown: {} } ],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_deploy_service_using_known_deployer_passes(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
		component.FakeComponent{FakeKind: component.DeployerKind, FakeName: "known", FakeSchema: "{}"},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [ { known: {} } ],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_deploy_service_with_releases_entry_not_matching_deployer_schema_fails(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
		component.FakeComponent{FakeKind: component.DeployerKind, FakeName: "type", FakeSchema: `{ "required": ["foo"] }`},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{
			type: {},
		}],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_deploy_service_with_releases_entry_matching_deployer_schema_passes(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
		component.FakeComponent{FakeKind: component.DeployerKind, FakeName: "type", FakeSchema: `{ "required": ["foo"] }`},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{
			type: { foo: true },
		}],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_getting_entry_types_list_from_deploy_service_works(t *testing.T) {
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	result := service.EntryTypes()

	assert.Equal(t, []string{"deploy"}, result)
}

func Test_getting_deploy_entries_returns_only_entries_defined_in_releases_property(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{ build: spec }],
		releases: [{ release1: spec1 }, { release2: spec2 }],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	result := service.Entries(component.DeployEntryType)

	assert.Len(t, result, 2)
	assert.Equal(t, 0, result[0].Index())
	assert.Equal(t, component.DeployerKind, result[0].ExecutorKind())
	assert.Equal(t, "release1", result[0].ExecutorName())
	assert.Equal(t, "spec1", result[0].Spec(collection))
	assert.Equal(t, 1, result[1].Index())
	assert.Equal(t, component.DeployerKind, result[1].ExecutorKind())
	assert.Equal(t, "release2", result[1].ExecutorName())
	assert.Equal(t, "spec2", result[1].Spec(collection))
}

func Test_getting_entries_of_unknown_type_from_deploy_service_returns_empty_slice(t *testing.T) {
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	result := service.Entries("unknown")

	assert.Empty(t, result)
}

func Test_getting_deploy_entry_spec_works(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{ name: spec }],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	entries := service.Entries(component.DeployEntryType)
	result := entries[0].Spec(collection)

	assert.Equal(t, "spec", result)
}

func Test_getting_deploy_entry_spec_fills_placeholders_using_values_from_service_environment_project_and_options(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{
			FakeKind:              component.ProjectKind,
			FakePlaceholderValues: map[string]interface{}{"Projects.Foo": "1"},
		},
		component.FakeComponent{
			FakeKind:              component.EnvironmentKind,
			FakePlaceholderValues: map[string]interface{}{"Environment.Bar": "2"},
		},
		component.FakeComponent{
			FakeKind:              component.OptionsKind,
			FakePlaceholderValues: map[string]interface{}{"Options.Egg": "3"},
		},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{ name: "{{.Service.Name}} {{.Projects.Foo}} {{.Environment.Bar}} {{.Options.Egg}}" }],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	entries := service.Entries(component.DeployEntryType)
	result := entries[0].Spec(collection)

	assert.Equal(t, "test 1 2 3", result)
}

func Test_validating_deploy_service_with_entries_fails_when_there_is_no_project_in_the_collection(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.EnvironmentKind},
		component.FakeComponent{FakeKind: component.OptionsKind},
		component.FakeComponent{FakeKind: component.DeployerKind, FakeName: "test"},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{ test: "" }],
	}`)

	service, _ := NewDeployService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}
