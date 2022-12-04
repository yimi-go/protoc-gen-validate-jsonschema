package generate

import (
	"fmt"
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfEnumFieldType(bc pgs.BuildContext, et pgs.Enum, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling enum field")
	fieldSo := &schemaObject{
		Type:    schemaType{"string"},
		Default: []byte(`"0"`),
	}
	for _, value := range et.Values() {
		fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.QuoteToASCII(value.Name().String())))
		fieldSo.Enum = append(fieldSo.Enum, []byte(fmt.Sprintf(`"%d"`, value.Value())))
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetEnum()
	if r == nil {
		bc.Debugf("enum rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %v", *r.Const)
		fieldSo.XConst = []byte(fmt.Sprintf(`"%d"`, *r.Const))
	}
	values := et.Values()
	in := make(map[int32]struct{}, len(values))
	notIn := make(map[int32]struct{}, len(values))
	if len(r.In) == 0 {
		for _, value := range values {
			in[value.Value()] = struct{}{}
		}
	} else {
		for _, i := range r.In {
			in[i] = struct{}{}
		}
	}
	for _, i := range r.NotIn {
		notIn[i] = struct{}{}
	}
	if len(in) == len(values) && len(notIn) == 0 {
		bc.Debugf("apply all enum values")
		return fieldSo
	}
	fieldSo.Enum = fieldSo.Enum[:0]
	for _, value := range values {
		_, isIn := in[value.Value()]
		_, isNotIn := notIn[value.Value()]
		if isIn && isNotIn {
			bc.Logf("ambiguous rule, enum value %d is both in and not_in, apply not_in rule", value.Value())
			continue
		}
		if isIn {
			bc.Debugf("accept enum value %d", value.Value())
			fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.QuoteToASCII(value.Name().String())))
			fieldSo.Enum = append(fieldSo.Enum, []byte(fmt.Sprintf(`"%d"`, value.Value())))
			continue
		}
		if isNotIn {
			bc.Logf("skip enum value %d, because not_in rule", value.Value())
			continue
		}
		bc.Logf("skip enum value %d, because in rule not contains it.", value.Value())
	}
	return fieldSo
}
