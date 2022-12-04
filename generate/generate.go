package generate

import (
	pgs "github.com/lyft/protoc-gen-star"
)

func Generate(bc pgs.BuildContext, msg pgs.Message) *SchemaFile {
	sf := &SchemaFile{
		Schema:      "http://json-schema.org/draft-04/schema#",
		Definitions: map[string]*schemaObject{},
	}
	sf.Ref = defOfMsg(bc, sf, msg)
	return sf
}
