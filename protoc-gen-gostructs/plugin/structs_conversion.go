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

		p.P("// Proto will convert our custom type to a struct that can be")
		p.P("// serialized as protobuf")
		p.Printf("func (from *%s) Proto() (*protodef.%s, error) {", message.GetName(), message.GetName())
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
					p.Printf("to.%s, err = from.%s.Proto()", fieldNameProto, messageName)
					p.PrintReturnErr()
					continue
				}

				// handle fields that are not repeated
				if field.IsRepeated() == false {

					// convert nested structs
					p.Printf("to.%s, err = from.%s.Proto()", fieldNameProto, fieldNameProto)
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
				p.Printf("to.%s[i], err = from.%s[i].Proto()", fieldNameProto, fieldNameProto)
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

		p.Printf("// FromProto%s will convert the proto optimized struct into", messageName)
		p.P("// an easier to use go structs")
		p.Printf("func (to *%s) FromProto(from *protodef.%s) error {", messageName, messageName)
		p.In()

		p.P("var err error\n")
		p.P("if from == nil {")
		p.In()
		p.P("return nil")
		p.Out()
		p.P("}\n")

		p.P("toTmp := ", messageName, "{}")

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
					p.Printf("toTmp.%s = from.Get%s()", fieldNameProto, fieldNameProto)
					continue
				}

				// check if the slice is nil
				p.Printf("if from.%s != nil {", fieldNameProto)
				p.In()

				// generate a slice of the same type with the same size
				// note: ftype will already generate slice of types
				p.Printf("toTmp.%s = make(%s, len(from.%s))", fieldNameProto, ftype, fieldNameProto)

				// go rough all items in the from struct
				p.Printf("for i := range from.%s {", fieldNameProto)
				p.In()

				// directly assign scalar fields
				p.Printf("toTmp.%s[i] = from.%s[i]", fieldNameProto, fieldNameProto)

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

					p.Printf("tmp := %s{}", fieldMessageName)
					p.Printf("err = tmp.FromProto(from.Get%s())", fieldNameProto)
					p.PrintReturnErr("return err")
					p.Printf("toTmp.%s = tmp", fieldMessageName)
					continue
				}

				// handle fields that are not repeated
				if field.IsRepeated() == false {

					// convert nested structs
					p.Printf("tmp := %s{}", fieldMessageName)
					p.Printf("err = tmp.FromProto(from.Get%s())", fieldNameProto)
					p.PrintReturnErr("return err")

					p.Printf("toTmp.%s = tmp", fieldNameProto)

					continue

				}

				// handle repeated fields

				// check if the slice is nil
				p.Printf("if from.Get%s() != nil {", fieldNameProto)
				p.In()

				// generate a slice of the same type with the same size
				// note: ftype will already generate a pointer to a slice
				p.Printf("toTmp.%s = make(%s, len(from.%s))", fieldNameProto, ftype, fieldNameProto)

				// go rough all items in the from struct
				p.Printf("for i := range from.%s {", fieldNameProto)
				p.In()

				// directly assign scalar fields
				p.Printf("tmp := %s{}", fieldMessageName)
				p.Printf("err = tmp.FromProto(from.%s[i])", fieldNameProto)
				p.PrintReturnErr("return err")

				p.Printf("toTmp.%s[i] = tmp", fieldNameProto)

				p.Out()
				p.P("}")

				p.Out()
				p.P("}")

			}

		}

		p.P()

		// return the error if not nil
		p.PrintReturnErr("return err")

		// assign the temporary struct to our struct
		// dereferencing of the pointer will make sure that we replace
		// the actual struct in memory
		p.P()
		p.P("*to = toTmp")
		p.P()
		p.P("return nil")

		p.Out()
		p.P("}")
		p.P()

	}

}
