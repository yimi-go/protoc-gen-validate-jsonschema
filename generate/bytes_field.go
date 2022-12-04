package generate

import (
	"encoding/base64"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfBytesFieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling bytes field")
	fieldSo := &schemaObject{
		Type:    schemaType{"string"},
		Format:  "bytes",
		Default: []byte(`""`),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetBytes()
	if r == nil {
		bc.Debugf("bytes rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %x", r.Const)
		encoded := base64.StdEncoding.EncodeToString(r.Const)
		fieldSo.XConst = []byte(`"` + encoded + `"`)
	}
	if r.MinLen != nil || r.Len != nil {
		if r.MinLen != nil && r.Len != nil {
			bc.Debugf("both minLen and len are set")
			if *r.MinLen > *r.Len {
				// false rule
				bc.Logf("false rule: minLen: %d, Len: %d", *r.MinLen, *r.Len)
			} else {
				bc.Debugf("apply x-bytesMinLength from len: %d", *r.Len)
				fieldSo.XBytesMinLength = r.Len
			}
		} else if r.MinLen == nil {
			// r.Len != nil
			bc.Debugf("apply x-bytesMinLength from len: %d", *r.Len)
			fieldSo.XBytesMinLength = r.Len
		} else {
			// r.Len == nil && r.MinLen != nil
			bc.Debugf("apply x-bytesMinLength from minLen: %d", *r.MinLen)
			fieldSo.XBytesMinLength = r.MinLen
		}
	}
	if r.MaxLen != nil || r.Len != nil {
		if r.MaxLen != nil && r.Len != nil {
			bc.Debugf("both maxLen and len are set")
			if *r.MaxLen < *r.Len {
				// false rule
				bc.Logf("false rule: maxLen: %d, Len: %d", *r.MaxLen, *r.Len)
			} else {
				bc.Debugf("apply x-bytesMaxLength from len: %d", *r.Len)
				fieldSo.XBytesMaxLength = r.Len
			}
		} else if r.MaxLen == nil {
			// r.Len != nil
			bc.Debugf("apply x-bytesMaxLength from len: %d", *r.Len)
			fieldSo.XBytesMaxLength = r.Len
		} else {
			// r.Len == nil && r.MaxLen != nil
			bc.Debugf("apply x-bytesMaxLength from maxLen: %d", *r.MaxLen)
			fieldSo.XBytesMaxLength = r.MaxLen
		}
	}
	if r.Pattern != nil && "" != *r.Pattern {
		bc.Debugf("how to describe bytes pattern in json schema?")
	}
	if len(r.Prefix) != 0 {
		bc.Debugf("how to describe bytes prefix in json schema?")
	}
	if len(r.Suffix) != 0 {
		bc.Debugf("how to describe bytes suffix in json schema?")
	}
	if len(r.Contains) != 0 {
		bc.Debugf("how to describe bytes contains in json schema?")
	}
	if len(r.In) != 0 {
		bc.Debugf("apply enum from in: %v", r.In)
		for _, s := range r.In {
			encoded := base64.StdEncoding.EncodeToString(s)
			fieldSo.Enum = append(fieldSo.Enum, []byte(`"`+encoded+`"`))
		}
	}
	if len(r.NotIn) != 0 {
		bc.Debugf("apply not enum from notIn: %v", r.NotIn)
		if fieldSo.Not == nil {
			fieldSo.Not = &schemaObject{}
		}
		for _, s := range r.NotIn {
			encoded := base64.StdEncoding.EncodeToString(s)
			fieldSo.Not.Enum = append(fieldSo.Not.Enum, []byte(`"`+encoded+`"`))
		}
	}
	if r.GetWellKnown() != nil {
		bc.Debugf("how to describe bytes well-known patterns in json schema?")
	}
	return fieldSo
}
