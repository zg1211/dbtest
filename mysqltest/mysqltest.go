package mysqltest

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

	db := sqlx.NewDb(creatorDB, "mysql")
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

	fixtures, err := testfixtures.NewFolder(creatorDB, &testfixtures.MySQL{}, fixturesPath)
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

	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbCreatorUser, dbCreatorUserPwd, dbHost, dbPort, dbName)
	return sql.Open("mysql", dbinfo)
}

func getTesterDB(testUser, testPwd string) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", testUser, testPwd, dbHost, dbPort, dbName)
	return sql.Open("mysql", dbinfo)
}
