package case29

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

func TestStringWrapper(t *testing.T) {
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
				"/module/internal/tests/case29/msg.pb.StringWrapperTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type                  string             `json:"type"`
			Default               *string            `json:"default"`
			Properties            map[string]*schema `json:"properties"`
			Description           string             `json:"description"`
			XConst                *string            `json:"x-const"`
			MinLength             uint64             `json:"minLength"`
			MaxLength             uint64             `json:"maxLength"`
			XBytesMinLength       uint64             `json:"x-bytesMinLength"`
			XBytesMaxLength       uint64             `json:"x-bytesMaxLength"`
			Pattern               string             `json:"pattern"`
			AllOf                 []*schema          `json:"allOf"`
			Not                   *schema            `json:"not"`
			Enum                  []string           `json:"enum"`
			Format                string             `json:"format"`
			AnyOf                 []*schema          `json:"anyOf"`
			XWellKnownRegex       string             `json:"x-wellKnownRegex"`
			XWellKnownRegexStrict bool               `json:"x-wellKnownRegexStrict"`
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
		assert.Equal(t, "#/definitions/case29.StringWrapperTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case29.StringWrapperTestMsg"]
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
			assert.Equal(t, "abc", *prop.XConst)
		}
		{
			prop, ok := ms.Properties["len"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.MinLength)
			assert.Equal(t, uint64(5), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["minLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(3), prop.MinLength)
			assert.Equal(t, uint64(0), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["minLenLtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.MinLength)
			assert.Equal(t, uint64(5), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["minLenEqLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.MinLength)
			assert.Equal(t, uint64(5), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["minLenGtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(0), prop.MinLength)
			assert.Equal(t, uint64(5), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["maxLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(0), prop.MinLength)
			assert.Equal(t, uint64(5), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["maxLenLtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.MinLength)
			assert.Equal(t, uint64(0), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["maxLenEqLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.MinLength)
			assert.Equal(t, uint64(5), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["maxLenGtLen"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.MinLength)
			assert.Equal(t, uint64(5), prop.MaxLength)
		}
		{
			prop, ok := ms.Properties["lenBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(0), prop.MinLength)
			assert.Equal(t, uint64(0), prop.MaxLength)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(0), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minBytesLtLenBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minBytesEqLenBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["minBytesGtLenBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(0), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(0), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxBytesLtLenBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(0), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxBytesEqLenBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["maxBytesGtLenBytes"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, uint64(5), prop.XBytesMinLength)
			assert.Equal(t, uint64(5), prop.XBytesMaxLength)
		}
		{
			prop, ok := ms.Properties["pattern"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "a.*b", prop.Pattern)
		}
		{
			prop, ok := ms.Properties["prefix"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, `^a\*.*$`, prop.Pattern)
		}
		{
			prop, ok := ms.Properties["suffix"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, `^.*\.z$`, prop.Pattern)
		}
		{
			prop, ok := ms.Properties["contains"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, `^.*\(op\)q.*$`, prop.Pattern)
		}
		{
			prop, ok := ms.Properties["prefixSuffix"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Empty(t, prop.Pattern)
			assert.Len(t, prop.AllOf, 2)
			assert.Equal(t, "^a.*$", prop.AllOf[0].Pattern)
			assert.Equal(t, "^.*z$", prop.AllOf[1].Pattern)
		}
		{
			prop, ok := ms.Properties["notContains"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Not)
			assert.Equal(t, `^.*z\*z.*$`, prop.Not.Pattern)
		}
		{
			prop, ok := ms.Properties["in"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, []string{"a", "b", "c"}, prop.Enum)
		}
		{
			prop, ok := ms.Properties["notIn"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Not)
			assert.Equal(t, []string{"x", "y", "z"}, prop.Not.Enum)
		}
		{
			prop, ok := ms.Properties["notInNotContains"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Not)
			assert.Equal(t, []string{"a", "b"}, prop.Not.Enum)
			assert.Equal(t, `^.*xyz.*$`, prop.Not.Pattern)
		}
		{
			prop, ok := ms.Properties["email"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "email", prop.Format)
		}
		{
			prop, ok := ms.Properties["hostname"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "hostname", prop.Format)
		}
		{
			prop, ok := ms.Properties["ip"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Len(t, prop.AnyOf, 2)
			assert.Equal(t, "ipv4", prop.AnyOf[0].Format)
			assert.Equal(t, "ipv6", prop.AnyOf[1].Format)
		}
		{
			prop, ok := ms.Properties["ipv4"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "ipv4", prop.Format)
		}
		{
			prop, ok := ms.Properties["ipv6"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "ipv6", prop.Format)
		}
		{
			prop, ok := ms.Properties["uri"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "uri", prop.Format)
		}
		{
			prop, ok := ms.Properties["uriRef"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "uri-reference", prop.Format)
		}
		{
			prop, ok := ms.Properties["address"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Len(t, prop.AnyOf, 3)
			assert.Equal(t, "ipv4", prop.AnyOf[0].Format)
			assert.Equal(t, "ipv6", prop.AnyOf[1].Format)
			assert.Equal(t, "hostname", prop.AnyOf[2].Format)
		}
		{
			prop, ok := ms.Properties["uuid"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "uuid", prop.Format)
		}
		{
			prop, ok := ms.Properties["wellKnownRegexHttpHeaderName"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "", prop.Format)
			assert.Equal(t, "HTTP_HEADER_NAME", prop.XWellKnownRegex)
		}
		{
			prop, ok := ms.Properties["wellKnownRegexHttpHeaderValue"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "", prop.Format)
			assert.Equal(t, "HTTP_HEADER_VALUE", prop.XWellKnownRegex)
		}
		{
			prop, ok := ms.Properties["wellKnownRegexUnknown"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "", prop.Format)
			assert.Equal(t, "", prop.XWellKnownRegex)
		}
		{
			prop, ok := ms.Properties["strict"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "", prop.Format)
			assert.Equal(t, "", prop.XWellKnownRegex)
			assert.Equal(t, true, prop.XWellKnownRegexStrict)
		}
	})
}
