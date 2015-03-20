package pubrecdb

import "testing"

func TestSetupDB(t *testing.T) {
	_, err := SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}

}
