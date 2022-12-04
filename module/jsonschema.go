package module

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/yimi-go/protoc-gen-validate-jsonschema/generate"
)

const (
	moduleName  = "pgv-jsonschema"
	moduleParam = "module"
)

type Module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func JsonSchema() pgs.Module { return &Module{ModuleBase: &pgs.ModuleBase{}} }

func (m *Module) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Module) Name() string { return moduleName }

func (m *Module) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	module := m.Parameters().Str(moduleParam)

	for _, f := range targets {
		m.Push(f.Name().String())

		for _, msg := range f.AllMessages() {
			m.Push(msg.Name().String())

			out := m.ctx.OutputPath(f)
			relativeName := strings.TrimPrefix(msg.FullyQualifiedName(), f.FullyQualifiedName())
			out = out.SetExt(fmt.Sprintf("%s.schema.json", relativeName))

			outPath := strings.TrimLeft(strings.ReplaceAll(filepath.ToSlash(out.String()), module, ""), "/")
			m.AddGeneratorFile(outPath, genJsonSchema(m, msg))

			m.Pop()
		}

		m.Pop()
	}

	return m.Artifacts()
}

func genJsonSchema(bc pgs.BuildContext, msg pgs.Message) string {
	sf := generate.Generate(bc, msg)
	bytes, _ := json.MarshalIndent(sf, "", "  ")
	return string(bytes)
}

var _ pgs.Module = (*Module)(nil)
