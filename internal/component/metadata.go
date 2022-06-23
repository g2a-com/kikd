package component

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Metadata interface {
	Filename() string
	Line() int
	String() string
}

type metadata struct {
	filename string `yaml:"-"`
	line     int    `yaml:"-"`
}

func NewMetadata(filename string, data *yaml.Node) Metadata {
	return metadata{filename, data.Line}
}

func (m metadata) String() string {
	if m.line != 0 {
		return fmt.Sprintf("%s:%d", m.filename, m.line)
	} else {
		return m.filename
	}
}

func (m metadata) Filename() string {
	return m.filename
}

func (m metadata) Line() int {
	return m.line
}
