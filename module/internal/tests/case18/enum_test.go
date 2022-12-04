package case18

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

func TestEnum(t *testing.T) {
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
				"/module/internal/tests/case18/msg.pb.EnumTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type        string             `json:"type"`
			Default     *string            `json:"default"`
			Properties  map[string]*schema `json:"properties"`
			Description string             `json:"description"`
			XConst      *string            `json:"x-const"`
			Enum        []string           `json:"enum"`
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
		assert.Equal(t, "#/definitions/case18.EnumTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case18.EnumTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, `0`, *prop.Default)

			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, `0`, *prop.Default)
			assert.Nil(t, prop.XConst)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, `0`, *prop.Default)
			assert.Equal(t, ms.Properties["noRule"], prop)
		}
		{
			prop, ok := ms.Properties["const"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, `0`, *prop.Default)
			assert.NotNil(t, prop.XConst)
			assert.Equal(t, "3", *prop.XConst)
		}
		{
			prop, ok := ms.Properties["in"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, `0`, *prop.Default)

			assert.Equal(t, []string{"E1_ONE", "1", "E1_TWO", "2"}, prop.Enum)
		}
		{
			prop, ok := ms.Properties["notIn"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, `0`, *prop.Default)

			assert.Equal(t, []string{"E1_UNSPECIFIED", "0", "E1_ONE", "1", "E1_TWO", "2", "E1_THREE", "3"}, prop.Enum)
		}
		{
			prop, ok := ms.Properties["inNotIn"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, `0`, *prop.Default)

			assert.Equal(t, []string{"E1_ONE", "1", "E1_TWO", "2"}, prop.Enum)
		}
	})
}
