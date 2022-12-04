package case30

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yimi-go/protoc-gen-validate-jsonschema/module/internal/tests/base"
)

func TestBytesWrapper(t *testing.T) {
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
				"/module/internal/tests/case30/msg.pb.BytesWrapperTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type            string             `json:"type"`
			Default         *string            `json:"default"`
			Properties      map[string]*schema `json:"properties"`
			Description     string             `json:"description"`
			XConst          *string            `json:"x-const"`
			XBytesMinLength uint64             `json:"x-bytesMinLength"`
			XBytesMaxLength uint64             `json:"x-bytesMaxLength"`
			AllOf           []*schema          `json:"allOf"`
			Not             *schema            `json:"not"`
			Enum            []string           `json:"enum"`
			Format          string             `json:"format"`
			AnyOf           []*schema          `json:"anyOf"`
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
		assert.Equal(t, "#/definitions/case30.BytesWrapperTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case30.BytesWrapperTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)
			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)
			assert.Nil(t, prop.XConst)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)
			assert.Equal(t, ms.Properties["noRule"], prop)
		}
		{
			prop, ok := ms.Properties["const"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)
			assert.NotNil(t, prop.XConst)
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("abc")), *prop.XConst)
		}
		{
			prop, ok := ms.Properties["len"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(3), prop.XBytesMinLength)
			assert.Equal(t, uint64(0), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minLenLtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minLenEqLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minLenGtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(0), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(0), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxLenLtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(0), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxLenEqLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxLenGtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["in"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, []string{
				base64.StdEncoding.EncodeToString([]byte("a")),
				base64.StdEncoding.EncodeToString([]byte("b")),
				base64.StdEncoding.EncodeToString([]byte("c")),
			}, prop.Enum)
		}
		{
			prop, ok := ms.Properties["notIn"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Not)
			assert.Equal(t, []string{
				base64.StdEncoding.EncodeToString([]byte("x")),
				base64.StdEncoding.EncodeToString([]byte("y")),
				base64.StdEncoding.EncodeToString([]byte("z")),
			}, prop.Not.Enum)
		}
	})
}
