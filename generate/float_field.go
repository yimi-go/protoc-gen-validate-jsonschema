package generate

import (
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfFloatFieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling float field")
	fieldSo := &schemaObject{
		Type:    schemaType{"number"},
		Format:  "float",
		Default: []byte("0"),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetFloat()
	if r == nil {
		bc.Debugf("float rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %v", *r.Const)
		fieldSo.XConst = []byte(strconv.FormatFloat(float64(*r.Const), 'f', -1, 64))
	}
	if r.Lt != nil || r.Lte != nil {
		if r.Lt != nil && r.Lte != nil {
			bc.Debugf("both lt and lte are set")
			if *r.Lt <= *r.Lte {
				applyFloatLt(bc, fieldSo, *r.Lt)
			} else {
				applyFloatLte(bc, fieldSo, *r.Lte)
			}
		} else if r.Lt == nil {
			// r.Lte != nil
			applyFloatLte(bc, fieldSo, *r.Lte)
		} else {
			// r.Lte == nil && r.Lt != nil
			applyFloatLt(bc, fieldSo, *r.Lt)
		}
	}
	if r.Gt != nil || r.Gte != nil {
		if r.Gt != nil && r.Gte != nil {
			bc.Debugf("both gt and gte are set")
			if *r.Gt >= *r.Gte {
				applyFloatGt(bc, fieldSo, *r.Gt)
			} else {
				applyFloatGte(bc, fieldSo, *r.Gte)
			}
		} else if r.Gt == nil {
			// r.Gte != nil
			applyFloatGte(bc, fieldSo, *r.Gte)
		} else {
			// r.Gte == nil && r.Gt != nil
			applyFloatGt(bc, fieldSo, *r.Gt)
		}
	}
	for _, in := range r.In {
		fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.FormatFloat(float64(in), 'f', -1, 64)))
	}
	not := &schemaObject{}
	for _, ni := range r.NotIn {
		not.Enum = append(not.Enum, []byte(strconv.FormatFloat(float64(ni), 'f', -1, 64)))
	}
	if len(not.Enum) > 0 {
		fieldSo.Not = not
	}
	return fieldSo
}

func applyFloatLt(bc pgs.BuildContext, so *schemaObject, lt float32) {
	bc.Debugf("apply maximum and exclusiveMaximum from lt: %f", lt)
	f64Lt := float64(lt)
	so.Maximum = &f64Lt
	so.ExclusiveMaximum = true
}

func applyFloatLte(bc pgs.BuildContext, so *schemaObject, lte float32) {
	bc.Debugf("apply maximum from lte: %f", lte)
	f64Lte := float64(lte)
	so.Maximum = &f64Lte
}

func applyFloatGt(bc pgs.BuildContext, so *schemaObject, gt float32) {
	bc.Debugf("apply minimum and exclusiveMinimum from lt: %f", gt)
	f64Gt := float64(gt)
	so.Minimum = &f64Gt
	so.ExclusiveMinimum = true
}

func applyFloatGte(bc pgs.BuildContext, so *schemaObject, gte float32) {
	bc.Debugf("apply minimum from lte: %f", gte)
	f64Gte := float64(gte)
	so.Minimum = &f64Gte
}
