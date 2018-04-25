// Test conversion from custom types
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	startrek "github.com/dkfbasel/protobuf/example/domain"
)

func main() {

	// initialize database connection
	db, err := sqlx.Connect("mysql", "commander:123456@tcp(0.0.0.0:3333)/startrek")
	if err != nil {
		log.Fatalln(err)
	}

	// create table starfleet and insert data
	setupStmt, err := ioutil.ReadFile("sql/setup.sql")
	if err != nil {
		panic(err.Error())
	}

	// initialize a startfleet ship struct
	starshipFleetShip := startrek.StarfleetShip{}

	databaseAlias := startrek.StarfleetShipAlias(starshipFleetShip)

	stmt := `
	SELECT name, passengers, mission, departure_time_of_ship
	FROM starfleet
	WHERE name = "USS Enterprise"`

	err = db.Get(&databaseAlias, stmt)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("%v", databaseAlias)

	_, err = db.Exec(string(setupStmt))
	if err != nil {
		panic(err.Error())
	}

	defer db.Close() // nolint: errcheck
}
