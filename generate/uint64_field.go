package generate

import (
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfUInt64FieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling uint64 field")
	fieldSo := &schemaObject{
		Type:    schemaType{"integer", "string"},
		Format:  "uint64",
		Default: []byte("0"),
	}
	applyUInt64Gte(bc, fieldSo, 0)
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetUint64()
	if r == nil {
		bc.Debugf("uint64 rule not found")
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
				applyUInt64Lt(bc, fieldSo, *r.Lt)
			} else {
				applyUInt64Lte(bc, fieldSo, *r.Lte)
			}
		} else if r.Lt == nil {
			// r.Lte != nil
			applyUInt64Lte(bc, fieldSo, *r.Lte)
		} else {
			// r.Lte == nil && r.Lt != nil
			applyUInt64Lt(bc, fieldSo, *r.Lt)
		}
	}
	if r.Gt != nil || r.Gte != nil {
		if r.Gt != nil && r.Gte != nil {
			bc.Debugf("both gt and gte are set")
			if *r.Gt >= *r.Gte {
				applyUInt64Gt(bc, fieldSo, *r.Gt)
			} else {
				applyUInt64Gte(bc, fieldSo, *r.Gte)
			}
		} else if r.Gt == nil {
			// r.Gte != nil
			applyUInt64Gte(bc, fieldSo, *r.Gte)
		} else {
			// r.Gte == nil && r.Gt != nil
			applyUInt64Gt(bc, fieldSo, *r.Gt)
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

func applyUInt64Lt(bc pgs.BuildContext, so *schemaObject, lt uint64) {
	bc.Debugf("apply maximum and exclusiveMaximum from lt: %d", lt)
	f64m := float64(lt)
	so.Maximum = &f64m
	so.ExclusiveMaximum = true
}

func applyUInt64Lte(bc pgs.BuildContext, so *schemaObject, lte uint64) {
	bc.Debugf("apply maximum from lte: %d", lte)
	f64m := float64(lte)
	so.Maximum = &f64m
}

func applyUInt64Gt(bc pgs.BuildContext, so *schemaObject, gt uint64) {
	bc.Debugf("apply minimum and exclusiveMinimum from lt: %d", gt)
	f64m := float64(gt)
	so.Minimum = &f64m
	so.ExclusiveMinimum = true
}

func applyUInt64Gte(bc pgs.BuildContext, so *schemaObject, gte uint64) {
	bc.Debugf("apply minimum from lte: %d", gte)
	f64m := float64(gte)
	so.Minimum = &f64m
}
