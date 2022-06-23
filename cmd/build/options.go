package main

import (
	"github.com/g2a-com/cicd/internal/component"
)

type options struct {
	component.Backbone

	Push        bool              `flag:"push" alias:"p" help:"Push artifacts to remote registry"`
	Services    []string          `flag:"services" alias:"s" help:"List of services to build (skip to build all services)"`
	Params      map[string]string `flag:"param" help:"Parameters to use in configuration files (key=value pairs)"`
	ProjectFile string            `flag:"project-file" alias:"f" help:"Path to project file"`
	ResultFile  string            `flag:"result-file" help:"Where to write result file"`
}

func (o options) Kind() component.Kind {
	return component.OptionsKind
}

func (o options) PlaceholderValues() map[string]interface{} {
	return map[string]interface{}{
		"Params": o.Params,
	}
}
