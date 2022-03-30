package blueprint

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/g2a-com/cicd/internal/object"
	"github.com/g2a-com/cicd/internal/placeholders"
	log "github.com/g2a-com/klio-logger-go"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

type Mode string

const (
	BuildMode  Mode = "build"
	DeployMode Mode = "deploy"
	RunMode    Mode = "run"
)

type Preprocessor func([]byte) ([]byte, error)

type Blueprint struct {
	Mode           Mode
	Services       []string
	Params         map[string]string
	Environment    string
	Tag            string
	Preprocessors  []Preprocessor
	objects        map[string]object.Object
	processedFiles map[string]bool
}

func (b *Blueprint) init() error {
	if b.processedFiles == nil {
		b.processedFiles = map[string]bool{}
	}
	if b.objects == nil {
		b.objects = map[string]object.Object{}
	}

	if b.Mode == "" {
		return errors.New("mode is not specified")
	}
	if b.Mode == DeployMode && b.Environment == "" {
		return errors.New("environment is requited in deploy mode")
	}
	return nil
}

func (b *Blueprint) Validate() (err error) {
	for _, obj := range b.objects {
		e := obj.Validate(b)
		if e != nil {
			err = multierror.Append(err, e)
		}
	}

	if b.Environment != "" {
		if b.GetObject(object.EnvironmentKind, b.Environment) == nil {
			err = multierror.Append(err, fmt.Errorf("environment %q does not exist, available environments: %s", b.Environment, strings.Join(b.getEnvironmentNames(), ", ")))
		}
	}

	for _, name := range b.getServiceNames() {
		if b.GetObject(object.ServiceKind, name) == nil {
			err = multierror.Append(err, fmt.Errorf("service %q does not exist, available services: %s", name, strings.Join(b.getServiceNames(), ", ")))
		}
	}

	return err
}

