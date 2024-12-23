package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // use _ to use the side effect of the package only
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:2423@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var globalDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	globalDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	testQueries = New(globalDB)
	os.Exit((m.Run()))
}
