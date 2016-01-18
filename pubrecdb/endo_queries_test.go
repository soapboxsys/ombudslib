package pubrecdb_test

import (
	"database/sql"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGetEndorsement(t *testing.T) {
	db, _ := SetupTestDB(true)

	zero := newSha("0000000000000000000000000000000000000000000000000000000000000000")
	_, err := db.GetEndorsement(zero)
	if err != sql.ErrNoRows {
		t.Fatalf("Query should return no rows not: %s", err)
	}

	// See if fakeWireEndo(5) is present
	txid := newSha("4bf52e816c845b40f71209e611fc3a1d352526d57f722a4c5fad7d8558611be3")
	e, err := db.GetEndorsement(txid)
	if err != nil {
		t.Fatal(err)
	}

	if e.Timestamp != 1234567895 {
		t.Fatal(spew.Sprintf("%s\n", e))
	}
}

func TestGetEndosByBid(t *testing.T) {
	db, _ := SetupTestDB(true)

	bid := newSha("c19fbeacb46e865bfee6db89e9b0a41019079efa305b477d14a35945442e9f45")
	endos, err := db.GetEndosByBid(bid)
	if err != nil {
		t.Fatal(err)
	}

	if len(endos) != 2 {
		spw(endos)
	}
}
