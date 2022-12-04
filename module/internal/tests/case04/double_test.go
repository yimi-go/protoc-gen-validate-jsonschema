package case04

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

func TestDouble(t *testing.T) {
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
				"/module/internal/tests/case04/msg.pb.DoubleTestMsg.schema.json",
			file.GetName())
		t.Logf("\n%s", file.GetContent())
		type schema struct {
			Type             string             `json:"type"`
			Format           string             `json:"format"`
			Default          *float64           `json:"default"`
			Properties       map[string]*schema `json:"properties"`
			Description      string             `json:"description"`
			XConst           *float64           `json:"x-const"`
			Maximum          *float64           `json:"maximum"`
			ExclusiveMaximum bool               `json:"ExclusiveMaximum"`
			Minimum          *float64           `json:"minimum"`
			ExclusiveMinimum bool               `json:"exclusiveMinimum"`
			Enum             []float64          `json:"enum"`
			Not              *schema            `json:"not"`
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
		assert.Equal(t, "#/definitions/case04.DoubleTestMsg", sf.Ref)
		assert.NotEmpty(t, sf.Definitions)
		ms, ok := sf.Definitions["case04.DoubleTestMsg"]
		if !ok {
			t.Fatal("missing root msg def")
		}
		assert.Equal(t, "object", ms.Type)

		{
			prop, ok := ms.Properties["desc"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.Equal(t, "aaa\n\nbbb", prop.Description)
		}

		{
			prop, ok := ms.Properties["noRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.Nil(t, prop.XConst)
			assert.Nil(t, prop.Maximum)
			assert.False(t, prop.ExclusiveMaximum)
			assert.Nil(t, prop.Minimum)
			assert.False(t, prop.ExclusiveMinimum)
		}

		{
			prop, ok := ms.Properties["blankRule"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.Equal(t, ms.Properties["noRule"], prop)
		}

		{
			prop, ok := ms.Properties["const"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.XConst)
			assert.Equal(t, float64(1), *prop.XConst)
		}

		{
			prop, ok := ms.Properties["lt"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Maximum)
			assert.Equal(t, float64(10), *prop.Maximum)
			assert.True(t, prop.ExclusiveMaximum)
		}

		{
			prop, ok := ms.Properties["lte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Maximum)
			assert.Equal(t, float64(10), *prop.Maximum)
			assert.False(t, prop.ExclusiveMaximum)
		}

		{
			prop, ok := ms.Properties["ltLtLte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Maximum)
			assert.Equal(t, float64(9), *prop.Maximum)
			assert.True(t, prop.ExclusiveMaximum)
		}

		{
			prop, ok := ms.Properties["ltEqLte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Maximum)
			assert.Equal(t, float64(10), *prop.Maximum)
			assert.True(t, prop.ExclusiveMaximum)
		}

		{
			prop, ok := ms.Properties["ltGtLte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Maximum)
			assert.Equal(t, float64(10), *prop.Maximum)
			assert.False(t, prop.ExclusiveMaximum)
		}

		{
			prop, ok := ms.Properties["gt"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Minimum)
			assert.Equal(t, float64(10), *prop.Minimum)
			assert.True(t, prop.ExclusiveMinimum)
		}

		{
			prop, ok := ms.Properties["gte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Minimum)
			assert.Equal(t, float64(10), *prop.Minimum)
			assert.False(t, prop.ExclusiveMinimum)
		}

		{
			prop, ok := ms.Properties["gtLtGte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Minimum)
			assert.Equal(t, float64(10), *prop.Minimum)
			assert.False(t, prop.ExclusiveMinimum)
		}

		{
			prop, ok := ms.Properties["gtEqGte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Minimum)
			assert.Equal(t, float64(10), *prop.Minimum)
			assert.True(t, prop.ExclusiveMinimum)
		}

		{
			prop, ok := ms.Properties["gtGtGte"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.NotNil(t, prop.Minimum)
			assert.Equal(t, float64(11), *prop.Minimum)
			assert.True(t, prop.ExclusiveMinimum)
		}

		{
			prop, ok := ms.Properties["in"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.Equal(t, float64(0), *prop.Default)
			assert.Equal(t, []float64{1, 2, 3, 4, 5}, prop.Enum)
		}

		{
			prop, ok := ms.Properties["notIn"]
			assert.True(t, ok)
			assert.NotNil(t, prop)
			assert.Equal(t, "number", prop.Type)
			assert.Equal(t, "double", prop.Format)
			assert.NotNil(t, prop.Default)
			assert.NotNil(t, prop.Not)
			assert.Equal(t, []float64{7, 8, 9}, prop.Not.Enum)
		}
	})
}
