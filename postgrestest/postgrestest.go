package postgrestest

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	testfixtures "gopkg.in/testfixtures.v2"
)

var (
	dbHost string
	dbPort string
	dbName string
)

func PrepareTestFixtures(fixturesPath, testUser, testPwd string, mustExecs []string) (*sql.DB, error) {
	creatorDB, err := getCreatorDB()
	if err != nil {
		return nil, err
	}
	defer creatorDB.Close()

	db := sqlx.NewDb(creatorDB, "postgres")
	tx := db.MustBegin()
	for _, mustExec := range mustExecs {
		tx.MustExec(mustExec)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	testerDB, err := getTesterDB(testUser, testPwd)
	if err != nil {
		return nil, err
	}

	fixtures, err := testfixtures.NewFolder(creatorDB, &testfixtures.PostgreSQL{}, fixturesPath)
	if err != nil {
		return nil, err
	}

	if err := fixtures.Load(); err != nil {
		return nil, err
	}

	return testerDB, nil
}

// dbHost = "127.0.0.1"
// dbPort = "3306"
// dbName = "demo_test"
// dbCreatorUser = "root"
// dbCreatorUserPwd = "1234567890"
func getCreatorDB() (*sql.DB, error) {
	dbHost = os.Getenv("DBTEST_HOST")
	dbPort = os.Getenv("DBTEST_PORT")
	dbName = os.Getenv("DBTEST_DB")
	dbCreatorUser := os.Getenv("DBTEST_CREATOR_USER")
	dbCreatorUserPwd := os.Getenv("DBTEST_CREATOR_PASSWORD")

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbCreatorUser, dbCreatorUserPwd, dbName)
	return sql.Open("postgres", dbinfo)
}

func getTesterDB(testUser, testPwd string) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, testUser, testPwd, dbName)
	return sql.Open("postgres", dbinfo)
}
