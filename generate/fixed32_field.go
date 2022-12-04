package generate

import (
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfFixed32FieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling fixed32 field")
	fieldSo := &schemaObject{
		Type:    schemaType{"integer"},
		Format:  "int64",
		Default: []byte("0"),
	}
	applyUInt32Gte(bc, fieldSo, 0)
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetFixed32()
	if r == nil {
		bc.Debugf("fixed32 rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %v", *r.Const)
		fieldSo.XConst = []byte(strconv.FormatInt(int64(*r.Const), 10))
	}
	if r.Lt != nil || r.Lte != nil {
		if r.Lt != nil && r.Lte != nil {
			bc.Debugf("both lt and lte are set")
			if *r.Lt <= *r.Lte {
				applyUInt32Lt(bc, fieldSo, *r.Lt)
			} else {
				applyUInt32Lte(bc, fieldSo, *r.Lte)
			}
		} else if r.Lt == nil {
			// r.Lte != nil
			applyUInt32Lte(bc, fieldSo, *r.Lte)
		} else {
			// r.Lte == nil && r.Lt != nil
			applyUInt32Lt(bc, fieldSo, *r.Lt)
		}
	}
	if r.Gt != nil || r.Gte != nil {
		if r.Gt != nil && r.Gte != nil {
			bc.Debugf("both gt and gte are set")
			if *r.Gt >= *r.Gte {
				applyUInt32Gt(bc, fieldSo, *r.Gt)
			} else {
				applyUInt32Gte(bc, fieldSo, *r.Gte)
			}
		} else if r.Gt == nil {
			// r.Gte != nil
			applyUInt32Gte(bc, fieldSo, *r.Gte)
		} else {
			// r.Gte == nil && r.Gt != nil
			applyUInt32Gt(bc, fieldSo, *r.Gt)
		}
	}
	for _, in := range r.In {
		fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.FormatInt(int64(in), 10)))
	}
	not := &schemaObject{}
	for _, ni := range r.NotIn {
		not.Enum = append(not.Enum, []byte(strconv.FormatInt(int64(ni), 10)))
	}
	if len(not.Enum) > 0 {
		fieldSo.Not = not
	}
	return fieldSo
}
