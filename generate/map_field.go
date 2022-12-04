package generate

import (
	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfMapFieldType(
	bc pgs.BuildContext, sf *SchemaFile, ft pgs.FieldType, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling map field")
	fieldSo := &schemaObject{
		Type:                 schemaType{"object"},
		AdditionalProperties: itemScheme(bc, sf, ft.Element(), valueRules(rules)),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetMap()
	if r == nil {
		bc.Debugf("map rule not found")
		return fieldSo
	}
	if r.MinPairs != nil {
		bc.Debugf("apply minProperties from min_pairs: %d", *r.MinPairs)
		fieldSo.MinProperties = *r.MinPairs
	}
	if r.MaxPairs != nil {
		bc.Debugf("apply maxProperties from max_pairs: %d", *r.MaxPairs)
		fieldSo.MaxProperties = *r.MaxPairs
	}
	bc.Logf("skip checking no_sparse")
	bc.Logf("skip checking keys")
	return fieldSo
}

func valueRules(rules *validate.FieldRules) *validate.FieldRules {
	if rules == nil {
		return nil
	}
	if rules.GetMap() == nil {
		return nil
	}
	return rules.GetMap().GetValues()
}
