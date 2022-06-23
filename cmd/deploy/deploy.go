package main

import (
	"fmt"
	"github.com/g2a-com/cicd/internal/component"
	"os"
	"path/filepath"
	"strings"

	. "github.com/g2a-com/cicd/internal/blueprint"
	"github.com/g2a-com/cicd/internal/flags"
	"github.com/g2a-com/cicd/internal/schema"
	"github.com/g2a-com/cicd/internal/script"
	"github.com/g2a-com/cicd/internal/utils"
	log "github.com/g2a-com/klio-logger-go/v2"
)

func main() {
	var err error

	// Exit nicely on panics
	defer utils.HandlePanics()

	// Parse options
	opts := options{
		ResultFile:  "deploy-result.json",
		ProjectFile: utils.FindProjectFile(),
	}
	flags.ParseArgs(&opts, os.Args)

	// Prepare logger
	l := log.StandardLogger()

	// Handle results
	result := &Result{}
	defer utils.SaveResult(opts.ResultFile, result)

	// Check if project file exists
	if !utils.FileExists(opts.ProjectFile) {
		panic("cannot find project.yaml")
	}

	// Load blueprint
	blueprint := Blueprint{
		Mode:        DeployMode,
		Environment: opts.Environment,
		Tag:         opts.Tag,
		Params:      opts.Params,
		Services:    opts.Services,
		Preprocessors: []Preprocessor{
			schema.Validate,
			schema.Migrate,
		},
	}
	err = blueprint.Load(filepath.Join(utils.FindCommandDirectory(), "assets", "executors", "*", "*.yaml"))
	assert(err == nil, err)
	err = blueprint.Load(opts.ProjectFile)
	assert(err == nil, err)
	err = blueprint.AddDocuments(opts)
	assert(err == nil, err)
	err = blueprint.Validate()
	assert(err == nil, err)

	// Change working directory
	err = os.Chdir(blueprint.GetProject().Directory())
	assert(err == nil, err)

	// Deploy
	l.Printf(`Deploying to environment %q...`, opts.Environment)

	environment, _ := blueprint.GetEnvironment(opts.Environment)

	for _, service := range blueprint.ListServices() {
		l := l.WithTags(service.Name())

		if len(service.Entries(component.DeployEntryType)) == 0 {
			l.WithLevel(log.VerboseLevel).Print("No releases to deploy")
			continue
		}

		l.Printf(`Deploying service %q...`, service.Name())

		for _, entry := range service.Entries(component.DeployEntryType) {
			e, ok := blueprint.GetExecutor(entry.ExecutorKind(), entry.ExecutorName())
			assert(ok, fmt.Errorf("%s %q does not exist", strings.ToLower(string(entry.ExecutorKind())), entry.ExecutorName()))

			s := script.New(e)
			s.Logger = l

			res, err := s.Run(DeployerInput{
				Spec:   entry.Spec(&blueprint),
				Force:  opts.Force,
				DryRun: opts.DryRun,
				Wait:   opts.Wait,
				Dirs: Dirs{
					Project:     blueprint.GetProject().Directory(),
					Environment: environment.Directory(),
					Service:     service.Directory(),
				},
			})
			assert(err == nil, err)

			result.addReleases(service, entry, res)
		}
	}

	// Print success message
	switch count := len(blueprint.ListServices()); count {
	case 0:
		l.Printf("There was nothing to deploy to environment %q", opts.Environment)
	case 1:
		l.Printf("Successfully deployed 1 service to environment %q", opts.Environment)
	default:
		l.Printf("Successfully deployed %v services to environment %q", count, opts.Environment)
	}
}

func assert(condition bool, err interface{}) {
	if !condition {
		panic(err)
	}
}
