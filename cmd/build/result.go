package main

import (
	"github.com/g2a-com/cicd/internal/component"
)

type ResultEntry struct {
	Service string `json:"service"`
	Entry   int    `json:"entry"`
	Result  string `json:"result"`
}

type Result struct {
	Tags            []ResultEntry `json:"tags"`
	Artifacts       []ResultEntry `json:"artifacts"`
	PushedArtifacts []ResultEntry `json:"pushedArtifacts"`
}

func (r *Result) getTags(service component.Component) (tags []string) {
	for _, r := range r.Tags {
		if r.Service == service.Name() {
			tags = append(tags, r.Result)
		}
	}
	return tags
}

func (r *Result) addTags(service component.Component, entry component.Entry, tags []string) {
	for _, tag := range tags {
		r.Tags = append(r.Tags, ResultEntry{service.Name(), entry.Index(), tag})
	}
}

func (r *Result) getArtifacts(service component.Component, entry component.Entry) (artifacts []string) {
	for _, r := range r.Artifacts {
		if r.Service == service.Name() && r.Entry == entry.Index() {
			artifacts = append(artifacts, r.Result)
		}
	}
	return artifacts
}

func (r *Result) addArtifacts(service component.Component, entry component.Entry, artifacts []string) {
	for _, artifact := range artifacts {
		r.Artifacts = append(r.Artifacts, ResultEntry{service.Name(), entry.Index(), artifact})
	}
}

func (r *Result) addPushedArtifacts(service component.Component, entry component.Entry, artifacts []string) {
	for _, artifact := range artifacts {
		r.PushedArtifacts = append(r.PushedArtifacts, ResultEntry{service.Name(), entry.Index(), artifact})
	}
}
