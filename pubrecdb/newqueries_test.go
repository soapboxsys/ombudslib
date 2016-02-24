package pubrecdb_test

import (
	"database/sql"
	"testing"
)

func TestGetBlock(t *testing.T) {
	db, _ := SetupTestDB(true)

	a_hash := tst_blk_a.Sha()

	a_blk, err := db.GetBlock(a_hash)
	if err != nil {
		t.Fatal(err)
	}

	if a_blk.Head.Hash != a_hash.String() && len(a_blk.Bulletins) != 1 {
		t.Fatal(spw(a_blk))
	}

	zero := newSha("0000000000000000000000000000000000000000000000000000000000000000")
	b, err := db.GetBlock(zero)
	if err != sql.ErrNoRows {
		t.Log("Query should have returned nothing, instead saw:")
		t.Fatalf(spw(b))
	}
}

func TestGetBlockTip(t *testing.T) {
	db, _ := SetupTestDB(true)

	tipHash := "c29afa6a9c333113f24d09368620c1eeb0943c65b92dc647cf80a51610a876d2"

	tipBlk, err := db.GetBlockTip()
	if err != nil {
		t.Fatal(err)
	}

	if tipBlk.Head.Hash != tipHash {
		t.Fatalf(spw(tipBlk))
	}
}
