package case21

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yimi-go/protoc-gen-validate-jsonschema/module/internal/tests/base"
)

func TestAny(t *testing.T) {
	base.FrameTest(t, "msg.pb.bin", func(t *testing.T, fs afero.Fs, res *bytes.Buffer) {
		resp := &pluginpb.CodeGeneratorResponse{}
		err := proto.Unmarshal(res.Bytes(), resp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, resp.GetFile(), 1)

		file := resp.GetFile()[0]
		assert.NotNil(t, file)
		assert.Equal(t,
			"github.com/yimi-go/protoc-gen-validate-jsonschema"+
				"/module/internal/tests/case21/msg.pb.AnyTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type                 any                `json:"type"`
			Default              any                `json:"default"`
			Properties           map[string]*schema `json:"properties"`
			Required             []string           `json:"required"`
			AdditionalProperties *schema            `json:"additionalProperties"`
			Description          string             `json:"description"`
			Enum                 []any              `json:"enum"`
			Not                  *schema            `json:"not"`
		}
		type schemaFile struct {
			Schema      string            `json:"$schema"`
			Ref         string            `json:"$ref"`
			Definitions map[string]schema `json:"definitions"`
		}
		var sf schemaFile
		err = json.Unmarshal([]byte(file.GetContent()), &sf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "http://json-schema.org/draft-04/schema#", sf.Schema)
		assert.Equal(t, "#/definitions/case21.AnyTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case21.AnyTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "aaa\n\nbbb", prop.Description)

			assert.Len(t, prop.Properties, 1)
			assert.Contains(t, prop.Properties, "@type")
			assert.NotNil(t, prop.Properties["@type"])
			assert.Equal(t, "string", prop.Properties["@type"].Type)
			assert.Equal(t, []string{"@type"}, prop.Required)
			assert.NotNil(t, prop.AdditionalProperties)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Len(t, prop.Properties, 1)
			assert.Contains(t, prop.Properties, "@type")
			assert.NotNil(t, prop.Properties["@type"])
			assert.Equal(t, "string", prop.Properties["@type"].Type)
			assert.Equal(t, []string{"@type"}, prop.Required)
			assert.NotNil(t, prop.AdditionalProperties)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Len(t, prop.Properties, 1)
			assert.Contains(t, prop.Properties, "@type")
			assert.NotNil(t, prop.Properties["@type"])
			assert.Equal(t, "string", prop.Properties["@type"].Type)
			assert.Equal(t, []string{"@type"}, prop.Required)
			assert.NotNil(t, prop.AdditionalProperties)
		}
		{
			prop, ok := ms.Properties["required"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Len(t, prop.Properties, 1)
			assert.Contains(t, prop.Properties, "@type")
			assert.NotNil(t, prop.Properties["@type"])
			assert.Equal(t, "string", prop.Properties["@type"].Type)
			assert.Equal(t, []string{"@type"}, prop.Required)
			assert.NotNil(t, prop.AdditionalProperties)

			assert.Contains(t, ms.Required, "required")
		}
		{
			prop, ok := ms.Properties["in"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Len(t, prop.Properties, 1)
			assert.Contains(t, prop.Properties, "@type")
			assert.NotNil(t, prop.Properties["@type"])
			assert.Equal(t, "string", prop.Properties["@type"].Type)
			assert.Equal(t, []string{"@type"}, prop.Required)
			assert.NotNil(t, prop.AdditionalProperties)

			assert.Equal(t, []any{"a/b", "a/c"}, prop.Enum)
		}
		{
			prop, ok := ms.Properties["notIn"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Len(t, prop.Properties, 1)
			assert.Contains(t, prop.Properties, "@type")
			assert.NotNil(t, prop.Properties["@type"])
			assert.Equal(t, "string", prop.Properties["@type"].Type)
			assert.Equal(t, []string{"@type"}, prop.Required)
			assert.NotNil(t, prop.AdditionalProperties)

			assert.NotNil(t, prop.Not)
			assert.Equal(t, []any{"x/y", "x/z"}, prop.Not.Enum)
		}
	})
}
