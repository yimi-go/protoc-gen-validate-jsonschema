package generate

import (
	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func applyOneOfOptions(so *schemaObject, msg pgs.Message) {
	for _, oneOf := range msg.RealOneOfs() {
		oo := &OneOf{Name: oneOf.Name().LowerCamelCase().String()}
		for _, field := range oneOf.Fields() {
			oo.Fields = append(oo.Fields, field.Name().LowerCamelCase().String())
		}
		_, _ = oneOf.Extension(validate.E_Required, &oo.Required)
		so.XOneOfs = append(so.XOneOfs, oo)
	}
}
