package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_migrating_empty_service_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4,
		kind: Service,
		name: test,
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Service,
		name: test
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_service_containing_unsupported_properties_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
		extra: true
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		extra: true
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_service_tagPolicy_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
		build: {
			artifacts: [],
			tagPolicy: {
				name: { config: true },
			}
		}
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{
			name: { config: true }
		}],
		artifacts: []
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_service_artifacts_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
		build: {
			tagPolicy: { tag: {} },
			artifacts: [{
				builder: { foo: bar }
			}]
		}
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ tag: {} }],
		artifacts: [{
			builder: { foo: bar }
		}]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_releases_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
		deploy: {
			releases: [{
				deployer: { foo: bar }
			}]
		}
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{
			deployer: { foo: bar }
		}]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_service_containing_only_hooks_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
		hooks: {
			pre-build: [ test ],
			post-build: [ test ],
			pre-deploy: [ test ],
			post-deploy: [ test ],
		}
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_build_hooks_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
		hooks: {
			pre-build: [pre, build],
			post-build: [post, build],
		},
		build: {
			tagPolicy: { tag: {} },
			artifacts: [{ builder: { foo: bar } }]
		}
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ tag: {} }],
		artifacts: [
			{ script: { sh: "set -e\npre\nbuild\n" }, push: false },
			{ builder: { foo: bar } },
			{ script: { sh: "set -e\npost\nbuild\n" }, push: false },
		]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_deploy_hooks_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
		hooks: {
			pre-deploy: [pre, deploy],
			post-deploy: [post, deploy],
		},
		deploy: {
			releases: [{ deployer: { foo: bar } }]
		}
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [
			{ script: { sh: "set -e\npre\ndeploy\n" } },
			{ deployer: { foo: bar } },
			{ script: { sh: "set -e\npost\ndeploy\n" } },
		]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_service_containing_restricted_property_names_from_v1beta4_to_v2_0(t *testing.T) {
	cases := []string{"artifacts", "releases", "tags", "tasks"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			input := testInput(fmt.Sprintf(`{
				apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
				%s: []
			}`, name))

			migrator := NewMigrator("g2a-cli/v2.0")
			_, err := migrator.Migrate([]byte(input))

			assert.Error(t, err)
		})
	}
}

func Test_migrating_empty_environment_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4,
		kind: Environment,
		name: test
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Environment,
		name: test
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_environment_variables_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Environment, name: test,
		variables: {
			foo: bar,
			egg: spam
		}
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Environment, name: test,
		variables: {
			foo: bar,
			egg: spam
		}
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_environment_deployServices_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Environment, name: test,
		deployServices: [foo, bar]
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Environment, name: test,
		deployServices: [foo, bar]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_empty_project_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4,
		kind: Project
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Project,
		name: project,
		files: [
			services/*/service.yaml,
			environments/*/environment.yaml,
		]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_project_services_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Project,
		services: [ ./ ],
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Project, name: project,
		files: [
			service.yaml,
			environments/*/environment.yaml,
		]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_project_environments_from_v1beta4_to_v2_0(t *testing.T) {
	input := testInput(`{
		apiVersion: g2a-cli/v1beta4, kind: Project,
		environments: [ ./ ],
	}`)
	expected := testInput(`{
		apiVersion: g2a-cli/v2.0, kind: Project, name: project,
		files: [
			services/*/service.yaml,
			environment.yaml,
		]
	}`)

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))

	assert.NoError(t, err)
	assert.YAMLEq(t, expected, string(result))
}

func Test_migrating_multiple_documents(t *testing.T) {
	input := testInput(
		`` +
			`{ apiVersion: g2a-cli/v1beta4, kind: Service, name: test }` +
			"\n---\n" +
			`{ apiVersion: g2a-cli/v1beta4, kind: Environment, name: test }`,
	)
	expected := []string{
		testInput(`{ apiVersion: g2a-cli/v2.0, kind: Service, name: test }`),
		testInput(`{ apiVersion: g2a-cli/v2.0, kind: Environment, name: test }`),
	}

	migrator := NewMigrator("g2a-cli/v2.0")
	result, err := migrator.Migrate([]byte(input))
	documents := strings.Split(string(result), "\n---\n")

	assert.NoError(t, err)
	assert.Len(t, documents, len(expected))
	for i, subDocument := range documents {
		assert.YAMLEq(t, expected[i], subDocument)
	}
}

func Test_migrating_service_from_v1beta4_to_v2_0_replaces_legacy_placeholders_with_new_ones(t *testing.T) {
	cases := [][2]string{
		{".Dirs.Project", ".Project.Dir"},
		{".Dirs.Environment", ".Environment.Dir"},
		{".Dirs.Service", ".Service.Dir"},
		{".Env.FooBar", ".Environment.Vars.FooBar"},
		{".Opts.Tag", ".Tag"},
		{".Invalid", ".Invalid"},
		{".Project.Dir", ".Project.Dir"},
	}
	for _, c := range cases {
		t.Run(c[0], func(t *testing.T) {
			input := testInput(fmt.Sprintf(`{
			apiVersion: g2a-cli/v1beta4, kind: Service, name: test,
			deploy: {
				releases: [{
					deployer: { foo: "{{ %s }}" }
				}]
			}
		}`, c[0]))
			expected := testInput(fmt.Sprintf(`{
			apiVersion: g2a-cli/v2.0, kind: Service, name: test,
			releases: [{
				deployer: { foo: "{{ %s }}" }
			}]
		}`, c[1]))

			migrator := NewMigrator("g2a-cli/v2.0")
			result, err := migrator.Migrate([]byte(input))

			assert.NoError(t, err)
			assert.YAMLEq(t, expected, string(result))
		})
	}
}

func Test_migrating_environment_from_v1beta4_to_v2_0_replaces_legacy_placeholders_with_new_ones(t *testing.T) {
	cases := [][2]string{
		{".Dirs.Project", ".Project.Dir"},
		{".Dirs.Environment", ".Environment.Dir"},
		{".Dirs.Service", ".Service.Dir"},
		{".Env.FooBar", ".Environment.Vars.FooBar"},
		{".Opts.Tag", ".Tag"},
		{".Invalid", ".Invalid"},
		{".Project.Dir", ".Project.Dir"},
	}
	for _, c := range cases {
		t.Run(c[0], func(t *testing.T) {
			input := testInput(fmt.Sprintf(`{
			apiVersion: g2a-cli/v1beta4,
			kind: Environment,
			name: test,
			variables: {
				name: "{{ %s }}"
			}
		}`, c[0]))
			expected := testInput(fmt.Sprintf(`{
			apiVersion: g2a-cli/v2.0,
			kind: Environment,
			name: test,
			variables: {
				name: "{{ %s }}"
			}
		}`, c[1]))

			migrator := NewMigrator("g2a-cli/v2.0")
			result, err := migrator.Migrate([]byte(input))

			assert.NoError(t, err)
			assert.YAMLEq(t, expected, string(result))
		})
	}
}

// testInput validates input against schema and returns it back. Use only in
// tests.
func testInput(input string) string {
	_, err := Validate([]byte(input))
	if err != nil {
		panic(err)
	}
	return input
}
