package pubrecdb

import "testing"

func TestWriteROFails(t *testing.T) {
	db, _ := SetupTestDB()

	cmd := `INSERT INTO bulletins (message, board, author) VALUES ('you', 'got', 'pwnd')`

	_, err := db.ExecuteROExpr(cmd)
	if err == nil {
		t.Fatal(err)
	}
}

func TestReadROWorks(t *testing.T) {
	db, _ := SetupTestDB()

	cmd := `SELECT * FROM bulletins`

	_, err := db.ExecuteROExpr(cmd)
	if err != nil {
		t.Fatal(err)
	}
}
