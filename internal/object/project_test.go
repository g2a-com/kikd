package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unmarshalling_empty_project(t *testing.T) {
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Project,
		name: test,
	}`)

	result, err := NewProject("dir/file.yaml", input)

	assert.NoError(t, err)
	assert.Equal(t, "dir/file.yaml", result.Metadata().Filename())
	assert.Equal(t, ProjectKind, result.Kind())
	assert.Equal(t, "test", result.Name())
	assert.Equal(t, "dir", result.Directory())
	assert.Equal(t, `project "test"`, result.DisplayName())
}

func Test_validating_empty_project_passes(t *testing.T) {
	collection := fakeCollection{}
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Project,
		name: test,
	}`)

	project, _ := NewProject("dir/file.yaml", input)
	err := project.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_duplicated_project_fails(t *testing.T) {
	collection := fakeCollection{
		fakeObject{kind: ProjectKind},
	}
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Project,
		name: test,
	}`)

	project, _ := NewProject("dir/file.yaml", input)
	err := project.Validate(collection)

	assert.Error(t, err)
}
