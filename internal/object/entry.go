package object

const (
	TagEntryType    = "tag"
	BuildEntryType  = "build"
	PushEntryType   = "push"
	DeployEntryType = "deploy"
)

type Entry interface {
	Index() int
	ExecutorKind() Kind
	ExecutorName() string
	Spec(ObjectCollection) interface{}
}
