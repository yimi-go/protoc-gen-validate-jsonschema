package generate

import (
	pgs "github.com/lyft/protoc-gen-star"
)

func defOfMsg(bc pgs.BuildContext, sf *SchemaFile, msg pgs.Message) string {
	sn := schemaName(msg)
	if _, ok := sf.Definitions[sn]; !ok {
		sf.Definitions[sn] = schemaOfMsg(bc, sf, msg)
	}
	return "#/definitions/" + sn
}

func schemaOfMsg(bc pgs.BuildContext, sf *SchemaFile, msg pgs.Message) *schemaObject {
	so := &schemaObject{
		Type:        schemaType{"object"},
		Description: description(msg),
		Properties:  make(map[string]*schemaObject),
	}
	applyOneOfOptions(so, msg)
	for _, field := range msg.Fields() {
		bc.Push(field.Name().String())

		defOfField(bc, sf, so, field)

		bc.Pop()
	}
	return so
}
