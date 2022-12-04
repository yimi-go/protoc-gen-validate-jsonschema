package generate

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func schemaOfStringFieldType(bc pgs.BuildContext, rules *validate.FieldRules) *schemaObject {
	bc.Debugf("handling string field")
	fieldSo := &schemaObject{
		Type:    schemaType{"string"},
		Default: []byte(`""`),
	}
	if rules == nil {
		bc.Debugf("no rules to apply")
		return fieldSo
	}
	r := rules.GetString_()
	if r == nil {
		bc.Debugf("string rule not found")
		return fieldSo
	}
	if r.Const != nil {
		bc.Debugf("const: %v", *r.Const)
		fieldSo.XConst = []byte(strconv.QuoteToASCII(*r.Const))
	}
	if r.MinLen != nil || r.Len != nil {
		if r.MinLen != nil && r.Len != nil {
			bc.Debugf("both minLen and len are set")
			if *r.MinLen > *r.Len {
				// false rule
				bc.Logf("false rule: len: %d, minLen: %d", *r.Len, *r.MinLen)
			} else {
				bc.Debugf("apply minLength from len: %d", *r.Len)
				fieldSo.MinLength = r.Len
			}
		} else if r.MinLen == nil {
			// r.Len != nil
			bc.Debugf("apply minLength from len: %d", *r.Len)
			fieldSo.MinLength = r.Len
		} else {
			// r.Len == nil && r.MinLen != nil
			bc.Debugf("apply minLength from minLen: %d", *r.MinLen)
			fieldSo.MinLength = r.MinLen
		}
	}
	if r.MaxLen != nil || r.Len != nil {
		if r.MaxLen != nil && r.Len != nil {
			bc.Debugf("both maxLen and len are set")
			if *r.MaxLen < *r.Len {
				// false rule
				bc.Logf("false rule: len: %d, maxLen: %d", *r.Len, *r.MaxLen)
			} else {
				bc.Debugf("apply maxLength from len: %d", *r.Len)
				fieldSo.MaxLength = r.Len
			}
		} else if r.MaxLen == nil {
			// r.Len != nil
			bc.Debugf("apply maxLength from len: %d", *r.Len)
			fieldSo.MaxLength = r.Len
		} else {
			// r.Len == nil && r.MaxLen != nil
			bc.Debugf("apply maxLength from maxLen: %d", *r.MaxLen)
			fieldSo.MaxLength = r.MaxLen
		}
	}
	if r.MinBytes != nil || r.LenBytes != nil {
		if r.MinBytes != nil && r.LenBytes != nil {
			bc.Debugf("both minBytes and lenBytes are set")
			if *r.MinBytes > *r.LenBytes {
				// false rule
				bc.Logf("false rule: minBytes: %d, LenBytes: %d", *r.MinBytes, *r.LenBytes)
			} else {
				bc.Debugf("apply x-bytesMinLength from lenBytes: %d", *r.LenBytes)
				fieldSo.XBytesMinLength = r.LenBytes
			}
		} else if r.MinBytes == nil {
			// r.LenBytes != nil
			bc.Debugf("apply x-bytesMinLength from lenBytes: %d", *r.LenBytes)
			fieldSo.XBytesMinLength = r.LenBytes
		} else {
			// r.LenBytes == nil && r.MinBytes != nil
			bc.Debugf("apply x-bytesMinLength from minBytes: %d", *r.MinBytes)
			fieldSo.XBytesMinLength = r.MinBytes
		}
	}
	if r.MaxBytes != nil || r.LenBytes != nil {
		if r.MaxBytes != nil && r.LenBytes != nil {
			bc.Debugf("both maxBytes and lenBytes are set")
			if *r.MaxBytes < *r.LenBytes {
				// false rule
				bc.Logf("false rule: maxBytes: %d, LenBytes: %d", *r.MaxBytes, *r.LenBytes)
			} else {
				bc.Debugf("apply x-bytesMaxLength from lenBytes: %d", *r.LenBytes)
				fieldSo.XBytesMaxLength = r.LenBytes
			}
		} else if r.MaxBytes == nil {
			// r.LenBytes != nil
			bc.Debugf("apply x-bytesMaxLength from lenBytes: %d", *r.LenBytes)
			fieldSo.XBytesMaxLength = r.LenBytes
		} else {
			// r.LenBytes == nil && r.MaxBytes != nil
			bc.Debugf("apply x-bytesMaxLength from maxBytes: %d", *r.MaxBytes)
			fieldSo.XBytesMaxLength = r.MaxBytes
		}
	}
	if r.Pattern != nil && "" != *r.Pattern {
		bc.Debugf("apply pattern: %s", *r.Pattern)
		fieldSo.AllOf = append(fieldSo.AllOf, &schemaObject{Pattern: *r.Pattern})
	}
	if r.Prefix != nil && "" != *r.Prefix {
		bc.Debugf("apply pattern from prefix: %s", *r.Prefix)
		prefix := regexp.QuoteMeta(*r.Prefix)
		fieldSo.AllOf = append(fieldSo.AllOf, &schemaObject{Pattern: fmt.Sprintf("^%s.*$", prefix)})
	}
	if r.Suffix != nil && "" != *r.Suffix {
		bc.Debugf("apply pattern from suffix: %s", *r.Suffix)
		suffix := regexp.QuoteMeta(*r.Suffix)
		fieldSo.AllOf = append(fieldSo.AllOf, &schemaObject{Pattern: fmt.Sprintf("^.*%s$", suffix)})
	}
	if r.Contains != nil && "" != *r.Contains {
		bc.Debugf("apply pattern from contains: %s", *r.Contains)
		suffix := regexp.QuoteMeta(*r.Contains)
		fieldSo.AllOf = append(fieldSo.AllOf, &schemaObject{Pattern: fmt.Sprintf("^.*%s.*$", suffix)})
	}
	if len(fieldSo.AllOf) == 1 {
		bc.Debugf("extract pattern from only allOf")
		fieldSo.Pattern = fieldSo.AllOf[0].Pattern
		fieldSo.AllOf = nil
	}
	if r.NotContains != nil && "" != *r.NotContains {
		bc.Debugf("apply not patten from notContains: %s", *r.NotContains)
		notContains := regexp.QuoteMeta(*r.NotContains)
		fieldSo.Not = &schemaObject{Pattern: fmt.Sprintf("^.*%s.*$", notContains)}
	}
	if len(r.In) != 0 {
		bc.Debugf("apply enum from in: %v", r.In)
		for _, s := range r.In {
			fieldSo.Enum = append(fieldSo.Enum, []byte(strconv.QuoteToASCII(s)))
		}
	}
	if len(r.NotIn) != 0 {
		bc.Debugf("apply not enum from notIn: %v", r.NotIn)
		if fieldSo.Not == nil {
			fieldSo.Not = &schemaObject{}
		}
		for _, s := range r.NotIn {
			fieldSo.Not.Enum = append(fieldSo.Not.Enum, []byte(strconv.QuoteToASCII(s)))
		}
	}
	switch wkr := r.GetWellKnown().(type) {
	case *validate.StringRules_Email:
		if wkr.Email {
			bc.Debugf("apply format to email cause wellKnown.email = true")
			fieldSo.Format = "email"
		}
	case *validate.StringRules_Hostname:
		if wkr.Hostname {
			bc.Debugf("apply format to hostname cause wellKnown.hostname = true")
			fieldSo.Format = "hostname"
		}
	case *validate.StringRules_Ip:
		if wkr.Ip {
			bc.Debugf("apply format to ipv4 or ipv6 cause wellKnown.ip = true")
			fieldSo.AnyOf = append(fieldSo.AnyOf, &schemaObject{Format: "ipv4"})
			fieldSo.AnyOf = append(fieldSo.AnyOf, &schemaObject{Format: "ipv6"})
		}
	case *validate.StringRules_Ipv4:
		if wkr.Ipv4 {
			bc.Debugf("apply format to ipv4 cause wellKnown.ipv4 = true")
			fieldSo.Format = "ipv4"
		}
	case *validate.StringRules_Ipv6:
		if wkr.Ipv6 {
			bc.Debugf("apply format to ipv6 cause wellKnown.ipv6 = true")
			fieldSo.Format = "ipv6"
		}
	case *validate.StringRules_Uri:
		if wkr.Uri {
			bc.Debugf("apply format to uri cause wellKnown.uri = true")
			fieldSo.Format = "uri"
		}
	case *validate.StringRules_UriRef:
		if wkr.UriRef {
			bc.Debugf("apply format to uri-reference cause wellKnown.uriRef = true")
			fieldSo.Format = "uri-reference"
		}
	case *validate.StringRules_Address:
		if wkr.Address {
			bc.Debugf("apply format to ip or hostname cause wellKnown.address = true")
			fieldSo.AnyOf = append(fieldSo.AnyOf, &schemaObject{Format: "ipv4"})
			fieldSo.AnyOf = append(fieldSo.AnyOf, &schemaObject{Format: "ipv6"})
			fieldSo.AnyOf = append(fieldSo.AnyOf, &schemaObject{Format: "hostname"})
		}
	case *validate.StringRules_Uuid:
		if wkr.Uuid {
			bc.Debugf("apply format to uuid cause wellKnown.uuid = true")
			fieldSo.Format = "uuid"
		}
	case *validate.StringRules_WellKnownRegex:
		switch wkr.WellKnownRegex {
		case validate.KnownRegex_HTTP_HEADER_NAME:
			bc.Debugf("apply x-wellKnownRegex to HTTP_HEADER_NAME")
			fieldSo.XWellKnownRegex = "HTTP_HEADER_NAME"
		case validate.KnownRegex_HTTP_HEADER_VALUE:
			bc.Debugf("apply x-wellKnownRegex to HTTP_HEADER_VALUE")
			fieldSo.XWellKnownRegex = "HTTP_HEADER_VALUE"
		}
	}
	if r.Strict != nil {
		bc.Debugf("apply x-wellKnownRegexStrict")
		fieldSo.XWellKnownRegexStrict = *r.Strict
	}
	return fieldSo
}
