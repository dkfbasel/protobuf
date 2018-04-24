package nullstring

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements the Scanner interface of the database driver
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.Text = ""
		ns.IsNull = true
		return nil
	}
	ns.IsNull = false
	ns.Text = value.(string)
	return nil
}

// Value implements the db driver Valuer interface
func (ns NullString) Value() (driver.Value, error) {
	if ns.IsNull {
		return nil, nil
	}
	return ns.Text, nil
}

// ImplementsGraphQLType is required by the graphql custom scalar interface
// this defines the name used in the schema to declare a null time type
func (ns *NullString) ImplementsGraphQLType(name string) bool {
	return name == "String"
}

// UnmarshalGraphQL is required by the graphql custom scalar interface
// this wraps the null string
func (ns *NullString) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {

	case NullString:
		nString := NullString(input)
		ns.IsNull = nString.IsNull
		ns.Text = nString.Text
		return nil

	case string:
		ns.Text = input
		ns.IsNull = false
		return nil

	default:
		fmt.Printf("%T\n", input)
		fmt.Println(input)
		return fmt.Errorf("wrong type")
	}
}
