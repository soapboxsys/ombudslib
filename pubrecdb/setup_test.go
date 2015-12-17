package pubrecdb_test

import (
	"os"
	"path"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	. "github.com/soapboxsys/ombudslib/pubrecdb"
)

// SetupTestDB exports setupTestDB for tests that live inside pubrecdb.
func SetupTestDB(add_rows bool) (*PublicRecord, error) {
	return setupTestDB(add_rows)
}

func setupTestDB(add_rows bool) (*PublicRecord, error) {
	db, err := InitDB(getPath(), &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}

	if add_rows {
		setupTestInsertBlocks(db)
		setupTestInsertBltns(db)
		setupTestInsertEndos(db)
	}

	return db, nil
}

func TestEmptySetupDB(t *testing.T) {
	_, err := setupTestDB(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupDB(t *testing.T) {
	_, err := InitDB(getPath(), &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
}

func getPath() (s string) {
	s = "/src/github.com/soapboxsys/ombudslib/pubrecdb/test/test.db"
	s = path.Join(os.Getenv("GOPATH"), s)
	return s
}

// setupTestInsertBlocks adds an initials set of blocks to the test db that
// other test functions can rely on.
func setupTestInsertBlocks(db *PublicRecord) {
	/*	blocks := []wire.MsgBlock{



		}


		for _, blk range blocks {
			db.insertBlock(blk)
		}
	*/
}

func setupTestInsertBltns(db *PublicRecord) {

}

func setupTestInsertEndos(db *PublicRecord) {

}
