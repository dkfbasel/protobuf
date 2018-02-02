package plugin

import (
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

// CustomPluginImports ...
type CustomPluginImports struct {
	generator *generator.Generator
	imports   map[string]string
}

// NewPluginImports will generate a new custom plugin import struct
func NewPluginImports(generator *generator.Generator) *CustomPluginImports {
	return &CustomPluginImports{
		generator,
		make(map[string]string),
	}
}

// NewImport is used to import a new pkg
func (plg *CustomPluginImports) NewImport(pkg string) generator.Single {

	parts := strings.Split(pkg, ":::")
	if len(parts) == 2 {
		plg.imports[parts[0]] = parts[1]
	}

	return nil
}

// GenerateImports is used to generete the import parts
func (plg *CustomPluginImports) GenerateImports(file *generator.FileDescriptor) {
	for alias, pkg := range plg.imports {
		plg.generator.PrintImport(alias, pkg)
	}
}
