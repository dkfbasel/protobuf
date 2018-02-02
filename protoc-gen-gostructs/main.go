package main

import (
	"io/ioutil"
	"os"
	"strings"

	"bitbucket.org/dkfbasel/dev.grpc-tags/protoc-gen-gostructs/plugin"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func main() {

	gen := generator.New()

	// read the proto file passed from protoc via stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		gen.Error(err, "could not read input")
	}

	err = proto.Unmarshal(data, gen.Request)
	if err != nil {
		gen.Error(err, "could not parse proto definition")
	}

	if len(gen.Request.FileToGenerate) == 0 {
		gen.Fail("no output files to generate")
	}

	// get command line paramenters
	gen.CommandLineParameters(gen.Request.GetParameter())

	if gen.Param["import"] == "" {
		gen.Fail("import parameter must be specified to access the proto structs")
	}

	// wrap all descriptor and file descriptors
	gen.WrapTypes()

	// set the pacakge name to be used
	gen.SetPackageNames()

	// build map of types
	gen.BuildTypeNameMap()

	// generate a plugin to handle all files
	gen.GeneratePlugin(plugin.New())

	// go through all input files and define the name of the output
	for i := 0; i < len(gen.Response.File); i++ {

		// modify the content after generation
		// NOTE: ideally this should be adapted in the template, however this
		// does currently not seem to be possible with gogo/protobuf
		var newContent = gen.Response.File[i].GetContent()

		// replace the generator name
		newContent = strings.Replace(newContent, "by protoc-gen-gogo", "by protoc-gen-gotags", -1)

		// remove proto imports
		newContent = strings.Replace(newContent, "import proto \"github.com/gogo/protobuf/proto\"\n", "", -1)

		// remove math package
		newContent = strings.Replace(newContent, "import math \"math\"\n", "", -1)

		// remove underscore definitions
		newContent = strings.Replace(newContent, "var _ = proto.Marshal\n", "", -1)
		newContent = strings.Replace(newContent, "var _ = math.Inf\n", "", -1)

		gen.Response.File[i].Content = &newContent

		newFileName := strings.Replace(*gen.Response.File[i].Name, ".pb.go", ".structs.pb.go", -1)

		gen.Response.File[i].Name = proto.String(newFileName)

		// gen.Error(fmt.Errorf("tmp"), gen.Response.File[i].GetName())
	}

	// return the return of the plugin to protoc via stdout
	data, err = proto.Marshal(gen.Response)
	if err != nil {
		gen.Error(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		gen.Error(err, "failed to pass output to protoc")
	}

}
