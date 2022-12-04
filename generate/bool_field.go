package generate

import (
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfBoolFieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling bool field")
	fieldSo := &schemaObject{
		Type:    schemaType{"boolean"},
		Default: []byte("false"),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetBool()
	if r == nil {
		bc.Debugf("bool rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %v", *r.Const)
		fieldSo.XConst = []byte(strconv.FormatBool(*r.Const))
	}
	return fieldSo
}
