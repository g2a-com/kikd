package project

import (
	"testing"

	"github.com/g2a-com/cicd/internal/component"
	"github.com/stretchr/testify/assert"
)

func Test_unmarshalling_empty_project(t *testing.T) {
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Project,
		name: test,
	}`)

	result, err := NewProject("dir/file.yaml", input)

	assert.NoError(t, err)
	assert.Equal(t, "dir/file.yaml", result.Metadata().Filename())
	assert.Equal(t, component.ProjectKind, result.Kind())
	assert.Equal(t, "test", result.Name())
	assert.Equal(t, "dir", result.Directory())
	assert.Equal(t, `project "test"`, result.DisplayName())
}

func Test_validating_empty_project_passes(t *testing.T) {
	collection := component.FakeCollection{}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Project,
		name: test,
	}`)

	project, _ := NewProject("dir/file.yaml", input)
	err := project.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_duplicated_project_fails(t *testing.T) {
	collection := component.FakeCollection{
		component.FakeComponent{FakeKind: component.ProjectKind},
	}
	input := component.PrepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Project,
		name: test,
	}`)

	project, _ := NewProject("dir/file.yaml", input)
	err := project.Validate(collection)

	assert.Error(t, err)
}
