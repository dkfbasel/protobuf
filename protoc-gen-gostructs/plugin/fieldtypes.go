package plugin

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

// fieldtypes is used to specify the types to use in go for specified
// proto field types
func fieldType(field *descriptor.FieldDescriptorProto, prefixes ...string) string {

	name := fieldTypeName(field)

	// add prefixes to the field name
	if len(prefixes) > 0 {
		name = fmt.Sprintf("%s.%s", strings.Join(prefixes, ""), name)
	}

	// nested messages are pointers
	if field.IsMessage() {
		name = fmt.Sprintf("*%s", name)
	}

	// repeated fields are slices
	if field.IsRepeated() {
		return fmt.Sprintf("[]%s", name)
	}

	return name

}

// fieldTypeName will return the name of the given field type
func fieldTypeName(field *descriptor.FieldDescriptorProto) string {

	name := field.GetType().String()

	switch name {

	case "TYPE_DOUBLE":
		return "float64"

	case "TYPE_FLOAT":
		return "float32"

	case "TYPE_INT32":
		return "int32"

	case "TYPE_INT64":
		return "int64"

	case "TYPE_UINT32":
		return "uint32"

	case "TYPE_UINT64":
		return "uint64"

	case "TYPE_SINT32":
		return "int32"

	case "TYPE_SINT64":
		return "int64"

	case "TYPE_FIXED32":
		return "uint32"

	case "TYPE_FIXED64":
		return "uint64"

	case "TYPE_SFIXED32":
		return "int32"

	case "TYPE_SFIXED64":
		return "int64"

	case "TYPE_BOOL":
		return "bool"

	case "TYPE_STRING":
		return "string"

	case "TYPE_BYTES":
		return "[]byte"

	case "TYPE_ENUM":
		return "int"

	case "TYPE_MESSAGE":
		return shortName(field.GetTypeName())

	default:
		return "UNDEFINED:" + name
	}

}

// get the short name of the type
func shortName(name string) string {
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