func (b *Blueprint) Load(glob string) error {
	err := b.init()
	if err != nil {
		return err
	}

	glob, err = filepath.Abs(glob)
	if err != nil {
		return err
	}

	globs := []string{glob}

	for i := 0; i < len(globs); i++ {
		glob := globs[i]

		paths, err := filepath.Glob(glob)
		if err != nil {
			return err
		}

		for _, p := range paths {
			if _, ok := b.processedFiles[p]; ok {
				continue
			} else {
				b.processedFiles[p] = true
			}

			docs, err := b.readFile(p, b.Mode)
			if err != nil {
				return fmt.Errorf(`file "%s" contains invalid document: %s`, p, err)
			}

			for _, obj := range docs {
				project, ok := obj.(object.Project)
				if ok {
					for _, entry := range project.Files {
						globs = append(globs, path.Join(project.Directory(), entry))
					}
				}
			}

			err = b.addDocuments(docs...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetProject gets Project object
func (b *Blueprint) GetProject() object.Project {
	return b.GetObjectsByKind(object.ProjectKind)[0].(object.Project)
}

// GetEnvironment gets Executor object by kind and name
func (b *Blueprint) GetExecutor(kind object.Kind, name string) (object.Executor, bool) {
	o := b.GetObject(kind, name)
	e, ok := o.(object.Executor)
	return e, o != nil && ok
}

// GetEnvironment gets Service object by the name
func (b *Blueprint) GetService(name string) (object.Service, bool) {
	o := b.GetObject(object.ServiceKind, name)
	return o.(object.Service), o != nil
}

// GetEnvironment gets Environment object by the name
func (b *Blueprint) GetEnvironment(name string) (object.Environment, bool) {
	o := b.GetObject(object.EnvironmentKind, name)
	return o.(object.Environment), o != nil
}

// ListServices returns all service objects in the blueprint
func (b *Blueprint) ListServices() []object.Service {
	names := b.getServiceNames()
	services := make([]object.Service, 0, len(names))
	for _, name := range b.getServiceNames() {
		s, _ := b.GetService(name)
		services = append(services, s)
	}
	return services
}

func (b *Blueprint) addDocuments(documents ...object.Object) error {
	for _, obj := range documents {
		key := string(obj.Kind()) + "/" + obj.Name()
		duplicate, ok := b.objects[key]
		if ok {
			return fmt.Errorf("%s is duplicated, it's defined in:\n\t* %s\n\t* %s", obj.DisplayName(), duplicate.Metadata(), obj.Metadata())
		}
		b.objects[key] = obj
	}

	return nil
}

func (b *Blueprint) ExpandPlaceholders() (err error) {
	project := b.GetProject()

	for _, obj := range b.objects {
		service, ok := obj.(object.Service)
		if !ok {
			continue
		}

		values := map[string]interface{}{
			"Project": map[string]interface{}{
				"Dir":  project.Directory(),
				"Name": project.Name(),
				"Vars": project.Variables,
			},
			"Service": map[string]interface{}{
				"Dir":  service.Directory(),
				"Name": service.Name(),
			},
			"Params": b.Params,
		}

		if b.Environment != "" {
			environment, _ := b.GetEnvironment(b.Environment)
			values["Environment"] = map[string]interface{}{
				"Dir":  environment.Directory(),
				"Name": environment.Name(),
				"Vars": environment.Variables,
			}
		}

		if b.Tag != "" {
			values["Tag"] = b.Tag
		}

		for i := range service.Build.Artifacts.ToBuild {
			service.Build.Artifacts.ToBuild[i].Spec, err = placeholders.ReplaceWithValues(service.Build.Artifacts.ToBuild[i].Spec, values)
			if err != nil {
				return err
			}
		}
		for i := range service.Build.Artifacts.ToPush {
			service.Build.Artifacts.ToPush[i].Spec, err = placeholders.ReplaceWithValues(service.Build.Artifacts.ToPush[i].Spec, values)
			if err != nil {
				return err
			}
		}
		for i := range service.Deploy.Releases {
			service.Deploy.Releases[i].Spec, err = placeholders.ReplaceWithValues(service.Deploy.Releases[i].Spec, values)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Blueprint) GetObject(kind object.Kind, name string) object.Object {
	key := string(kind) + "/" + name
	obj, ok := b.objects[key]
	if !ok {
		return nil
	}
	return obj
}

func (b *Blueprint) GetObjectsByKind(kind object.Kind) []object.Object {
	keys := make([]string, 0, len(b.objects))
	for key, obj := range b.objects {
		if obj.Kind() == kind {
			keys = append(keys, key)
		}
	}
	keys = sort.StringSlice(keys)

	objects := make([]object.Object, 0, len(keys))
	for _, key := range keys {
		objects = append(objects, b.objects[key])
	}

	return objects
}

func (b *Blueprint) getServiceNames() (names []string) {
	if len(b.Services) > 0 {
		return b.Services
	}
	for _, obj := range b.objects {
		if obj.Kind() == object.ServiceKind {
			names = append(names, obj.Name())
		}
	}
	sort.Strings(names)
	return names
}

func (b *Blueprint) getEnvironmentNames() (names []string) {
	for _, obj := range b.objects {
		if obj.Kind() == object.EnvironmentKind {
			names = append(names, obj.Name())
		}
	}
	sort.Strings(names)
	return names
}

func (b *Blueprint) readFile(filename string, mode Mode) ([]object.Object, error) {
	log.Debugf("Loading file: %s", filename)

	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	for _, preprocessor := range b.Preprocessors {
		buf, err = preprocessor(buf)
		if err != nil {
			return nil, fmt.Errorf(`file "%s" contains invalid document: %s`, filename, err)
		}
	}

	var documents []object.Object

	reader := bytes.NewReader(buf)
	decoder := yaml.NewDecoder(reader)

	for i := 0; true; i++ {
		var content yaml.Node

		err := decoder.Decode(&content)
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf(`file "%s" contains invalid document: %s`, filename, err)
			}
			break
		}

		doc, err := object.NewObject(filename, &content)
		if err != nil {
			return nil, err
		}

		if service, ok := doc.(object.Service); ok {
			if mode != BuildMode {
				service.Build.Artifacts.ToBuild = []object.ServiceEntry{}
				service.Build.Artifacts.ToPush = []object.ServiceEntry{}
				service.Build.Tags = []object.ServiceEntry{}
			}
			if mode != DeployMode {
				service.Deploy.Releases = []object.ServiceEntry{}
			}
			doc = service
		}

		if mode != DeployMode && doc.Kind() == object.EnvironmentKind {
			continue
		}

		documents = append(documents, doc)
	}

	return documents, nil
}
