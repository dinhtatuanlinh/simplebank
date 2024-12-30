package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"simplebank/util"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

//var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config file")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())

	//var err error
	//testDb, err = pgxpool.New(context.Background(), dbSource)
	//if err != nil {
	//	log.Fatal("Cannot connect to db:", err)
	//}
	//testQueries = New(testDb)
	//os.Exit(m.Run())
}
