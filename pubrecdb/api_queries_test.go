package pubrecdb_test

import (
	"database/sql"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGetBulletin(t *testing.T) {
	db, _ := SetupTestDB(true)

	b, err := db.GetBulletin("I am not in the db")
	if err != sql.ErrNoRows {
		t.Fatalf("Query should error: %s, %s", b, err)
	}

	// See if the fakeWireBltn(3) txid is present
	txid := "73532d0280dc80bd7b8477522d17cd648eae067d5759cd758b0939159d57dfab"
	bltn, err := db.GetBulletin(txid)
	if err != nil {
		t.Fatal(err)
	}

	if bltn.Author != "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy" &&
		bltn.Message != "Blah" &&
		bltn.NumEndos != 0 &&
		bltn.BlockRef.Hash != "0000000000000000036f69604b2f9074571814702400dbb5d5cf6a78fd1dad40" {
		t.Fatal(spew.Sprintf("bltn(3) returned %s", bltn))
	}

	// Check for fakeWireBltn(4)'s presence.
	txid = "c19fbeacb46e865bfee6db89e9b0a41019079efa305b477d14a35945442e9f45"
	bltn, err = db.GetBulletin(txid)
	if err != nil {
		t.Fatal(err)
	}

	// Check that num endos is right
	if bltn.NumEndos != 2 {
		t.Fatal(spew.Sprintf("bltn(4) query returned %s", bltn))
	}

}
