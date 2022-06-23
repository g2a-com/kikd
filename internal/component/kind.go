package component

type Kind string

const (
	BuilderKind     Kind = "builder"
	DeployerKind    Kind = "deployer"
	EnvironmentKind Kind = "environment"
	ProjectKind     Kind = "project"
	PusherKind      Kind = "pusher"
	ServiceKind     Kind = "service"
	TaggerKind      Kind = "tagger"
	OptionsKind     Kind = "options"
)
