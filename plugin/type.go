package plugin

import (
	"fmt"
	"strings"
)

func (p *Plugin) extractType(t string) string {
	if strings.HasPrefix(t, "."+p.file.GetPackage()+".") {
		return strings.TrimPrefix(t, "."+p.file.GetPackage()+".")
	}

	return ""
}

func (p *Plugin) importedType(t string) string {

	// FIXME: do not deal with imported types
	return ""

	if t == "" {
		return ""
	}

	obj := p.ObjectNamed(t)

	objType := strings.TrimPrefix(strings.Join(obj.TypeName(), ""), "."+obj.File().GetPackage()+".")

	for i, v := range p.importedList {
		if v == string(obj.GoImportPath()) {
			return fmt.Sprintf("google_protobuf%d.%s", i+1, objType)
		}
	}

	p.importedList = append(p.importedList, string(obj.GoImportPath()))

	return fmt.Sprintf("google_protobuf%d.%s", len(p.importedList), objType)
}

func (p *Plugin) getType(t string) (string, bool) {
	if v := p.extractType(t); v != "" {
		return v, false
	} else if v := p.importedType(t); v != "" {
		return v, true
	}

	return "", false
}

