package base

import (
	"bytes"
	"os"
	"testing"

	pgs "github.com/lyft/protoc-gen-star"
	"github.com/spf13/afero"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yimi-go/protoc-gen-validate-jsonschema/module"
)

func FrameTest(t *testing.T, binFile string, assertion func(t *testing.T, fs afero.Fs, res *bytes.Buffer)) {
	req, err := os.Open(binFile)
	if err != nil {
		t.Fatal(err)
	}

	fs := afero.NewMemMapFs()
	res := &bytes.Buffer{}

	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.Init(
		pgs.DebugMode(),
		pgs.SupportedFeatures(&optional),
		pgs.ProtocInput(req),  // use the pre-generated request
		pgs.ProtocOutput(res), // capture CodeGeneratorResponse
		pgs.FileSystem(fs),    // capture any custom files written directly to disk
	).RegisterModule(module.JsonSchema()).Render()

	// check res and the fs for output
	assertion(t, fs, res)
}
