package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unmarshalling_empty_service(t *testing.T) {
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Service,
		name: test,
	}`)

	result, err := NewService("dir/file.yaml", input)

	assert.NoError(t, err)
	assert.Equal(t, "dir/file.yaml", result.Metadata().Filename())
	assert.Equal(t, ServiceKind, result.Kind())
	assert.Equal(t, "test", result.Name())
	assert.Equal(t, "dir", result.Directory())
	assert.Equal(t, `service "test"`, result.DisplayName())
}

func Test_validating_empty_service_passes(t *testing.T) {
	collection := testCollection([]Object{})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0,
		kind: Service,
		name: test,
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_using_unknown_tagger_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: TaggerKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ unknown: {} }],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_using_known_tagger_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: TaggerKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ known: {} }],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_using_unknown_builder_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			unknown: {},
			push: false,
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_using_known_builder_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			known: {},
			push: false,
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_using_unknown_pusher_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "known", schema: "{}"},
		fakeObject{kind: PusherKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			known: {},
			push: { unknown: {} }
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_using_known_pusher_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "known", schema: "{}"},
		fakeObject{kind: PusherKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			known: {},
			push: { known: {} }
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_using_unknown_deployer_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: DeployerKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [ { unknown: {} } ],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_using_known_deployer_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: DeployerKind, name: "known", schema: "{}"},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [ { known: {} } ],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_with_tags_entry_not_matching_tagger_schema_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: TaggerKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{
			type: {},
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_with_tags_entry_matching_tagger_schema_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: TaggerKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{
			type: { foo: true },
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_with_artifacts_entry_not_matching_builder_schema_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			type: {},
			push: false,
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_with_artifacts_entry_matching_builder_schema_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			type: { foo: true },
			push: false,
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_with_artifacts_entry_not_matching_pusher_schema_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "type", schema: `{}`},
		fakeObject{kind: PusherKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			type: {},
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_with_artifacts_entry_matching_pusher_schema_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: BuilderKind, name: "type", schema: `{}`},
		fakeObject{kind: PusherKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{
			type: { foo: true },
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_validating_service_with_releases_entry_not_matching_deployer_schema_fails(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: DeployerKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{
			type: {},
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.Error(t, err)
}

func Test_validating_service_with_releases_entry_matching_deployer_schema_passes(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: DeployerKind, name: "type", schema: `{ "required": ["foo"] }`},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		releases: [{
			type: { foo: true },
		}],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	err := service.Validate(collection)

	assert.NoError(t, err)
}

func Test_getting_entry_types_list_works(t *testing.T) {
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
	}`)

	service, _ := NewService("dir/file.yaml", input)
	result := service.EntryTypes()

	assert.Equal(t, []string{"tag", "build", "push", "deploy"}, result)
}

func Test_getting_tag_entries_returns_only_entries_defined_in_tags_property(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: ProjectKind},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ tag1: spec1 }, { tag2: spec2 }],
		artifacts: [{ build: spec }],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	result := service.Entries(TagEntryType)

	assert.Len(t, result, 2)
	assert.Equal(t, 0, result[0].Index())
	assert.Equal(t, TaggerKind, result[0].ExecutorKind())
	assert.Equal(t, "tag1", result[0].ExecutorName())
	assert.Equal(t, "spec1", result[0].Spec(collection))
	assert.Equal(t, 1, result[1].Index())
	assert.Equal(t, TaggerKind, result[1].ExecutorKind())
	assert.Equal(t, "tag2", result[1].ExecutorName())
	assert.Equal(t, "spec2", result[1].Spec(collection))
}

func Test_getting_build_entries_returns_only_entries_defined_in_artifacts_property(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: ProjectKind},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ tag: spec }],
		artifacts: [{ build1: spec1 }, { build2: spec2 }],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	result := service.Entries(BuildEntryType)

	assert.Len(t, result, 2)
	assert.Equal(t, 0, result[0].Index())
	assert.Equal(t, BuilderKind, result[0].ExecutorKind())
	assert.Equal(t, "build1", result[0].ExecutorName())
	assert.Equal(t, "spec1", result[0].Spec(collection))
	assert.Equal(t, 1, result[1].Index())
	assert.Equal(t, BuilderKind, result[1].ExecutorKind())
	assert.Equal(t, "build2", result[1].ExecutorName())
	assert.Equal(t, "spec2", result[1].Spec(collection))
}

func Test_getting_push_entries_returns_only_entries_defined_in_artifacts_property_preffering_definition_from_push_property_when_specified(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: ProjectKind},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ tag: spec }],
		artifacts: [
			{ build1: spec1 },
			{ build2: spec, push: { push2: spec2 } },
		],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	result := service.Entries(PushEntryType)

	assert.Len(t, result, 2)
	assert.Equal(t, 0, result[0].Index())
	assert.Equal(t, PusherKind, result[0].ExecutorKind())
	assert.Equal(t, "build1", result[0].ExecutorName())
	assert.Equal(t, "spec1", result[0].Spec(collection))
	assert.Equal(t, 1, result[1].Index())
	assert.Equal(t, PusherKind, result[1].ExecutorKind())
	assert.Equal(t, "push2", result[1].ExecutorName())
	assert.Equal(t, "spec2", result[1].Spec(collection))
}

func Test_getting_push_entries_ignores_entries_with_push_property_set_to_false(t *testing.T) {
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		tags: [{ tag: spec }],
		artifacts: [
			{ build1: spec1 },
			{ build2: spec, push: false },
		],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	result := service.Entries(PushEntryType)

	assert.Len(t, result, 1)
	assert.Equal(t, "build1", result[0].ExecutorName())
}

func Test_getting_deploy_entries_returns_only_entries_defined_in_releases_property(t *testing.T) {
	collection := testCollection([]Object{
		fakeObject{kind: ProjectKind},
	})
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
		artifacts: [{ build: spec }],
		releases: [{ release1: spec1 }, { release2: spec2 }],
	}`)

	service, _ := NewService("dir/file.yaml", input)
	result := service.Entries(DeployEntryType)

	assert.Len(t, result, 2)
	assert.Equal(t, 0, result[0].Index())
	assert.Equal(t, DeployerKind, result[0].ExecutorKind())
	assert.Equal(t, "release1", result[0].ExecutorName())
	assert.Equal(t, "spec1", result[0].Spec(collection))
	assert.Equal(t, 1, result[1].Index())
	assert.Equal(t, DeployerKind, result[1].ExecutorKind())
	assert.Equal(t, "release2", result[1].ExecutorName())
	assert.Equal(t, "spec2", result[1].Spec(collection))
}

func Test_getting_entries_of_unknown_type_returns_empty_slice(t *testing.T) {
	input := prepareTestInput(`{
		apiVersion: g2a-cli/v2.0, kind: Service, name: test,
	}`)

	service, _ := NewService("dir/file.yaml", input)
	result := service.Entries("unknown")

	assert.Empty(t, result)
}
