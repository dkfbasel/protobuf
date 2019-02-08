package nullstring

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
)

// IsNull will return if the current string is null
func (ns *NullString) IsNull() bool {
	// we use IsNotNull instead of IsNull to make sure that a timestamp is
	// initialized as null value
	return ns.IsNotNull == false
}

// Set will set the null string to the given value
func (ns *NullString) Set(value string) {

	if ns == nil {
		*ns = NullString{}
	}

	ns.Text = value
	ns.IsNotNull = true

}

// SetNull will set the nullstring to null
func (ns *NullString) SetNull() {

	if ns == nil {
		*ns = NullString{}
	}

	ns.Text = ""
	ns.IsNotNull = false

}

// Scan implements the Scanner interface of the database driver
func (ns *NullString) Scan(value interface{}) error {

	// if the nullstring is nil, initialize it
	if ns == nil {
		*ns = NullString{}
	}

	// if the value is nil, reset the data of the nullstring
	if value == nil {
		ns.Text = ""
		ns.IsNotNull = false
		return nil
	}

	// create a sql NullString to use the not exported convertAssign-method
	// of the golang sql package
	sqlString := sql.NullString{}

	// scan the value, using the sql package
	err := sqlString.Scan(value)
	if err != nil {
		return err
	}

	ns.IsNotNull = true
	ns.Text = sqlString.String
	return nil
}

// Value implements the db driver Valuer interface
func (ns NullString) Value() (driver.Value, error) {
	if ns.IsNull() {
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
		ns.IsNotNull = input.IsNotNull
		ns.Text = input.Text
		return nil

	case string:
		ns.Text = input
		ns.IsNotNull = true
		return nil

	default:
		fmt.Printf("%T\n", input)
		fmt.Println(input)
		return fmt.Errorf("wrong type")
	}
}

// UnmarshalJSON is required to parse json input to string value
func (ns *NullString) UnmarshalJSON(input []byte) error {
	ns.Set(string(input))
	return nil
}

// MarshalJSON will return the content as json value, this is also called
// by graphql to generate the response
func (ns *NullString) MarshalJSON() ([]byte, error) {

	if ns.IsNull() {
		return []byte("null"), nil
	}

	// escape quotes inline to guarantee that the result is a single string
	return []byte(strconv.Quote(ns.Text)), nil
}
