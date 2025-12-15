package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/syralon/coconut/internal/protoc"
)

var VERSION = "v0.0.1"

func main() {
	v := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if v != nil && *v {
		_, file := path.Split(strings.ReplaceAll(os.Args[0], "\\", "/"))
		fmt.Printf("%s %s", strings.Split(file, ".")[0], VERSION)
		os.Exit(0)
	}
	builder := protoc.NewErrorBuilder(VERSION)
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		return builder.Build(plugin)
	})
}
