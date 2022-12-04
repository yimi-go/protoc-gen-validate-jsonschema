package generate

import (
	"encoding/json"
)

type SchemaFile struct {
	Schema      string                   `json:"$schema"`
	Ref         string                   `json:"$ref"`
	Definitions map[string]*schemaObject `json:"definitions"`
}

type schemaObject struct {
	Description string `json:"description,omitempty"`

	Ref string `json:"$ref,omitempty"`

	Type    schemaType        `json:"type,omitempty"`
	Format  string            `json:"format,omitempty"`
	Default json.RawMessage   `json:"default,omitempty"`
	Enum    []json.RawMessage `json:"enum,omitempty"`

	Properties           map[string]*schemaObject `json:"properties,omitempty"`
	AdditionalProperties *schemaObject            `json:"additionalProperties,omitempty"`

	MultipleOf       float64  `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`

	MaxLength *uint64 `json:"maxLength,omitempty"`
	MinLength *uint64 `json:"minLength,omitempty"`
	Pattern   string  `json:"pattern,omitempty"`

	MaxItems    uint64        `json:"maxItems,omitempty"`
	MinItems    uint64        `json:"minItems,omitempty"`
	UniqueItems bool          `json:"uniqueItems,omitempty"`
	Items       *schemaObject `json:"items,omitempty"`

	MaxProperties uint64   `json:"maxProperties,omitempty"`
	MinProperties uint64   `json:"minProperties,omitempty"`
	Required      []string `json:"required,omitempty"`

	Not   *schemaObject   `json:"not,omitempty"`
	AllOf []*schemaObject `json:"allOf,omitempty"`
	AnyOf []*schemaObject `json:"anyOf,omitempty"`

	XOneOfs               []*OneOf        `json:"x-oneOfs,omitempty"`
	XConst                json.RawMessage `json:"x-const,omitempty"`
	XBytesMaxLength       *uint64         `json:"x-bytesMaxLength,omitempty"`
	XBytesMinLength       *uint64         `json:"x-bytesMinLength,omitempty"`
	XWellKnownRegex       string          `json:"x-wellKnownRegex,omitempty"`
	XWellKnownRegexStrict bool            `json:"x-wellKnownRegexStrict,omitempty"`
	XDurationLt           json.RawMessage `json:"x-durationLt,omitempty"`
	XDurationLte          json.RawMessage `json:"x-durationLte,omitempty"`
	XDurationGt           json.RawMessage `json:"x-durationGt,omitempty"`
	XDurationGte          json.RawMessage `json:"x-durationGte,omitempty"`
	XTimestampLt          json.RawMessage `json:"x-timestampLt,omitempty"`
	XTimestampLte         json.RawMessage `json:"x-timestampLte,omitempty"`
	XTimestampGt          json.RawMessage `json:"x-timestampGt,omitempty"`
	XTimestampGte         json.RawMessage `json:"x-timestampGte,omitempty"`
	XTimestampLtNow       bool            `json:"x-timestampLtNow,omitempty"`
	XTimestampGtNow       bool            `json:"x-timestampGtNow,omitempty"`
	XTimestampWithin      json.RawMessage `json:"x-timestampWithin,omitempty"`
}

type schemaType []string

func (st schemaType) MarshalJSON() ([]byte, error) {
	if len(st) == 0 {
		return []byte("null"), nil
	}
	if len(st) == 1 {
		return json.Marshal(st[0])
	}
	return json.Marshal([]string(st))
}

type OneOf struct {
	Name     string   `json:"name"`
	Required bool     `json:"required,omitempty"`
	Fields   []string `json:"fields"`
}
