package case20

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yimi-go/protoc-gen-validate-jsonschema/module/internal/tests/base"
)

func TestMap(t *testing.T) {
	base.FrameTest(t, "msg.pb.bin", func(t *testing.T, fs afero.Fs, res *bytes.Buffer) {
		resp := &pluginpb.CodeGeneratorResponse{}
		err := proto.Unmarshal(res.Bytes(), resp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, resp.GetFile(), 2)

		var file *pluginpb.CodeGeneratorResponse_File
		for _, f := range resp.GetFile() {
			if strings.Contains(f.GetName(), "Map") {
				file = f
				break
			}
		}
		assert.NotNil(t, file)
		assert.Equal(t,
			"github.com/yimi-go/protoc-gen-validate-jsonschema"+
				"/module/internal/tests/case20/msg.pb.MapTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type                 any                `json:"type"`
			Format               string             `json:"format"`
			Default              any                `json:"default"`
			Properties           map[string]*schema `json:"properties"`
			AdditionalProperties *schema            `json:"additionalProperties"`
			MaxProperties        uint64             `json:"maxProperties"`
			MinProperties        uint64             `json:"minProperties"`
			Description          string             `json:"description"`
			XConst               any                `json:"x-const"`
			Items                *schema            `json:"items"`
			MinItems             uint64             `json:"minItems"`
			MaxItems             uint64             `json:"maxItems"`
			Minimum              *float64           `json:"minimum"`
			ExclusiveMinimum     bool               `json:"exclusiveMinimum"`
			Maximum              *float64           `json:"maximum"`
			UniqueItems          bool               `json:"uniqueItems"`
			Enum                 []any              `json:"enum"`
			Ref                  string             `json:"$ref"`
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
		assert.Equal(t, "#/definitions/case20.MapTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case20.MapTestMsg"]
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

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.Equal(t, []any{"integer", "string"}, prop.AdditionalProperties.Type)
			assert.Equal(t, "int64", prop.AdditionalProperties.Format)
			assert.Equal(t, 0.0, prop.AdditionalProperties.Default)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.Equal(t, "integer", prop.AdditionalProperties.Type)
			assert.Equal(t, "int64", prop.AdditionalProperties.Format)
			assert.Equal(t, 0.0, prop.AdditionalProperties.Default)
			assert.NotNil(t, prop.AdditionalProperties.Minimum)
			assert.Equal(t, 0.0, *prop.AdditionalProperties.Minimum)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.Equal(t, []any{"integer", "string"}, prop.AdditionalProperties.Type)
			assert.Equal(t, "uint64", prop.AdditionalProperties.Format)
			assert.Equal(t, 0.0, prop.AdditionalProperties.Default)
			assert.NotNil(t, prop.AdditionalProperties.Minimum)
			assert.Equal(t, 0.0, *prop.AdditionalProperties.Minimum)
		}
		{
			prop, ok := ms.Properties["maxPairs"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.Equal(t, "integer", prop.AdditionalProperties.Type)
			assert.Equal(t, "int32", prop.AdditionalProperties.Format)
			assert.Equal(t, 0.0, prop.AdditionalProperties.Default)

			assert.Equal(t, uint64(5), prop.MaxProperties)
		}
		{
			prop, ok := ms.Properties["minPairs"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.Equal(t, []any{"integer", "string"}, prop.AdditionalProperties.Type)
			assert.Equal(t, "int64", prop.AdditionalProperties.Format)
			assert.Equal(t, 0.0, prop.AdditionalProperties.Default)

			assert.Equal(t, uint64(5), prop.MinProperties)
		}
		{
			prop, ok := ms.Properties["values"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.Equal(t, "boolean", prop.AdditionalProperties.Type)
			assert.Equal(t, "", prop.AdditionalProperties.Format)
			assert.Equal(t, false, prop.AdditionalProperties.Default)

			assert.NotNil(t, prop.AdditionalProperties.XConst)
			assert.Equal(t, true, prop.AdditionalProperties.XConst)
		}
		{
			prop, ok := ms.Properties["evs"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.Equal(t, "string", prop.AdditionalProperties.Type)
			assert.Equal(t, "", prop.AdditionalProperties.Format)
			assert.Equal(t, "0", prop.AdditionalProperties.Default)
			assert.Equal(t, []any{"E3_ONE", "1", "E3_TWO", "2", "E3_THREE", "3"}, prop.AdditionalProperties.Enum)
		}
		{
			prop, ok := ms.Properties["mvs"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "object", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Nil(t, prop.Properties)
			assert.NotNil(t, prop.AdditionalProperties)
			assert.NotEmpty(t, prop.AdditionalProperties.Ref)
			def := strings.TrimPrefix(prop.AdditionalProperties.Ref, "#/definitions/")
			_, ok = sf.Definitions[def]
			assert.True(t, ok)
		}
	})
}
