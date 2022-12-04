package case15

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

func TestBool(t *testing.T) {
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
				"/module/internal/tests/case15/msg.pb.BoolTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type        string             `json:"type"`
			Default     *bool              `json:"default"`
			Properties  map[string]*schema `json:"properties"`
			Description string             `json:"description"`
			XConst      *bool              `json:"x-const"`
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
		assert.Equal(t, "#/definitions/case15.BoolTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case15.BoolTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "boolean", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, false, *prop.Default)
			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}

		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "boolean", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, false, *prop.Default)
			assert.Nil(t, prop.XConst)
		}

		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "boolean", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, false, *prop.Default)
			assert.Equal(t, ms.Properties["noRule"], prop)
		}

		{
			prop, ok := ms.Properties["const"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "boolean", prop.Type)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, false, *prop.Default)
			assert.NotNil(t, prop.XConst)
			assert.Equal(t, true, *prop.XConst)
		}
	})
}
