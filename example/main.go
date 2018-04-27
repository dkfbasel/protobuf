// Test conversion from custom types
package main

import (
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	mappingFn "github.com/segmentio/go-snakecase"

	startrek "github.com/dkfbasel/protobuf/example/domain"
	"github.com/dkfbasel/protobuf/types/nullstring"
)

func main() {

	fmt.Println("- run example")

	// initialize db connection
	db, err := databaseConnect()

	// could not find database connection
	if err != nil {
		panic(err.Error())
	}

	defer db.Close() // nolint: errcheck

	// adapt the db mapper function to use snake case
	// refer to sqlx documentation for this

	// note: this should be a rather fast function, as it is called every
	// time that named structs are used with sqlx
	db.MapperFunc(mappingFn.Snakecase)

	// create table starfleet and insert some test data
	setupStmt, err := ioutil.ReadFile("sql/setup.sql")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(string(setupStmt))
	if err != nil {
		panic(err.Error())
	}

	// select the uss enterprise from the database
	stmt := `
	SELECT name, no_of_passengers, mission_statement, we_are_leaving_at
	FROM starfleet
	WHERE name = "USS Enterprise"
	ORDER BY id DESC
	LIMIT 1;`

	// initialize a grcp struct that will be filled from the database
	// directly
	var starshipFleetShip startrek.StarfleetShip
	err = db.Get(&starshipFleetShip, stmt)
	if err != nil {
		panic(err.Error())
	}

	// the USS Entripse has a passenger capacity of 1012 persons
	// and is not on a mission right now
	fmt.Println()
	fmt.Println("---- USS Enterprise: Without mission ----")
	fmt.Printf("%+v\n", starshipFleetShip)

	// set a mission for the USS Enterprise
	starshipFleetShip.MissionStatement = &nullstring.NullString{}
	starshipFleetShip.MissionStatement.Text = "Training mission"
	starshipFleetShip.MissionStatement.IsNull = false

	// update the USS Enterprise information in the database
	stmt = `
	INSERT INTO starfleet (name, no_of_passengers, mission_statement, we_are_leaving_at)
	VALUES (:name, :no_of_passengers, :mission_statement, :we_are_leaving_at);`

	_, err = db.NamedExec(stmt, starshipFleetShip)
	if err != nil {
		panic(err.Error())
	}

	// select the uss enterprise from the database agan
	stmt = `
	SELECT name, no_of_passengers, mission_statement, we_are_leaving_at
	FROM starfleet
	WHERE name = "USS Enterprise"
	ORDER BY id DESC
	LIMIT 1;`

	err = db.Get(&starshipFleetShip, stmt)
	if err != nil {
		panic(err.Error())
	}

	// print USS Enterprise with mission
	fmt.Println()
	fmt.Println("--- USS Enterprise: With mission ----")
	fmt.Printf("%+v\n", starshipFleetShip)

}

// databaseConnect handles the database connection and returns a sqlx database
// if the connection is ready
func databaseConnect() (*sqlx.DB, error) {

	// try to initialize database connection
	db, err := sqlx.Connect(
		"mysql",
		"commander:123456@tcp(starfleet_mysql:3306)/startrek?multiStatements=true&parseTime=true")
	iter := 0

	for iter < 100 && err != nil {
		time.Sleep(time.Second)
		db, err = sqlx.Connect(
			"mysql",
			"commander:123456@tcp(starfleet_mysql:3306)/startrek?multiStatements=true&parseTime=true")
	}

	// could not connect to the database
	if err != nil {
		return nil, err
	}

	// try to connect to the database
	err = db.Ping()
	iter = 0
	for iter < 100 && err != nil {
		time.Sleep(time.Second)
		err = db.Ping()
	}

	// could not ping the database
	if err != nil {
		return nil, err
	}

	return db, nil
}
