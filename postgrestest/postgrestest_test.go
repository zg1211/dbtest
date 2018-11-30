package postgrestest

import (
	"testing"

	"github.com/zg1211/dbtest/postgrestest/schema"
)

func Test(t *testing.T) {
	db, err := PrepareTestFixtures("fixtures/postgres", "root", "1234567890", []string{schema.DropTable(), schema.CreateTable()})
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
}
