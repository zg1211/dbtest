package mysqltest

import (
	"testing"

	"github.com/zg1211/dbtest/mysqltest/schema"
)

func Test(t *testing.T) {
	db, err := PrepareTestFixtures("fixtures/mysql", "root", "1234567890", []string{schema.DropTable(), schema.CreateTable()})
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
}
