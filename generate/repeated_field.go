package generate

import (
	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfRepeatedFieldType(
	bc pgs.BuildContext, sf *SchemaFile, ft pgs.FieldType, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling repeated field")
	fieldSo := &schemaObject{
		Type:  schemaType{"array"},
		Items: itemScheme(bc, sf, ft.Element(), itemRules(rules)),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetRepeated()
	if r == nil {
		bc.Debugf("repeated rule not found")
		return fieldSo
	}
	if r.MinItems != nil {
		bc.Debugf("apply minItems from min_items: %d", *r.MinItems)
		fieldSo.MinItems = *r.MinItems
	}
	if r.MaxItems != nil {
		bc.Debugf("apply maxItems from max_items: %d", *r.MaxItems)
		fieldSo.MaxItems = *r.MaxItems
	}
	if r.Unique != nil {
		bc.Debugf("apply uniqueItems from unique: %t", *r.Unique)
		fieldSo.UniqueItems = *r.Unique
	}
	return fieldSo
}

func itemRules(rules *validate.FieldRules) *validate.FieldRules {
	if rules == nil {
		return nil
	}
	if rules.GetRepeated() == nil {
		return nil
	}
	return rules.GetRepeated().Items
}

func itemScheme(
	bc pgs.BuildContext, sf *SchemaFile, fte pgs.FieldTypeElem, itemRules *validate.FieldRules) *schemaObject {
	bc.Push("repeated items")
	defer bc.Pop()
	if fte.IsEnum() {
		return schemaOfEnumFieldType(bc, fte.Enum(), itemRules)
	}
	if fte.IsEmbed() {
		ref := defOfMsg(bc, sf, fte.Embed())
		return &schemaObject{Ref: ref}
	}
	// fte must be scalar type.
	return schemaOfScalarFieldType(bc, fte.ProtoType(), itemRules)
}
