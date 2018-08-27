package main

import (
	"github.com/gogo/protobuf/vanity/command"
	"github.com/infobloxopen/protoc-gen-atlas-validate/plugin"
)

func main() {
	plugin := &plugin.Plugin{}
	response := command.GeneratePlugin(command.Read(), plugin, ".pb.atlas.validate.go")
	command.Write(response)
}
