package plugin

import (
	"fmt"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

// generateConvertFnCustom will generate converting functions for custom types
// to proto types
func generateConvertFnCustom(p *plugin, file *generator.FileDescriptor,
	comments map[string]*descriptor.SourceCodeInfo_Location, pathIndex map[string]string) {

	for mIndex, message := range file.Messages() {

		p.P("// Convert will convert the proto struct to a custom type with")
		p.P("// tags and embedded structs if specified")
		p.Printf("func (from *%s) Convert() (*protodef.%s, error) {", message.GetName(), message.GetName())
		p.In()

		p.P("var err error\n")
		p.P("if from == nil {")
		p.In()
		p.P("return nil, nil")
		p.Out()
		p.P("}\n")

		p.P("to := protodef.", message.Name, "{}")

		// iterate through all struct fields
		fields := message.GetField()
		for fIndex, field := range fields {

			p.P()

			// get the name of the field
			fieldNameProto := field.GetName()

			// get the type of the field
			ftype := fieldType(field)

			// get the type of the proto field
			ftypeProto := fieldType(field, "protodef")

			// create the path to the field comment
			// 4 is for messages, 2 for fields
			path := fmt.Sprintf("4,%d,2,%d", mIndex, fIndex)

			// get the comments for this field
			commentline := comments[path]

			// parse any field tags
			dta := parseTags(commentline.GetTrailingComments())

			if field.IsScalar() || field.IsString() {

				// handle all non repeated fields
				if field.IsRepeated() == false {
					// convert the field
					p.Printf("to.%s = from.%s", fieldNameProto, fieldNameProto)
					continue
				}

				// check if the slice is nil
				p.Printf("if from.%s != nil {", fieldNameProto)
				p.In()

				// generate a slice of the same type with the same size
				// note: ftype will already generate slice of types
				p.Printf("to.%s = make(%s, len(from.%s))", fieldNameProto, ftype, fieldNameProto)

				// go rough all items in the from struct
				p.Printf("for i := range from.%s {", fieldNameProto)
				p.In()

				// directly assign scalar fields
				p.Printf("to.%s[i] = from.%s[i]", fieldNameProto, fieldNameProto)

				p.Out()
				p.P("}")

				p.Out()
				p.P("}")

			}

			// handle all non repeated fields
			if field.IsMessage() {

				// handle embbeded structs
				// convert custom embedded struct and assign nested message in proto struct
				if dta.IsEmbedded() {

					// embedded fields cannot be repeated
					if field.IsRepeated() {
						p.Fail("embedded structs cannot be repeated: ", field.GetName())
						continue
					}

					messageName := fieldTypeName(field)
					p.Printf("to.%s, err = from.%s.Convert()", fieldNameProto, messageName)
					p.PrintReturnErr()
					continue
				}

				// handle fields that are not repeated
				if field.IsRepeated() == false {

					// convert nested structs
					p.Printf("to.%s, err = from.%s.Convert()", fieldNameProto, fieldNameProto)
					p.PrintReturnErr()
					continue

				}

				// handle repeated fields

				// check if the slice is nil
				p.Printf("if from.%s != nil {", fieldNameProto)
				p.In()

				// generate a slice of the same type with the same size
				// note: ftype will already generate a pointer to a slice

				p.Printf("to.%s = make(%s, len(from.%s))", fieldNameProto, ftypeProto, fieldNameProto)

				// go rough all items in the from struct
				p.Printf("for i := range from.%s {", fieldNameProto)
				p.In()

				// directly assign scalar fields
				p.Printf("to.%s[i], err = from.%s[i].Convert()", fieldNameProto, fieldNameProto)
				p.PrintReturnErr()

				p.Out()
				p.P("}")

				p.Out()
				p.P("}")

			}

		}

		p.P()
		p.P("return &to, err")

		p.Out()
		p.P("}")
		p.P()

	}

}

// generateConvertFnProto will generate converting functions from proto types to
// custom types
func generateConvertFnProto(p *plugin, file *generator.FileDescriptor,
	comments map[string]*descriptor.SourceCodeInfo_Location, pathIndex map[string]string) {

	for mIndex, message := range file.Messages() {

		messageName := message.GetName()

		p.Printf("// Convert%s will convert the proto optimized struct into", messageName)
		p.P("// an easier to use go structs")
		p.Printf("func Convert%s(from *protodef.%s) (*%s, error) {", messageName, messageName, messageName)
		p.In()

		p.P("var err error\n")
		p.P("if from == nil {")
		p.In()
		p.P("return nil, nil")
		p.Out()
		p.P("}\n")

		p.P("to := ", messageName, "{}")

		// iterate through all struct fields
		fields := message.GetField()
		for fIndex, field := range fields {

			p.P()

			// get the name of the field
			fieldNameProto := field.GetName()

			// get the type of the field
			ftype := fieldType(field)

			// create the path to the field comment
			// 4 is for messages, 2 for fields
			path := fmt.Sprintf("4,%d,2,%d", mIndex, fIndex)

			// get the comments for this field
			commentline := comments[path]

			// parse any field tags
			dta := parseTags(commentline.GetTrailingComments())

			if field.IsScalar() || field.IsString() {

				// handle all non repeated fields
				if field.IsRepeated() == false {
					// convert the field
					p.Printf("to.%s = from.Get%s()", fieldNameProto, fieldNameProto)
					continue
				}

				// check if the slice is nil
				p.Printf("if from.%s != nil {", fieldNameProto)
				p.In()

				// generate a slice of the same type with the same size
				// note: ftype will already generate slice of types
				p.Printf("to.%s = make(%s, len(from.%s))", fieldNameProto, ftype, fieldNameProto)

				// go rough all items in the from struct
				p.Printf("for i := range from.%s {", fieldNameProto)
				p.In()

				// directly assign scalar fields
				p.Printf("to.%s[i] = from.%s[i]", fieldNameProto, fieldNameProto)

				p.Out()
				p.P("}")

				p.Out()
				p.P("}")

			}

			// handle all non repeated fields
			if field.IsMessage() {

				// get the message name of the nested field
				fieldMessageName := fieldTypeName(field)

				// handle embbeded structs
				// convert custom embedded struct and assign nested message in proto struct
				if dta.IsEmbedded() {

					// embedded fields cannot be repeated
					if field.IsRepeated() {
						p.Fail("embedded structs cannot be repeated: ", field.GetName())
						continue
					}

					p.Printf("tmp, err := Convert%s(from.Get%s())", fieldMessageName, fieldNameProto)
					p.PrintReturnErr()
					p.Printf("to.%s = *tmp", fieldMessageName)
					continue
				}

				// handle fields that are not repeated
				if field.IsRepeated() == false {

					// convert nested structs
					p.Printf("to.%s, err = Convert%s(from.Get%s())", fieldNameProto, fieldMessageName, fieldNameProto)
					p.PrintReturnErr()
					continue

				}

				// handle repeated fields

				// check if the slice is nil
				p.Printf("if from.Get%s() != nil {", fieldNameProto)
				p.In()

				// generate a slice of the same type with the same size
				// note: ftype will already generate a pointer to a slice
				p.Printf("to.%s = make(%s, len(from.%s))", fieldNameProto, ftype, fieldNameProto)

				// go rough all items in the from struct
				p.Printf("for i := range from.%s {", fieldNameProto)
				p.In()

				// directly assign scalar fields
				p.Printf("to.%s[i], err = Convert%s(from.%s[i])", fieldNameProto, fieldMessageName, fieldNameProto)
				p.PrintReturnErr()

				p.Out()
				p.P("}")

				p.Out()
				p.P("}")

			}

		}

		p.P()
		p.P("return &to, err")

		p.Out()
		p.P("}")
		p.P()

	}

}
