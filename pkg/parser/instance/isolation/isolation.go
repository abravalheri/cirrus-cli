package isolation

import (
	"github.com/cirruslabs/cirrus-ci-agent/api"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/nameable"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/node"
	"github.com/cirruslabs/cirrus-cli/pkg/parser/parseable"
	jsschema "github.com/lestrrat-go/jsschema"
)

type Isolation struct {
	proto api.Isolation

	parseable.DefaultParser
}

func NewIsolation(mergedEnv map[string]string) *Isolation {
	isolation := &Isolation{}

	parallelsSchema := NewParallels(mergedEnv).Schema()
	isolation.OptionalField(nameable.NewSimpleNameable("parallels"), parallelsSchema, func(node *node.Node) error {
		parallels := NewParallels(mergedEnv)

		if err := parallels.Parse(node); err != nil {
			return err
		}

		isolation.proto.Type = parallels.Proto()

		return nil
	})

	containerSchema := NewContainer(mergedEnv).Schema()
	isolation.OptionalField(nameable.NewSimpleNameable("container"), containerSchema, func(node *node.Node) error {
		container := NewContainer(mergedEnv)

		if err := container.Parse(node); err != nil {
			return err
		}

		isolation.proto.Type = container.Proto()

		return nil
	})

	return isolation
}

func (isolation *Isolation) Parse(node *node.Node) error {
	return isolation.DefaultParser.Parse(node)
}

func (isolation *Isolation) Proto() *api.Isolation {
	return &isolation.proto
}

func (isolation *Isolation) Schema() *jsschema.Schema {
	modifiedSchema := isolation.DefaultParser.Schema()

	modifiedSchema.Type = jsschema.PrimitiveTypes{jsschema.ObjectType}
	modifiedSchema.Description = "Persistent Worker isolation."

	return modifiedSchema
}
