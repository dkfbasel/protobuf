// Test conversion from custom types
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	startrek "github.com/dkfbasel/protobuf/example/domain"
	"github.com/dkfbasel/protobuf/types/nullstring"
)

func main() {

	// initialize database connection
	db, err := sqlx.Connect("mysql", "commander:123456@tcp(0.0.0.0:3333)/startrek?multiStatements=true&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close() // nolint: errcheck

	// create table starfleet and insert data
	setupStmt, err := ioutil.ReadFile("sql/setup.sql")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(string(setupStmt))
	if err != nil {
		panic(err.Error())
	}

	// initialize a grcp struct
	starshipFleetShip := startrek.StarfleetShip{}

	// create an alias containing db tags
	databaseAlias := startrek.StarfleetShipAlias(starshipFleetShip)

	// select the uss enterprise from the database
	stmt := `
	SELECT name, passengers, mission, departure_time_of_ship
	FROM starfleet
	WHERE name = "USS Enterprise"
	ORDER BY id DESC
	LIMIT 1;`

	err = db.Get(&databaseAlias, stmt)
	if err != nil {
		panic(err.Error())
	}

	// the USS Entripse has a passenger capacity of 1012 persons
	// and is not on a mission right now
	starshipFleetShip = startrek.StarfleetShip(databaseAlias)
	fmt.Println()
	fmt.Println("---- USS Enterprise: Without mission ----")
	fmt.Printf("%+v\n", starshipFleetShip)

	// set a mission for the USS Enterprise
	starshipFleetShip.Mission = &nullstring.NullString{}
	starshipFleetShip.Mission.Text = "Training mission"
	starshipFleetShip.Mission.IsNull = false

	// save the USS Enterprise in the database with the new mission
	databaseAlias = startrek.StarfleetShipAlias(starshipFleetShip)
	stmt = `
	INSERT INTO starfleet (name, passengers, mission, departure_time_of_ship)
	VALUES (:name, :passengers, :mission, :departure_time_of_ship);`
	_, err = db.NamedExec(stmt, databaseAlias)
	if err != nil {
		panic(err.Error())
	}

	// reinitialize the ship
	starshipFleetShip = startrek.StarfleetShip{}

	// create an alias containing db tags
	databaseAlias = startrek.StarfleetShipAlias(starshipFleetShip)

	// select the uss enterprise from the database
	stmt = `
	SELECT name, passengers, mission, departure_time_of_ship
	FROM starfleet
	WHERE name = "USS Enterprise"
	ORDER BY id DESC
	LIMIT 1;`

	err = db.Get(&databaseAlias, stmt)
	if err != nil {
		panic(err.Error())
	}

	// print USS Enterprise with mission
	starshipFleetShip = startrek.StarfleetShip(databaseAlias)
	fmt.Println()
	fmt.Println("--- USS Enterprise: With mission ----")
	fmt.Printf("%+v\n", databaseAlias)

	// remove the table
	clearStmt, err := ioutil.ReadFile("sql/clear.sql")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(string(clearStmt))
	if err != nil {
		panic(err.Error())
	}
}
