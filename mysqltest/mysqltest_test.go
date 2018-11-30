package mysqltest

import "testing"

func Test(t *testing.T) {
	db, err := PrepareTestFixtures("fixtures/mysql", "root", "1234567890", []string{})
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
}
