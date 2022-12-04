package case31

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

func TestDuration(t *testing.T) {
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
				"/module/internal/tests/case31/msg.pb.DurationTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type         any                `json:"type"`
			Format       string             `json:"format"`
			Default      any                `json:"default"`
			Properties   map[string]*schema `json:"properties"`
			Required     []string           `json:"required"`
			Description  string             `json:"description"`
			XConst       any                `json:"x-const"`
			Enum         []any              `json:"enum"`
			Not          *schema            `json:"not"`
			XDurationLt  string             `json:"x-durationLt"`
			XDurationLte string             `json:"x-durationLte"`
			XDurationGt  string             `json:"x-durationGt"`
			XDurationGte string             `json:"x-durationGte"`
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
		assert.Equal(t, "#/definitions/case31.DurationTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case31.DurationTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}
		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)
		}
		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)
		}
		{
			prop, ok := ms.Properties["required"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Contains(t, ms.Required, "required")
		}
		{
			prop, ok := ms.Properties["in"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, []any{"60s", "1.500s"}, prop.Enum)
		}
		{
			prop, ok := ms.Properties["notIn"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.NotNil(t, prop.Not)
			assert.Equal(t, []any{"60s", "1.500s"}, prop.Not.Enum)
		}
		{
			prop, ok := ms.Properties["const"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "60s", prop.XConst)
		}
		{
			prop, ok := ms.Properties["lt"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "60s", prop.XDurationLt)
		}
		{
			prop, ok := ms.Properties["lte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "60s", prop.XDurationLte)
		}
		{
			prop, ok := ms.Properties["gt"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "60s", prop.XDurationGt)
		}
		{
			prop, ok := ms.Properties["gte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "string", prop.Type)
			assert.Equal(t, "duration", prop.Format)
			assert.Nil(t, prop.Default)

			assert.Equal(t, "60s", prop.XDurationGte)
		}
	})
}
