// Test conversion from custom types
package main

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root@/startrek")
	if err != nil {
		panic(err.Error())
	}

	setupStmt, err := ioutil.ReadFile("sql/setup.sql")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(string(setupStmt))
	if err != nil {
		panic(err.Error())
	}

	defer db.Close() // nolint: errcheck
}
