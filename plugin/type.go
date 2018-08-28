package plugin

import (
	"fmt"
	"strings"
)

func (p *Plugin) importedType(t string) string {

	if t == "" {
		return ""
	}

	obj := p.ObjectNamed(t)

	objType := strings.TrimPrefix(strings.Join(obj.TypeName(), ""), "."+obj.File().GetPackage()+".")

	for i, v := range p.imports {
		if v == string(obj.GoImportPath()) {
			return fmt.Sprintf("google_protobuf%d.%s", i+1, objType)
		}
	}

	p.imports = append(p.imports, string(obj.GoImportPath()))
	return fmt.Sprintf("google_protobuf%d.%s", len(p.imports), objType)
}

func (p *Plugin) getProtoType(t string) string {
	if strings.HasPrefix(t, p.pkgPrefix()) {
		return p.trimPkgPrefix(t)
	}

	return ""
}

func (p *Plugin) getGoType(t string) string {
	if strings.HasPrefix(t, p.pkgPrefix()) {
		return strings.Join(strings.Split(p.getProtoType(t), "."), "_")
	} else if !strings.HasPrefix(t, ".") {
		return strings.Join(strings.Split(t, "."), "_")
	}

	return ""
}
