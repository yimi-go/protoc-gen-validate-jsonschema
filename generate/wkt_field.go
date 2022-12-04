package generate

import (
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
	"google.golang.org/protobuf/encoding/protojson"
)

func wktSchema(bc pgs.BuildContext, ft pgs.FieldType, rules *validate.FieldRules) (*schemaObject, bool) {
	if fieldSo := unwrapWkt(bc, ft.Embed().WellKnownType(), rules); fieldSo != nil {
		fieldSo.Default = nil
		return fieldSo, false
	}
	switch ft.Embed().WellKnownType() {
	case pgs.AnyWKT:
		return schemaOfWktAny(bc, rules)
	case pgs.DurationWKT:
		return schemaOfWktDuration(bc, rules)
	case pgs.TimestampWKT:
		return schemaOfWktTimestamp(bc, rules)
	case pgs.EmptyWKT:
		return &schemaObject{Type: schemaType{"object"}}, false
	case pgs.ListValueWKT:
		return &schemaObject{Type: schemaType{"array"}, Items: &schemaObject{}}, false
	case pgs.StructWKT:
		return &schemaObject{Type: schemaType{"object"}}, false
	case pgs.ValueWKT:
		return &schemaObject{}, false
	}
	if ft.Embed().Package().ProtoName() == pgs.WellKnownTypePackage {
		switch ft.Embed().Name() {
		case "FieldMask":
			return &schemaObject{Type: schemaType{"string"}}, false
		}
	}
	bc.Debugf("not well-known types")
	return nil, false
}

func unwrapWkt(
	bc pgs.BuildContext, wkt pgs.WellKnownType, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("try handle well-known type as wrapper")

	switch wkt {
	case pgs.FloatValueWKT:
		return schemaOfFloatFieldType(bc, rules)
	case pgs.DoubleValueWKT:
		return schemaOfDoubleFieldType(bc, rules)
	case pgs.Int32ValueWKT:
		return schemaOfInt32FieldType(bc, rules)
	case pgs.Int64ValueWKT:
		return schemaOfInt64FieldType(bc, rules)
	case pgs.UInt32ValueWKT:
		return schemaOfUInt32FieldType(bc, rules)
	case pgs.UInt64ValueWKT:
		return schemaOfUInt64FieldType(bc, rules)
	case pgs.BoolValueWKT:
		return schemaOfBoolFieldType(bc, rules)
	case pgs.StringValueWKT:
		return schemaOfStringFieldType(bc, rules)
	case pgs.BytesValueWKT:
		return schemaOfBytesFieldType(bc, rules)
	default:
		bc.Debugf("not wrapper")
		return nil
	}
}

func schemaOfWktAny(bc pgs.BuildContext, rules *validate.FieldRules) (*schemaObject, bool) {
	bc.Debugf("handling Any field")
	typeSo := &schemaObject{Type: schemaType{"string"}}
	fieldSo := &schemaObject{
		Type:                 schemaType{"object"},
		Required:             []string{"@type"},
		Properties:           map[string]*schemaObject{"@type": typeSo},
		AdditionalProperties: &schemaObject{},
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo, false
	}
	r := rules.GetAny()
	if r == nil {
		bc.Debugf("any rule not found")
		return fieldSo, false
	}
	if len(r.In) > 0 {
		for _, s := range r.In {
			fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.QuoteToASCII(s)))
		}
	}
	if len(r.NotIn) > 0 {
		not := &schemaObject{}
		for _, s := range r.NotIn {
			not.Enum = append(not.Enum, []byte(strconv.QuoteToASCII(s)))
		}
		fieldSo.Not = not
	}
	return fieldSo, r.GetRequired()
}

func schemaOfWktDuration(bc pgs.BuildContext, rules *validate.FieldRules) (*schemaObject, bool) {
	bc.Debugf("handling Duration field")
	fieldSo := &schemaObject{
		Type:   schemaType{"string"},
		Format: "duration",
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo, false
	}
	r := rules.GetDuration()
	if r == nil {
		bc.Debugf("duration rule not found")
		return fieldSo, false
	}
	if r.Const != nil {
		bc.Debugf("apply x-const: %v", r.Const)
		fieldSo.XConst = []byte(protojson.Format(r.Const))
	}
	if r.Lt != nil {
		bc.Debugf("apply x-durationLt: %v", r.Lt)
		fieldSo.XDurationLt = []byte(protojson.Format(r.Lt))
	}
	if r.Lte != nil {
		bc.Debugf("apply x-durationLte: %v", r.Lte)
		fieldSo.XDurationLte = []byte(protojson.Format(r.Lte))
	}
	if r.Gt != nil {
		bc.Debugf("apply x-durationGt: %v", r.Gt)
		fieldSo.XDurationGt = []byte(protojson.Format(r.Gt))
	}
	if r.Gte != nil {
		bc.Debugf("apply x-durationGte: %v", r.Gte)
		fieldSo.XDurationGte = []byte(protojson.Format(r.Gte))
	}
	if len(r.In) > 0 {
		bc.Debugf("apply enum from in: %v", r.In)
		for _, d := range r.In {
			fieldSo.Enum = append(fieldSo.Enum, []byte(protojson.Format(d)))
		}
	}
	if len(r.NotIn) > 0 {
		bc.Debugf("apply not.enum from not_in: %v", r.In)
		not := &schemaObject{}
		for _, d := range r.NotIn {
			not.Enum = append(not.Enum, []byte(protojson.Format(d)))
		}
		fieldSo.Not = not
	}
	return fieldSo, r.GetRequired()
}

func schemaOfWktTimestamp(bc pgs.BuildContext, rules *validate.FieldRules) (*schemaObject, bool) {
	bc.Debugf("handling Timestamp field")
	fieldSo := &schemaObject{
		Type:   schemaType{"string"},
		Format: "date-time",
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo, false
	}
	r := rules.GetTimestamp()
	if r == nil {
		bc.Debugf("timestamp rule not found")
		return fieldSo, false
	}
	if r.Const != nil {
		bc.Debugf("apply x-const: %v", r.Const)
		fieldSo.XConst = []byte(protojson.Format(r.Const))
	}
	if r.Lt != nil {
		bc.Debugf("apply x-timestampLt: %v", r.Lt)
		fieldSo.XTimestampLt = []byte(protojson.Format(r.Lt))
	}
	if r.Lte != nil {
		bc.Debugf("apply x-timestampLte: %v", r.Lte)
		fieldSo.XTimestampLte = []byte(protojson.Format(r.Lte))
	}
	if r.Gt != nil {
		bc.Debugf("apply x-timestampGt: %v", r.Gt)
		fieldSo.XTimestampGt = []byte(protojson.Format(r.Gt))
	}
	if r.Gte != nil {
		bc.Debugf("apply x-timestampGte: %v", r.Gte)
		fieldSo.XTimestampGte = []byte(protojson.Format(r.Gte))
	}
	if r.Within != nil {
		bc.Debugf("apply x-timestampWithin: %v", r.Within)
		fieldSo.XTimestampWithin = []byte(protojson.Format(r.Within))
		if r.GetLtNow() {
			bc.Debugf("apply x-timestampLtNow: %t", true)
			fieldSo.XTimestampLtNow = true
		}
		if r.GetGtNow() {
			bc.Debugf("apply x-timestampGtNow: %t", true)
			fieldSo.XTimestampGtNow = true
		}
	}
	return fieldSo, r.GetRequired()
}
