package nullint

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// Scan implements the Scanner interface of the database driver
func (ni *NullInt) Scan(value interface{}) error {

	// if the nullstring is nil, initialize it
	if ni == nil {
		*ni = NullInt{}
	}

	// if the value is nil, reset the data of the nullstring
	if value == nil {

		ni.Int = 0
		ni.IsNull = true
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
	ni.IsNull = false
	ni.Int = sqlInt.Int64
	return nil
}

// Value implements the db driver Valuer interface
func (ni NullInt) Value() (driver.Value, error) {
	if ni.IsNull {
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
		ni.IsNull = nInt.IsNull
		ni.Int = nInt.Int
		return nil

	case int:
		ni.Int = int64(input)
		ni.IsNull = false
		return nil

	case int32:
		ni.Int = int64(input)
		ni.IsNull = false
		return nil

	case int64:
		ni.Int = input
		ni.IsNull = false
		return nil

	default:
		fmt.Printf("%T\n", input)
		fmt.Println(input)
		return fmt.Errorf("wrong type")
	}
}
