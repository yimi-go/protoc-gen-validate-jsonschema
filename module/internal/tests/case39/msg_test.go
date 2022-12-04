package case39

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

func TestMsg(t *testing.T) {
	base.FrameTest(t, "msg.pb.bin", func(t *testing.T, fs afero.Fs, res *bytes.Buffer) {
		resp := &pluginpb.CodeGeneratorResponse{}
		err := proto.Unmarshal(res.Bytes(), resp)
		if err != nil {
			t.Fatal(err)
		}
		assert.Len(t, resp.GetFile(), 2)

		var file *pluginpb.CodeGeneratorResponse_File
		for _, f := range resp.GetFile() {
			if strings.Contains(f.GetName(), "MsgTestMsg") {
				file = f
				break
			}
		}
		assert.NotNil(t, file)
		assert.Equal(t,
			"github.com/yimi-go/protoc-gen-validate-jsonschema"+
				"/module/internal/tests/case39/msg.pb.MsgTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type        any                `json:"type"`
			Properties  map[string]*schema `json:"properties"`
			Required    []string           `json:"required"`
			Description string             `json:"description"`
			Ref         string             `json:"$ref"`
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
		assert.Equal(t, "#/definitions/case39.MsgTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case39.MsgTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			def := strings.TrimPrefix(prop.Ref, "#/definitions/")
			_, ok = sf.Definitions[def]
			assert.True(t, ok)

			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			def := strings.TrimPrefix(prop.Ref, "#/definitions/")
			_, ok = sf.Definitions[def]
			assert.True(t, ok)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			def := strings.TrimPrefix(prop.Ref, "#/definitions/")
			_, ok = sf.Definitions[def]
			assert.True(t, ok)
		}
		{
			prop, ok := ms.Properties["required"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Nil(t, prop.Type)
			def := strings.TrimPrefix(prop.Ref, "#/definitions/")
			_, ok = sf.Definitions[def]
			assert.True(t, ok)

			assert.Contains(t, ms.Required, "required")
		}
	})
}
