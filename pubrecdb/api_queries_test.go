package pubrecdb_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/davecgh/go-spew/spew"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/ombwire"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
)

func newSha(s string) *wire.ShaHash {
	h, _ := wire.NewShaHashFromStr(s)
	return h
}

func TestGetBulletin(t *testing.T) {
	db, _ := SetupTestDB(true)

	txid := newSha("000000000000000000000000000000000000000000000000000000000000000")
	b, err := db.GetBulletin(txid)
	if err != sql.ErrNoRows {
		t.Fatalf("Query should error: %s, %s", b, err)
	}

	// See if the fakeWireBltn(3) txid is present
	txid = newSha("73532d0280dc80bd7b8477522d17cd648eae067d5759cd758b0939159d57dfab")
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
	txid = newSha("c19fbeacb46e865bfee6db89e9b0a41019079efa305b477d14a35945442e9f45")
	bltn, err = db.GetBulletin(txid)
	if err != nil {
		t.Fatal(err)
	}

	// Check that num endos is right
	if bltn.NumEndos != 3 {
		t.Fatal(spew.Sdump("bltn(4) query returned:", bltn))
	}

}

// Tests a bulletins whole round trip, to see if the location values where set
// to null properly.
func TestNilLocBulletin(t *testing.T) {
	db, _ := SetupTestDB(true)
	msg := "This bltn has no location"
	ts := uint64(1234567890)

	noloc := ombwire.Bulletin{
		Timestamp: &ts,
		Message:   &msg,
	}
	bltn := ombutil.Bulletin{
		Tx:     fakeMsgTx(10),
		Block:  peg.GetStartBlock(),
		Author: ombutil.Author("1asdfasdfasfdafads"),
		Wire:   &noloc,
	}
	if err, ok := db.InsertBulletin(&bltn); err != nil || !ok {
		t.Fatal(err)
	}

	txsha := bltn.Tx.TxSha()
	jsonBltn, err := db.GetBulletin(&txsha)
	if err != nil {
		t.Fatal(err)
	}

	if jsonBltn.Location != nil {
		t.Fatalf(spw(jsonBltn))
	}

	b, err := json.Marshal(jsonBltn)
	if err != nil {
		t.Fatal(err)
	}

	s := fmt.Sprintf("%s", b)
	if !strings.Contains(s, `"loc":null,`) {
		t.Fatalf("Output json does not comform to spec:\n %s", spw(jsonBltn))
	}
}

func TestGetTag(t *testing.T) {
	db, _ := SetupTestDB(true)

	page, err := db.GetTag(ombutil.Tag("#wistful"))
	if err != nil {
		t.Fatal(err)
	}

	if len(page.Bulletins) != 0 {
		t.Fatal("Page should be empty")
	}

	page, err = db.GetTag(ombutil.Tag("#preflight"))
	if err != nil {
		t.Fatal(err)
	}

	if len(page.Bulletins) != 2 {
		t.Fatal("Page should have two bulletins in it")
	}
}

func TestGetRange(t *testing.T) {
	db, _ := SetupTestDB(true)

	// Height of start is pegStart + 3
	start := newSha("c29afa6a9c333113f24d09368620c1eeb0943c65b92dc647cf80a51610a876d2")
	// Stop is the hash of the blk before the peg blk.
	stop := newSha("000000000000000002dfcd5cd05cd4f80d792e51ecdc5942cd6cec1365b22a2d")

	page, err := db.QueryRange(start, stop)
	if err != nil {
		t.Fatal(err)
	}

	if len(page.Bulletins) != 5 {
		t.Fatalf("Query failed: %s\n", spw(page))
	}
}

func TestGetAuthorResp(t *testing.T) {
	db, _ := SetupTestDB(true)

	net := chaincfg.MainNetParams
	auth, _ := btcutil.DecodeAddress("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", &net)

	authResp, err := db.GetAuthor(auth)
	if err != nil {
		t.Fatal(err)
	}

	if len(authResp.Bulletins) != 5 || authResp.Summary.LastBlkTs != 1451606601 ||
		len(authResp.Endorsements) != 1 {
		t.Fatalf(spw(authResp))
	}
}

func TestGetNearbyBltns(t *testing.T) {
	db, _ := SetupTestDB(true)

	b, err := db.GetNearbyBltns(45.0, 44.0, 20000)
	if err != nil {
		t.Fatal(err)
	}

	if len(b) != 5 {
		t.Fatal(spw(b))
	}
}

func TestGetMostEndorsedBltns(t *testing.T) {
	db, _ := SetupTestDB(true)

	bltns, err := db.GetMostEndorsedBltns(100)
	if err != nil {
		t.Fatal(err)
	}

	if len(bltns) != 1 {
		t.Fatalf(spw(bltns))
	}

}

func spw(t interface{}) string {
	return spew.Sdump(t)
}
