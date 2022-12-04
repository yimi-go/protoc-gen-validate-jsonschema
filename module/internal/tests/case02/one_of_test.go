package case02

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

func TestOneOf(t *testing.T) {
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
				"/module/internal/tests/case02/msg.pb.OneOfMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type oneOf struct {
			Name     string   `json:"name"`
			Required bool     `json:"required"`
			Fields   []string `json:"fields"`
		}
		type schema struct {
			Type       string            `json:"type"`
			Properties map[string]schema `json:"properties"`
			XOneOfs    []oneOf           `json:"x-oneOfs"`
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
		assert.Equal(t, "#/definitions/case02.OneOfMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case02.OneOfMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)
		assert.NotEmpty(t, ms.Properties)
		assert.NotEmpty(t, ms.XOneOfs)
		assert.Len(t, ms.XOneOfs, 2)
		for _, oo := range ms.XOneOfs {
			assert.NotNil(t, oo)
			if oo.Name == "required" {
				assert.True(t, oo.Required)
				assert.Equal(t, []string{"a", "b"}, oo.Fields)
			} else {
				assert.Equal(t, "n", oo.Name)
				assert.False(t, oo.Required)
				assert.Equal(t, []string{"c", "d"}, oo.Fields)
			}
		}
	})
}
