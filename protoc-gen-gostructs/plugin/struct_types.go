package plugin

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type messageInfo struct {
	Message *generator.Descriptor
	Path    string
}

// generateStructs will generate custom structs with tags
func generateStructs(p *plugin, file *generator.FileDescriptor,
	comments map[string]*descriptor.SourceCodeInfo_Location, pathIndex map[string]string) {

	// now handle all messages and create respective custom structs
	for _, message := range file.Messages() {

		path := pathIndex[message.GetName()]
		comment := comments[path].GetLeadingComments()

		// get the name of the struct
		structName := message.GetName()

		// print a default comment if no comment was specified
		if comment == "" {
			p.printComment(structName, "...")

		} else {
			// print the comment that was specified, replacing the struct name
			p.printComment(strings.Replace(comment, *message.Name, structName, -1))
		}

		// print the struct name
		p.Printf("type %s struct{", structName)
		p.In()

		// print the message fields
		generateMessageFields(p, file, pathIndex, message, comments)

		p.Out()
		p.P("}\n")

	}

}

// generateMessageFields will generate the fields for the given message
func generateMessageFields(p *plugin, file *generator.FileDescriptor,
	pathIndex map[string]string, message *generator.Descriptor,
	comments map[string]*descriptor.SourceCodeInfo_Location) {

	// iterate through all struct fields
	fields := message.GetField()
	for fIndex, field := range fields {

		// create the path to the field comment
		// 4 is for messages, 2 for fields
		path := fmt.Sprintf("%s,2,%d", pathIndex[message.GetName()], fIndex)

		// get the comments for this field
		commentline := comments[path]

		// print leading comments (if any)
		p.printComment(commentline.GetLeadingComments())

		// get the trailing comment
		trailing := commentline.GetTrailingComments()

		// print fields without any trailing comments
		if trailing == "" {
			p.P(generator.CamelCase(field.GetName()), " ", fieldType(field))
			continue
		}

		// parse the trailing comment
		dta := parseTags(trailing)

		// handle non message type fields
		if field.IsMessage() == false {
			// print the field with comments and tags
			p.P(generator.CamelCase(field.GetName()), " ", fieldType(field), dta.Output)
			continue
		}

		// handle nested message types
		if dta.IsEmbedded() {

			if field.IsRepeated() {
				p.Fail("embedded structs cannot be repeated: ", field.GetName())
			}

			p.P(fieldTypeName(field), dta.Output)
			continue
		}

		// print the field with comments and tags
		p.P(generator.CamelCase(field.GetName()), " ", fieldType(field), dta.Output)

	}

}
