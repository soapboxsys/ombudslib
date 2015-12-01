package pubrecdb_test

import (
	"os"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	. "github.com/soapboxsys/ombudslib/pubrecdb"
)

func setupTestDB(add_rows bool) (*PublicRecord, error) {

	dbpath := os.Getenv("GOPATH") + "/src/github.com/soapboxsys/ombudslib/pubrecdb/test/test.db"

	db, err := InitDB(dbpath, chaincfg.MainNetParams)
	if err != nil {
		return nil, err
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
	_, err := setupTestDB(true)
	if err != nil {
		t.Fatal(err)
	}
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
