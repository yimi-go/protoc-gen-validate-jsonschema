package case35

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

func TestValue(t *testing.T) {
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
				"/module/internal/tests/case35/msg.pb.ValueTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type        any                `json:"type"`
			Default     any                `json:"default"`
			Properties  map[string]*schema `json:"properties"`
			Required    []string           `json:"required"`
			Description string             `json:"description"`
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
		assert.Equal(t, "#/definitions/case35.ValueTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case35.ValueTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			assert.Nil(t, prop.Default)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			assert.Nil(t, prop.Default)
		}
		{
			prop, ok := ms.Properties["required"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			assert.Nil(t, prop.Default)

			assert.Contains(t, ms.Required, "required")
		}
	})
}