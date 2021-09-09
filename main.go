package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// response := command.GeneratePlugin(command.Read(), plugin, ".pb.atlas.validate.go")

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var request pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(input, &request)
	if err != nil {
		panic(err)
	}

	opts := protogen.Options{}

	plugin, err := opts.New(&request)
	if err != nil {
		panic(err)
	}

	resp := generate(plugin)
	out, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(out))
}

func generate(plugin *protogen.Plugin) *pluginpb.CodeGeneratorResponse {
	fmt.Fprintf(os.Stderr, "running generate\n")
	return &pluginpb.CodeGeneratorResponse{}
}
