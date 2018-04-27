package nullstring

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// Scan implements the Scanner interface of the database driver
func (ns *NullString) Scan(value interface{}) error {

	// if the nullstring is nil, initialize it
	if ns == nil {
		*ns = NullString{}
	}

	// if the value is nil, reset the data of the nullstring
	if value == nil {

		ns.Text = ""
		ns.IsNull = true
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

	ns.IsNull = false
	ns.Text = sqlString.String
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
