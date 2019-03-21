package plugin

import (
	"fmt"
	"path"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	bytesPkgPath  = "bytes"
	ctxPkgPath    = "context"
	fmtPkgPath    = "fmt"
	httpPkgPath   = "net/http"
	ioutilPkgPath = "io/ioutil"
	jsonPkgPath   = "encoding/json"

	metadataPkgPath  = "google.golang.org/grpc/metadata"
	gwruntimePkgPath = "github.com/grpc-ecosystem/grpc-gateway/runtime"

	runtimePkgPath = "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
)

var wkt = map[string]bool{
	// ptypes
	".google.protobuf.Timestamp": true,
	".google.protobuf.Duration":  true,
	".google.protobuf.Empty":     true,
	".google.protobuf.Any":       true,
	".google.protobuf.Struct":    true,

	// nillable values
	".google.protobuf.StringValue": true,
	".google.protobuf.BytesValue":  true,
	".google.protobuf.Int32Value":  true,
	".google.protobuf.UInt32Value": true,
	".google.protobuf.Int64Value":  true,
	".google.protobuf.UInt64Value": true,
	".google.protobuf.FloatValue":  true,
	".google.protobuf.DoubleValue": true,
	".google.protobuf.BoolValue":   true,
}

// initPluginImports function initializes plugin imports with set of prior-known
// packages.
func (p *Plugin) initPluginImports(g *generator.Generator) {

	pi := NewPluginImports(g)

	for _, v := range []string{

		// std packages
		bytesPkgPath,
		ctxPkgPath,
		fmtPkgPath,
		httpPkgPath,
		ioutilPkgPath,
		jsonPkgPath,

		// external packages
		metadataPkgPath,
		gwruntimePkgPath,

		// local packages
		runtimePkgPath,
	} {
		pi.AddImport(v)
	}

	p.pluginImports = pi
}

// importPkg structure represents one import atom that is a path, imported name and flag
// that indicates whether this import package is used or not.
type importPkg struct {
	path string
	used bool
	name string
}

// IsUsed function returns flag that indiciates whether this package was used or not.
func (p *importPkg) IsUsed() bool {
	return p.used
}

// Use function marks import package as used and returns its import name.
func (p *importPkg) Use() string {
	p.used = true
	return p.name
}

// NewPluginImports returns initialized pluginImport object.
func NewPluginImports(g *generator.Generator) *pluginImports {
	return &pluginImports{
		generator: g,
		pkgs:      make([]*importPkg, 0),
	}
}

// pluginImports structure represents a set of package imports.
type pluginImports struct {
	generator *generator.Generator
	pkgs      []*importPkg
}

// getPkg function seeks for a package among added packages.
func (p *pluginImports) getPkg(pkgPath string) (*importPkg, bool) {
	for _, v := range p.pkgs {
		if pkgPath == v.path {
			return v, true
		}
	}

	return nil, false
}

// AddImport function adds a new unused import package to a list of imports.
// It uses basename of the package path as package name, in case of conflicts
// it renders package name as <basename><unique-number> e. g. runtime1.
func (p *pluginImports) AddImport(pkgPath string) *importPkg {

	var (
		pkgBaseName = path.Base(pkgPath)
		pkgName     = pkgBaseName
		pkgCount    = 0
	)

	if imp, ok := p.getPkg(pkgPath); ok {
		return imp
	}

	for _, v := range p.pkgs {
		// since all items are added in order we cannot meet basename2 before
		// we encounter basename1.
		if v.name == pkgName {
			pkgCount++
			pkgName = fmt.Sprintf("%s%d", pkgBaseName, pkgCount)
		}
	}

	imp := &importPkg{
		path: pkgPath,
		used: false,
		name: pkgName,
	}

	p.pkgs = append(p.pkgs, imp)
	return imp
}

func (p *pluginImports) GenerateImports(file *generator.FileDescriptor) {
	for _, v := range p.pkgs {
		if !v.used {
			continue
		}
		p.generator.PrintImport(generator.GoPackageName(v.name), generator.GoImportPath(v.path))
	}
}

// Import functions locates predefined import and fails if later does not exist.
func (p *pluginImports) Import(pkgPath string) *importPkg {
	var (
		imp *importPkg
		ok  bool
	)

	if imp, ok = p.getPkg(pkgPath); !ok {
		p.generator.Fail(`unable to find entry for import path `, pkgPath, `: include import in InitPluginImports`)
	}

	return imp
}

func (p *Plugin) isLocal(o generator.Object) bool {
	return p.DefaultPackageName(o) == ""
}

func (p *Plugin) isWKT(t string) bool {
	return wkt[t]
}

func (p *Plugin) objectNamed(name string) generator.Object {
	obj := p.ObjectNamed(name)

	if !p.isLocal(obj) {
		p.AddImport(string(obj.GoImportPath())).Use()
	}

	return obj
}
