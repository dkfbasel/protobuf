package nullint

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// IsNull will return if the current null int is null
func (ni *NullInt) IsNull() bool {

	if ni == nil {
		return true
	}

	// we use IsNotNull instead of IsNull to make sure that a timestamp is
	// initialized as null value
	return ni.IsNotNull == false
}

// Set will set the null string to the given value
func (ni *NullInt) Set(value int64) {

	if ni == nil {
		*ni = NullInt{}
	}

	ni.Int = value
	ni.IsNotNull = true

}

// SetNull will set the null int to null
func (ni *NullInt) SetNull() {

	if ni == nil {
		*ni = NullInt{}
	}

	ni.Int = 0
	ni.IsNotNull = false

}

// Scan implements the Scanner interface of the database driver
func (ni *NullInt) Scan(value interface{}) error {

	// if the nullstring is nil, initialize it
	if ni == nil {
		*ni = NullInt{}
	}

	// if the value is nil, reset the data of the nullstring
	if value == nil {

		ni.Int = 0
		ni.IsNotNull = false
		return nil

	}

	// create a sql NullInt64 to use the not exported convertAssign-method
	// of the golang sql package
	sqlInt := sql.NullInt64{}

	// scan the value, using the sql package
	err := sqlInt.Scan(value)
	if err != nil {
		return err
	}
	ni.IsNotNull = true
	ni.Int = sqlInt.Int64
	return nil
}

// Value implements the db driver Valuer interface
func (ni NullInt) Value() (driver.Value, error) {
	if ni.IsNull() {
		return nil, nil
	}
	return ni.Int, nil
}

// ImplementsGraphQLType is required by the graphql custom scalar interface
// this defines the name used in the schema to declare a null time type
func (ni *NullInt) ImplementsGraphQLType(name string) bool {
	return name == "Int"
}

// UnmarshalGraphQL is required by the graphql custom scalar interface
// this wraps the null integer
func (ni *NullInt) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {

	case NullInt:
		nInt := NullInt(input)
		ni.IsNotNull = nInt.IsNotNull
		ni.Int = nInt.Int
		return nil

	case int:
		ni.Int = int64(input)
		ni.IsNotNull = true
		return nil

	case int32:
		ni.Int = int64(input)
		ni.IsNotNull = true
		return nil

	case int64:
		ni.Int = input
		ni.IsNotNull = true
		return nil

	default:
		fmt.Printf("%T\n", input)
		fmt.Println(input)
		return fmt.Errorf("wrong type")
	}
}
