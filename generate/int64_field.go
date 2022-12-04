package generate

import (
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfInt64FieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling int64 field")
	fieldSo := &schemaObject{
		Type:    schemaType{"integer", "string"},
		Format:  "int64",
		Default: []byte("0"),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetInt64()
	if r == nil {
		bc.Debugf("int64 rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %v", *r.Const)
		fieldSo.XConst = []byte(strconv.FormatInt(*r.Const, 10))
	}
	if r.Lt != nil || r.Lte != nil {
		if r.Lt != nil && r.Lte != nil {
			bc.Debugf("both lt and lte are set")
			if *r.Lt <= *r.Lte {
				applyInt64Lt(bc, fieldSo, *r.Lt)
			} else {
				applyInt64Lte(bc, fieldSo, *r.Lte)
			}
		} else if r.Lt == nil {
			// r.Lte != nil
			applyInt64Lte(bc, fieldSo, *r.Lte)
		} else {
			// r.Lte == nil && r.Lt != nil
			applyInt64Lt(bc, fieldSo, *r.Lt)
		}
	}
	if r.Gt != nil || r.Gte != nil {
		if r.Gt != nil && r.Gte != nil {
			bc.Debugf("both gt and gte are set")
			if *r.Gt >= *r.Gte {
				applyInt64Gt(bc, fieldSo, *r.Gt)
			} else {
				applyInt64Gte(bc, fieldSo, *r.Gte)
			}
		} else if r.Gt == nil {
			// r.Gte != nil
			applyInt64Gte(bc, fieldSo, *r.Gte)
		} else {
			// r.Gte == nil && r.Gt != nil
			applyInt64Gt(bc, fieldSo, *r.Gt)
		}
	}
	for _, in := range r.In {
		fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.FormatInt(in, 10)))
	}
	not := &schemaObject{}
	for _, ni := range r.NotIn {
		not.Enum = append(not.Enum, []byte(strconv.FormatInt(ni, 10)))
	}
	if len(not.Enum) > 0 {
		fieldSo.Not = not
	}
	return fieldSo
}

func applyInt64Lt(bc pgs.BuildContext, so *schemaObject, lt int64) {
	bc.Debugf("apply maximum and exclusiveMaximum from lt: %d", lt)
	f64m := float64(lt)
	so.Maximum = &f64m
	so.ExclusiveMaximum = true
}

func applyInt64Lte(bc pgs.BuildContext, so *schemaObject, lte int64) {
	bc.Debugf("apply maximum from lte: %d", lte)
	f64m := float64(lte)
	so.Maximum = &f64m
}

func applyInt64Gt(bc pgs.BuildContext, so *schemaObject, gt int64) {
	bc.Debugf("apply minimum and exclusiveMinimum from lt: %d", gt)
	f64m := float64(gt)
	so.Minimum = &f64m
	so.ExclusiveMinimum = true
}

func applyInt64Gte(bc pgs.BuildContext, so *schemaObject, gte int64) {
	bc.Debugf("apply minimum from lte: %d", gte)
	f64m := float64(gte)
	so.Minimum = &f64m
}
