package plugin

import (
	"strconv"
	"strings"

	descriptor "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func extractComments(file *generator.FileDescriptor) map[string]*descriptor.SourceCodeInfo_Location {

	comments := make(map[string]*descriptor.SourceCodeInfo_Location)

	// go through all locations in the file
	for _, loc := range file.GetSourceCodeInfo().GetLocation() {

		// initalize the separate path parts
		var p []string

		// combine the path to nodes
		for _, n := range loc.Path {
			p = append(p, strconv.Itoa(int(n)))
		}

		// save the path to the comment
		comments[strings.Join(p, ",")] = loc
	}

	return comments

}

// printComment will print the given comment (and split newlines)
func (p *plugin) printComment(comments ...string) {

	comment := strings.Join(comments, " ")

	if comment == "" {
		return
	}

	text := strings.TrimSuffix(comment, "\n")
	for _, line := range strings.Split(text, "\n") {
		p.P("// ", strings.TrimPrefix(line, " "))
	}

}
