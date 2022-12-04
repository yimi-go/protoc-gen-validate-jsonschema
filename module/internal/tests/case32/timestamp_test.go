package case32

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

func TestTimestamp(t *testing.T) {
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
				"/module/internal/tests/case32/msg.pb.TimestampTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type             any                `json:"type"`
			Format           string             `json:"format"`
			Default          any                `json:"default"`
			Properties       map[string]*schema `json:"properties"`
			Required         []string           `json:"required"`
			Description      string             `json:"description"`
			XConst           any                `json:"x-const"`
			Enum             []any              `json:"enum"`
			Not              *schema            `json:"not"`
			XTimestampLt     string             `json:"x-timestampLt"`
			XTimestampLte    string             `json:"x-timestampLte"`
			XTimestampGt     string             `json:"x-timestampGt"`
			XTimestampGte    string             `json:"x-timestampGte"`
			XTimestampWithin string             `json:"x-timestampWithin"`
			XTimestampLtNow  bool               `json:"x-timestampLtNow"`
			XTimestampGtNow  bool               `json:"x-timestampGtNow"`
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
		assert.Equal(t, "#/definitions/case32.TimestampTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case32.TimestampTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)
		}
		{
			prop, ok := ms.Properties["required"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Contains(t, ms.Required, "required")
		}
		{
			prop, ok := ms.Properties["const"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "2022-11-19T12:00:00Z", prop.XConst)
		}
		{
			prop, ok := ms.Properties["lt"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "2022-11-19T12:00:00Z", prop.XTimestampLt)
		}
		{
			prop, ok := ms.Properties["lte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "2022-11-19T12:00:00Z", prop.XTimestampLte)
		}
		{
			prop, ok := ms.Properties["gt"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "2022-11-19T12:00:00Z", prop.XTimestampGt)
		}
		{
			prop, ok := ms.Properties["gte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "2022-11-19T12:00:00Z", prop.XTimestampGte)
		}
		{
			prop, ok := ms.Properties["within"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "600s", prop.XTimestampWithin)
		}
		{
			prop, ok := ms.Properties["ltNow"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "600s", prop.XTimestampWithin)
			assert.True(t, prop.XTimestampLtNow)
		}
		{
			prop, ok := ms.Properties["gtNow"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "600s", prop.XTimestampWithin)
			assert.True(t, prop.XTimestampGtNow)
		}
		{
			prop, ok := ms.Properties["ltNowOnly"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "", prop.XTimestampWithin)
			assert.False(t, prop.XTimestampLtNow)
		}
		{
			prop, ok := ms.Properties["gtNowOnly"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "date-time", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "", prop.XTimestampWithin)
			assert.False(t, prop.XTimestampGtNow)
		}
	})
}
