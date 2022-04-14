package main

import "github.com/g2a-com/cicd/internal/object"

type options struct {
	object.GenericObject

	Environment string            `flag:"environment" alias:"e" help:"Name of an environment to deploy to" required:"true"`
	Tag         string            `flag:"tag" alias:"t" help:"Tag (version) of service to deploy"`
	Force       bool              `flag:"force" help:"Force release update"`
	DryRun      bool              `flag:"dry-run" help:"Simulate a deploy"`
	Wait        int               `flag:"wait" default:"0" help:"Maximum time in seconds to wait for deploy to complete, 0 - don't wait"`
	Services    []string          `flag:"services" alias:"s" help:"List of services to deploy (overrides environment configuration)"`
	Params      map[string]string `flag:"param" help:"Parameters to use in configuration files (key=value pairs)"`
	ProjectFile string            `flag:"project-file" alias:"f" help:"Path to project file"`
	ResultFile  string            `flag:"result-file" help:"Where to write result file"`
}

func (o options) Kind() object.Kind {
	return object.OptionsKind
}

func (o options) PlaceholderValues() map[string]interface{} {
	return map[string]interface{}{
		"Params": o.Params,
		"Tag":    o.Tag,
	}
}
