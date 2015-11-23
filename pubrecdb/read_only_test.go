package pubrecdb_test

import "testing"

func TestWriteROFails(t *testing.T) {
	db, _ := setupTestDB(true)

	cmd := `INSERT INTO bulletins (message, board, author) VALUES ('you', 'got', 'pwnd')`

	_, err := db.ExecuteROExpr(cmd)
	if err == nil {
		t.Fatal(err)
	}
}

func TestReadROWorks(t *testing.T) {
	db, _ := setupTestDB(true)

	cmd := `SELECT * FROM bulletins`

	_, err := db.ExecuteROExpr(cmd)
	if err != nil {
		t.Fatal(err)
	}
}
