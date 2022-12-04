package generate

import (
	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func defOfField(bc pgs.BuildContext, sf *SchemaFile, msgSo *schemaObject, field pgs.Field) {
	rules := fieldRules(bc, field)
	fieldSo, required := schemaOfFieldType(bc, sf, field.Type(), rules)
	fieldSo.Description = description(field)
	pn := propertyName(field)
	msgSo.Properties[pn] = fieldSo
	if required {
		requiredProperty(msgSo, pn)
	}
	if field.Type().IsEmbed() && rules != nil && rules.Message != nil && rules.Message.GetRequired() {
		requiredProperty(msgSo, pn)
	}
}

func requiredProperty(msgSo *schemaObject, pn string) {
	hasRequired := false
	for i := range msgSo.Required {
		hasRequired = pn == msgSo.Required[i]
		if hasRequired {
			break
		}
	}
	if !hasRequired {
		msgSo.Required = append(msgSo.Required, pn)
	}
}

func fieldRules(bc pgs.BuildContext, field pgs.Field) *validate.FieldRules {
	rules := &validate.FieldRules{}
	if ok, _ := field.Extension(validate.E_Rules, rules); !ok {
		bc.Debugf("field rules not found")
		return nil
	}
	bc.Debugf("found field rules")
	return rules
}

func schemaOfFieldType(
	bc pgs.BuildContext, sf *SchemaFile, ft pgs.FieldType, rules *validate.FieldRules) (*schemaObject, bool) {
	if ft.IsEnum() {
		if ft.Enum().Package().ProtoName() == pgs.WellKnownTypePackage {
			switch ft.Enum().Name() {
			case "NullValue":
				return &schemaObject{Type: schemaType{"null"}}, false
			}
		}
		return schemaOfEnumFieldType(bc, ft.Enum(), rules), false
	}
	if ft.IsRepeated() {
		return schemaOfRepeatedFieldType(bc, sf, ft, rules), false
	}
	if ft.IsMap() {
		return schemaOfMapFieldType(bc, sf, ft, rules), false
	}
	if ft.IsEmbed() {
		if wktSo, required := wktSchema(bc, ft, rules); wktSo != nil {
			return wktSo, required
		}
		ref := defOfMsg(bc, sf, ft.Embed())
		required := false
		if rules != nil && rules.Message != nil {
			required = rules.Message.GetRequired()
		}
		return &schemaObject{Ref: ref}, required
	}

	if fieldSo := schemaOfScalarFieldType(bc, ft.ProtoType(), rules); fieldSo != nil {
		return fieldSo, false
	}
	bc.Logf("unimplemented for type %s", ft.ProtoType().String())
	return &schemaObject{}, false // return a fully free schema for unknown type field.
}

func schemaOfScalarFieldType(bc pgs.BuildContext, fpt pgs.ProtoType, rules *validate.FieldRules) *schemaObject {
	switch fpt {
	case pgs.FloatT:
		return schemaOfFloatFieldType(bc, rules)
	case pgs.DoubleT:
		return schemaOfDoubleFieldType(bc, rules)
	case pgs.Int32T:
		return schemaOfInt32FieldType(bc, rules)
	case pgs.Int64T:
		return schemaOfInt64FieldType(bc, rules)
	case pgs.UInt32T:
		return schemaOfUInt32FieldType(bc, rules)
	case pgs.UInt64T:
		return schemaOfUInt64FieldType(bc, rules)
	case pgs.SInt32:
		return schemaOfSInt32FieldType(bc, rules)
	case pgs.SInt64:
		return schemaOfSInt64FieldType(bc, rules)
	case pgs.Fixed32T:
		return schemaOfFixed32FieldType(bc, rules)
	case pgs.Fixed64T:
		return schemaOfFixed64FieldType(bc, rules)
	case pgs.SFixed32:
		return schemaOfSFixed32FieldType(bc, rules)
	case pgs.SFixed64:
		return schemaOfSFixed64FieldType(bc, rules)
	case pgs.BoolT:
		return schemaOfBoolFieldType(bc, rules)
	case pgs.StringT:
		return schemaOfStringFieldType(bc, rules)
	case pgs.BytesT:
		return schemaOfBytesFieldType(bc, rules)
	default:
		return nil
	}
}
