package generate

import (
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfDoubleFieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling double field")
	fieldSo := &schemaObject{
		Type:    schemaType{"number"},
		Format:  "double",
		Default: []byte("0"),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetDouble()
	if r == nil {
		bc.Debugf("double rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %v", *r.Const)
		fieldSo.XConst = []byte(strconv.FormatFloat(*r.Const, 'f', -1, 64))
	}
	if r.Lt != nil || r.Lte != nil {
		if r.Lt != nil && r.Lte != nil {
			bc.Debugf("both lt and lte are set")
			if *r.Lt <= *r.Lte {
				applyDoubleLt(bc, fieldSo, *r.Lt)
			} else {
				applyDoubleLte(bc, fieldSo, *r.Lte)
			}
		} else if r.Lt == nil {
			// r.Lte != nil
			applyDoubleLte(bc, fieldSo, *r.Lte)
		} else {
			// r.Lte == nil && r.Lt != nil
			applyDoubleLt(bc, fieldSo, *r.Lt)
		}
	}
	if r.Gt != nil || r.Gte != nil {
		if r.Gt != nil && r.Gte != nil {
			bc.Debugf("both gt and gte are set")
			if *r.Gt >= *r.Gte {
				applyDoubleGt(bc, fieldSo, *r.Gt)
			} else {
				applyDoubleGte(bc, fieldSo, *r.Gte)
			}
		} else if r.Gt == nil {
			// r.Gte != nil
			applyDoubleGte(bc, fieldSo, *r.Gte)
		} else {
			// r.Gte == nil && r.Gt != nil
			applyDoubleGt(bc, fieldSo, *r.Gt)
		}
	}
	for _, in := range r.In {
		fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.FormatFloat(in, 'f', -1, 64)))
	}
	not := &schemaObject{}
	for _, ni := range r.NotIn {
		not.Enum = append(not.Enum, []byte(strconv.FormatFloat(ni, 'f', -1, 64)))
	}
	if len(not.Enum) > 0 {
		fieldSo.Not = not
	}
	return fieldSo
}

func applyDoubleLt(bc pgs.BuildContext, so *schemaObject, lt float64) {
	bc.Debugf("apply maximum and exclusiveMaximum from lt: %f", lt)
	so.Maximum = &lt
	so.ExclusiveMaximum = true
}

func applyDoubleLte(bc pgs.BuildContext, so *schemaObject, lte float64) {
	bc.Debugf("apply maximum from lte: %f", lte)
	so.Maximum = &lte
}

func applyDoubleGt(bc pgs.BuildContext, so *schemaObject, gt float64) {
	bc.Debugf("apply minimum and exclusiveMinimum from lt: %f", gt)
	so.Minimum = &gt
	so.ExclusiveMinimum = true
}

func applyDoubleGte(bc pgs.BuildContext, so *schemaObject, gte float64) {
	bc.Debugf("apply minimum from lte: %f", gte)
	so.Minimum = &gte
}
