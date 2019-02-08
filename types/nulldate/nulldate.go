package nulldate

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

// IsNull will return if the current timestamp is null
func (dt *NullDate) IsNull() bool {

	if dt == nil {
		return true
	}

	// we use IsNotNull instead of IsNull to make sure that a date is always
	// initialized as null value
	return dt.IsNotNull == false
}

// Set will set the date to the given time
func (dt *NullDate) Set(value string) {

	if dt == nil {
		*dt = NullDate{}
	}

	if value == "" {
		dt.SetNull()
		return
	}

	dt.Date = value
	dt.IsNotNull = true
}

// SetNull will clear the date
func (dt *NullDate) SetNull() {

	if dt == nil {
		return
	}

	dt.Date = ""
	dt.IsNotNull = false
}

// Time returns a golang time object of the date
func (dt *NullDate) Time() time.Time {

	if dt.IsNull() || dt.Date == "" {
		return time.Time{}
	}

	asTime, err := time.Parse("2006-01-02", dt.Date)
	if err != nil {
		return time.Time{}
	}

	return asTime
}

// Scan implements the Scanner interface of the database driver
func (dt *NullDate) Scan(value interface{}) error {

	// initialize timestamp if pointer is nil
	if dt == nil {
		*dt = NullDate{}
	}

	switch input := value.(type) {
	case time.Time:
		dt.Date = input.Format("2006-01-02")
		dt.IsNotNull = true
		return nil

	case string:
		dt.Date = input
		dt.IsNotNull = true
		return nil

	case nil:
		dt.SetNull()
		return nil

	default:
		return fmt.Errorf("unkown type for NullDate: %T", input)
	}
}

// Value implements the db driver Valuer interface
func (dt *NullDate) Value() (driver.Value, error) {

	if dt.IsNull() || dt.Date == "" {
		return nil, nil
	}

	return dt.Date, nil
}

// ImplementsGraphQLType is required by the graphql custom scalar interface
// this defines the name used in the schema to declare a null time type
func (dt *NullDate) ImplementsGraphQLType(name string) bool {
	return name == "Date"
}

// UnmarshalGraphQL is required by the graphql custom scalar interface
// this wraps the null date
func (dt *NullDate) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {

	case NullDate:
		dt.IsNotNull = input.IsNotNull
		dt.Date = input.Date
		return nil

	case time.Time:
		dt.Date = input.Format("2006-01-02")
		dt.IsNotNull = true
		return nil

	case string:
		// try to parse the information as date to make sure it is valid
		_, err := time.Parse("2006-01-02", input)
		if err != nil {
			return fmt.Errorf("date must be valid and in the format yyyy-mm-dd")
		}

		dt.Set(input)
		return nil

	case nil:
		dt.SetNull()
		return nil

	default:
		return fmt.Errorf("unkown type for NullDate: %T", input)
	}
}

// UnmarshalJSON is used to convert the json representation into a null date
func (dt *NullDate) UnmarshalJSON(input []byte) error {
	// trim the leading and trailing quotes from the timestamp
	cleanInput := bytes.Trim(input, "\"")
	asString := string(cleanInput)
	dt.Set(asString)
	return nil
}

// MarshalJSON will return the content as json value, this is also called
// by graphql to generate the response
func (dt NullDate) MarshalJSON() ([]byte, error) {

	if dt.IsNull() {
		return []byte("null"), nil
	}

	// format the timestamp in iso compatible time format
	formatted := fmt.Sprintf("\"%s\"", dt.Date)

	return []byte(formatted), nil
}
