package pubrecdb_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
	. "github.com/soapboxsys/ombudslib/pubrecdb"
)

// SetupTestDB exports setupTestDB for tests that live inside pubrecdb.
func SetupTestDB(add_rows bool) (*PublicRecord, error) {
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
	_, err := SetupTestDB(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupDB(t *testing.T) {
	_, err := SetupTestDB(true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitDB(t *testing.T) {
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

// TODO(nskelsey)
// setupTestInsertBlocks adds an initial set of blocks to the test db that
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

	// bltn txid is
	// [73532d0280dc80bd7b8477522d17cd648eae067d5759cd758b0939159d57dfab]
	bltn := fakeUBltn(3)
	if err, _ := db.InsertBulletin(bltn); err != nil {
		log.Fatal(err)
	}

	// bltn(4) txid is:
	// [c19fbeacb46e865bfee6db89e9b0a41019079efa305b477d14a35945442e9f45]
	bltn = fakeUBltn(4)
	if err, _ := db.InsertBulletin(bltn); err != nil {
		log.Fatal(err)
	}

	// bltn(7) txid is:
	bltn = fakeUBltn(7)
	var m string = "We can change content #preflight"
	bltn.Wire.Message = &m
	if err, _ := db.InsertBulletin(bltn); err != nil {
		log.Fatal(err)
	}

	// bltn(8) txid is:
	bltn = fakeUBltn(8)
	m = "This is another message #preflight"
	bltn.Wire.Message = &m
	if err, _ := db.InsertBulletin(bltn); err != nil {
		log.Fatal(err)
	}
}

func setupTestInsertEndos(db *PublicRecord) {

	// All endorsements contained within point to bltn(4)
	bid, _ := wire.NewShaHashFromStr("c19fbeacb46e865bfee6db89e9b0a41019079efa305b477d14a35945442e9f45")

	// endo(5) txid is:
	// [4bf52e816c845b40f71209e611fc3a1d352526d57f722a4c5fad7d8558611be3]
	endo := fakeUEndo(5, bid)
	if err, _ := db.InsertEndorsement(endo); err != nil {
		log.Fatal(err)
	}

	// endo(6) txid is:
	// [c471a636cc5698cda96fdab9caa946e9c741d0bdda8636a292c12879d6620e01]
	endo = fakeUEndo(6, bid)
	if err, _ := db.InsertEndorsement(endo); err != nil {
		log.Fatal(err)
	}
}

func fakeUBltn(seed int) *ombutil.Bulletin {
	wirebltn := fakeWireBltn(seed)
	pegBlk := peg.GetStartBlock()
	auth := ombutil.Author("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy")

	bltn := &ombutil.Bulletin{
		Tx:     fakeMsgTx(seed),
		Author: auth,
		Wire:   &wirebltn,
		Block:  pegBlk,
	}
	return bltn
}

func fakeUEndo(seed int, bid *wire.ShaHash) *ombutil.Endorsement {
	wireEndo := fakeWireEndo(seed, bid)
	pegBlk := peg.GetStartBlock()
	auth := ombutil.Author("4end0")

	endo := &ombutil.Endorsement{
		Tx:     fakeMsgTx(seed),
		Author: auth,
		Wire:   wireEndo,
		Block:  pegBlk,
	}
	return endo
}
