package main

import (
	"fmt"
	"github.com/g2a-com/cicd/internal/component"
	"github.com/g2a-com/cicd/internal/component/executor"
	"os"
	"path/filepath"

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
		ResultFile:  "build-result.json",
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
		Mode:     BuildMode,
		Params:   opts.Params,
		Services: opts.Services,
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

	// Helper for getting executors
	getExecutor := func(kind component.Kind, name string) executor.Executor {
		e, ok := blueprint.GetExecutor(kind, name)
		assert(ok, fmt.Errorf("%s %q does not exist", kind, name))
		return e
	}

	// Build
	for _, service := range blueprint.ListServices() {
		l := l.WithTags(service.Name())

		if len(service.Entries(component.BuildEntryType)) == 0 {
			l.WithLevel(log.VerboseLevel).Print("No artifacts to build")
			continue
		}

		// Generate tags
		for _, entry := range service.Entries(component.TagEntryType) {
			s := script.New(getExecutor(entry.ExecutorKind(), entry.ExecutorName()))
			s.Logger = l

			res, err := s.Run(TaggerInput{
				Spec: entry.Spec(&blueprint),
				Dirs: Dirs{
					Project: blueprint.GetProject().Directory(),
					Service: service.Directory(),
				},
			})
			assert(err == nil, err)

			result.addTags(service, entry, res)
		}

		if len(result.getTags(service)) == 0 {
			l.WithLevel(log.WarnLevel).Print("No tags to build")
			continue
		}

		// Build artifacts
		for _, entry := range service.Entries(component.BuildEntryType) {
			s := script.New(getExecutor(entry.ExecutorKind(), entry.ExecutorName()))
			s.Logger = l

			res, err := s.Run(BuilderInput{
				Spec: entry.Spec(&blueprint),
				Tags: result.getTags(service),
				Dirs: Dirs{
					Project: blueprint.GetProject().Directory(),
					Service: service.Directory(),
				},
			})
			assert(err == nil, err)

			result.addArtifacts(service, entry, res)
		}
	}

	// Push artifacts
	if opts.Push {
		for _, service := range blueprint.ListServices() {
			l := l.WithTags("push", service.Name())

			for _, entry := range service.Entries(component.PushEntryType) {
				s := script.New(getExecutor(entry.ExecutorKind(), entry.ExecutorName()))
				s.Logger = l

				res, err := s.Run(PusherInput{
					Spec:      entry.Spec(&blueprint),
					Tags:      result.getTags(service),
					Artifacts: result.getArtifacts(service, entry),
					Dirs: Dirs{
						Project: blueprint.GetProject().Directory(),
						Service: service.Directory(),
					},
				})
				assert(err == nil, err)

				result.addPushedArtifacts(service, entry, res)
			}
		}
	}

	// Print success message
	switch count := len(blueprint.ListServices()); count {
	case 0:
		l.Print("There was nothing to build")
	case 1:
		l.Print("Successfully built 1 service")
	default:
		l.Printf("Successfully built %v services", count)
	}
}

func assert(condition bool, err interface{}) {
	if !condition {
		panic(err)
	}
}
