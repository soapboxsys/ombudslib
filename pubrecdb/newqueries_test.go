package pubrecdb_test

import (
	"database/sql"
	"testing"
	"time"
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

func TestGetStatistics(t *testing.T) {
	db, _ := SetupTestDB(true)

	before := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	after := time.Date(2020, 12, 30, 0, 0, 0, 0, time.UTC)

	stats, err := db.GetStatistics(before, after)
	if err != nil {
		t.Fatal(err)
	}

	if stats.NumBltns != 5 && stats.NumEndos != 3 && stats.NumBlks != 4 {
		t.Fatal(spw(stats))
	}
}
