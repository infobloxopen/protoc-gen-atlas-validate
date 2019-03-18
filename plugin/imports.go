package plugin

import (
	"fmt"
	"path"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	httpPkgPath      = "net/http"
	jsonPkgPath      = "encoding/json"
	ctxPkgPath       = "context"
	fmtPkgPath       = "fmt"
	ioutilPkgPath    = "io/ioutil"
	metadataPkgPath  = "google.golang.org/grpc/metadata"
	bytesPkgPath     = "bytes"
	runtimePkgPath   = "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
	gwruntimePkgPath = "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func (p *Plugin) initPluginImports(g *generator.Generator) {

	pi := NewPluginImports(g)

	for _, v := range []string{
		httpPkgPath,
		jsonPkgPath,
		ctxPkgPath,
		fmtPkgPath,
		ioutilPkgPath,
		metadataPkgPath,
		bytesPkgPath,
		runtimePkgPath,
		gwruntimePkgPath,
	} {
		pi.AddImport(v, "")
	}

	p.pluginImports = pi
}

type importPkg struct {
	path string
	used bool
	name string
}

func (p *importPkg) IsUsed() bool {
	return p.used
}

func (p *importPkg) Use() string {
	p.used = true
	return p.name
}

func (p *importPkg) Name() string {
	return p.name
}

func NewPluginImports(g *generator.Generator) *pluginImports {
	return &pluginImports{
		generator: g,
		pkgs:      make([]*importPkg, 0),
		pkgsMap:   make(map[string]*importPkg),
	}
}

type pluginImports struct {
	generator *generator.Generator
	pkgs      []*importPkg
	pkgsMap   map[string]*importPkg
}

func (p *pluginImports) AddImport(pkg string, pkgName string) *importPkg {

	var (
		pkgBaseName = path.Base(pkg)
		pkgCount    = 0
	)

	if imp, ok := p.pkgsMap[pkg]; ok {
		return imp
	}

	if pkgName == "" {
		pkgName = pkgBaseName
	}

	for _, v := range p.pkgs {
		if v.name == pkgName {
			pkgCount++
			pkgName = fmt.Sprintf("%s%d", pkgBaseName, pkgCount)
		}
	}

	imp := &importPkg{
		path: pkg,
		used: false,
		name: pkgName,
	}

	p.pkgs = append(p.pkgs, imp)
	p.pkgsMap[pkg] = imp

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

func (p *pluginImports) NewImport(pkg string) *importPkg {
	if _, ok := p.pkgsMap[pkg]; !ok {
		p.generator.Fail(`unable to find entry for import path `, pkg, `: include import in InitPluginImports`)
	}

	return p.pkgsMap[pkg]
}

func (p *Plugin) isLocal(o generator.Object) bool {
	return p.DefaultPackageName(o) == ""
}

func (p *Plugin) objectNamed(name string) generator.Object {
	obj := p.ObjectNamed(name)

	if !p.isLocal(obj) {
		p.AddImport(string(obj.GoImportPath()), "").Use()
	}

	return obj
}

func (p *Plugin) objectFieldNamed(o generator.Object, t string, f string) generator.Object {
	return p.objectNamed(o.File().GetMessage(t).GetFieldDescriptor(f).GetTypeName())
}
