package timestamp

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time returns a golang time object
func (ts *Timestamp) Time() time.Time {
	if ts.IsNull {
		return time.Time{}
	}
	return time.Unix(0, ts.Milliseconds*1000*1000)
}

// Scan implements the Scanner interface of the database driver
func (ts *Timestamp) Scan(value interface{}) error {

	// initialize timestamp if pointer is nil
	if ts == nil {
		*ts = Timestamp{}
	}

	dbTime, isNotNull := value.(time.Time)

	if isNotNull {
		ts.Milliseconds = dbTime.UnixNano() / 1000 / 1000
		ts.IsNull = false
		return nil
	}
	ts.Milliseconds = 0
	ts.IsNull = true
	return nil
}

// Value implements the db driver Valuer interface
func (ts Timestamp) Value() (driver.Value, error) {
	if ts.IsNull {
		return nil, nil
	}
	return time.Unix(0, ts.Milliseconds*1000*1000), nil
}

// ImplementsGraphQLType is required by the graphql custom scalar interface
// this defines the name used in the schema to declare a null time type
func (ts *Timestamp) ImplementsGraphQLType(name string) bool {
	return name == "Date"
}

// UnmarshalGraphQL is required by the graphql custom scalar interface
// this wraps the null time
func (ts *Timestamp) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {

	case Timestamp:
		time := Timestamp(input)
		ts.IsNull = time.IsNull
		ts.Milliseconds = time.Milliseconds
		return nil

	case time.Time:
		time := &Timestamp{}
		time.Milliseconds = input.UnixNano() / 1000 / 1000
		time.IsNull = false
		return nil

	case string:

		// try to parse the information as date
		timepoint, err := time.Parse("2006-01-02", input)

		if err != nil {
			return fmt.Errorf("format for time must be yyyy-mm-dd")
		}

		ts.Milliseconds = timepoint.UnixNano() / 1000 / 1000
		ts.IsNull = false
		return nil

	default:
		fmt.Printf("%T\n", input)
		fmt.Println(input)
		return fmt.Errorf("wrong type")
	}
}
