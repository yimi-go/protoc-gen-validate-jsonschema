package case19

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

func TestRepeated(t *testing.T) {
	base.FrameTest(t, "msg.pb.bin", func(t *testing.T, fs afero.Fs, res *bytes.Buffer) {
		resp := &pluginpb.CodeGeneratorResponse{}
		err := proto.Unmarshal(res.Bytes(), resp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, resp.GetFile(), 2)

		var file *pluginpb.CodeGeneratorResponse_File
		for _, f := range resp.GetFile() {
			if strings.Contains(f.GetName(), "Repeated") {
				file = f
				break
			}
		}
		assert.NotNil(t, file)
		assert.Equal(t,
			"github.com/yimi-go/protoc-gen-validate-jsonschema"+
				"/module/internal/tests/case19/msg.pb.RepeatedTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type             any                `json:"type"`
			Format           string             `json:"format"`
			Default          any                `json:"default"`
			Properties       map[string]*schema `json:"properties"`
			Description      string             `json:"description"`
			XConst           any                `json:"x-const"`
			Items            *schema            `json:"items"`
			MinItems         uint64             `json:"minItems"`
			MaxItems         uint64             `json:"maxItems"`
			Minimum          *float64           `json:"minimum"`
			ExclusiveMinimum bool               `json:"exclusiveMinimum"`
			Maximum          *float64           `json:"maximum"`
			UniqueItems      bool               `json:"uniqueItems"`
			Enum             []any              `json:"enum"`
			Ref              string             `json:"$ref"`
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
		assert.Equal(t, "#/definitions/case19.RepeatedTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case19.RepeatedTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "aaa\n\nbbb", prop.Description)

			assert.NotNil(t, prop.Items)
			assert.Equal(t, "number", prop.Items.Type)
			assert.Equal(t, "float", prop.Items.Format)
			assert.Equal(t, 0.0, prop.Items.Default)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Items)
			assert.Equal(t, "number", prop.Items.Type)
			assert.Equal(t, "double", prop.Items.Format)
			assert.Equal(t, 0.0, prop.Items.Default)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Items)
			assert.Equal(t, "integer", prop.Items.Type)
			assert.Equal(t, "int32", prop.Items.Format)
			assert.Equal(t, 0.0, prop.Items.Default)
		}
		{
			prop, ok := ms.Properties["maxItems"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Items)
			assert.Equal(t, []any{"integer", "string"}, prop.Items.Type)
			assert.Equal(t, "int64", prop.Items.Format)
			assert.Equal(t, 0.0, prop.Items.Default)

			assert.Equal(t, uint64(3), prop.MaxItems)
		}
		{
			prop, ok := ms.Properties["minItems"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Items)
			assert.Equal(t, "integer", prop.Items.Type)
			assert.Equal(t, "int64", prop.Items.Format)
			assert.Equal(t, 0.0, prop.Items.Default)
			assert.NotNil(t, prop.Items.Minimum)
			assert.Equal(t, 0.0, *prop.Items.Minimum)

			assert.Equal(t, uint64(3), prop.MinItems)
		}
		{
			prop, ok := ms.Properties["unique"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Items)
			assert.Equal(t, []any{"integer", "string"}, prop.Items.Type)
			assert.Equal(t, "uint64", prop.Items.Format)
			assert.Equal(t, 0.0, prop.Items.Default)
			assert.NotNil(t, prop.Items.Minimum)
			assert.Equal(t, 0.0, *prop.Items.Minimum)

			assert.True(t, prop.UniqueItems)
		}
		{
			prop, ok := ms.Properties["si32"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Items)
			assert.Equal(t, "integer", prop.Items.Type)
			assert.Equal(t, "int32", prop.Items.Format)
			assert.Equal(t, 0.0, prop.Items.Default)
			assert.Equal(t, []any{float64(1), float64(2)}, prop.Items.Enum)
		}
		{
			prop, ok := ms.Properties["es"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.True(t, prop.UniqueItems)
			assert.NotNil(t, prop.Items)
			assert.Equal(t, "string", prop.Items.Type)
			assert.Equal(t, "", prop.Items.Format)
			assert.Equal(t, "0", prop.Items.Default)
			assert.Equal(t, []any{"E2_TWO", "2", "E2_THREE", "3"}, prop.Items.Enum)
		}
		{
			prop, ok := ms.Properties["ms"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "array", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Items)
			assert.NotEmpty(t, prop.Items.Ref)
			def := strings.TrimPrefix(prop.Items.Ref, "#/definitions/")
			_, ok = sf.Definitions[def]
			assert.True(t, ok)
		}
	})
}
