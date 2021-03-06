package pubrecdb_test

import (
	"encoding/hex"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
	. "github.com/soapboxsys/ombudslib/pubrecdb"
)

var tst_blk_a *btcutil.Block
var testdb *PublicRecord

// SetupTestDB exports setupTestDB for tests that live inside pubrecdb.
func SetupTestDB(add_rows bool) (*PublicRecord, error) {
	var err error
	if testdb == nil {
		testdb, err = InitDB(getPath(), &chaincfg.MainNetParams)
	}

	if err != nil {
		panic(err)
	}

	testdb.EmptyTables()
	if err = ExecPragma(testdb, false); err != nil {
		panic(err)
	}
	err = testdb.InsertGenesisBlk(chaincfg.MainNetParams.Net)
	if err != nil {
		panic(err)
	}
	if err = ExecPragma(testdb, true); err != nil {
		panic(err)
	}
	if add_rows {
		setupTestInsertBlocks(testdb)
		setupTestInsertBltns(testdb)
		setupTestInsertEndos(testdb)
	}

	return testdb, nil
}

func TestEmptySetupDB(t *testing.T) {
	_, err := SetupTestDB(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupDB(t *testing.T) {
	db, err := SetupTestDB(true)
	if err != nil {
		t.Fatal(err)
	}
	if db == nil {
		t.Fatal("DB not initialized")
	}
	blk, err := db.GetBlockTip()
	if err != nil {
		t.Fatal(err)
	}
	if blk.Head.Hash != "c29afa6a9c333113f24d09368620c1eeb0943c65b92dc647cf80a51610a876d2" {
		t.Fatal(spw(blk))
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

	peg_h := peg.GetStartBlock().Sha()
	zero := *newSha("0000000000000000000000000000000000000000000000000000000000000000")
	ts := time.Unix(123456789, 0)

	// a_blk hash is:
	// [89fc0ab80e1427696ef6019d0a284101eab6f26b091918e8d2226961f1f3cac1]
	a_blk := wire.MsgBlock{
		Header: wire.BlockHeader{
			MerkleRoot: zero,
			PrevBlock:  *peg_h,
			Timestamp:  ts,
			Bits:       0,
			Nonce:      0,
		},
	}
	tst_blk_a = btcutil.NewBlock(&a_blk)
	tst_blk_a.SetHeight(peg.StartHeight + 1)

	// b_blk hash is:
	// [4dab5d688eb366516e5a64afc7bf8536cdeca5b4d5fbc9fce12b71ac6a11a687]
	b_blk := wire.MsgBlock{
		Header: wire.BlockHeader{
			MerkleRoot: zero,
			PrevBlock:  a_blk.BlockSha(),
			Timestamp:  ts,
			Bits:       0,
			Nonce:      0,
		},
	}

	// c_blk hash is:
	// [c29afa6a9c333113f24d09368620c1eeb0943c65b92dc647cf80a51610a876d2]
	c_blk := wire.MsgBlock{
		Header: wire.BlockHeader{
			MerkleRoot: zero,
			PrevBlock:  b_blk.BlockSha(),
			Timestamp:  ts,
			Bits:       0,
			Nonce:      0,
		},
	}

	blocks := []wire.MsgBlock{a_blk, b_blk, c_blk}

	for i, blk := range blocks {
		ublk := btcutil.NewBlock(&blk)
		// Get the height of the blocks in sequence
		ublk.SetHeight(int32(i+1) + peg.StartHeight)
		db.InsertBlockHead(ublk)
	}
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
	// [196d5c0848253dd156740f5db875a78c4f1fcb384104f395c3ebcc241250f8df]
	bltn = fakeUBltn(7)
	var m string = "We can change content #preflight"
	bltn.Wire.Message = &m
	if err, _ := db.InsertBulletin(bltn); err != nil {
		log.Fatal(err)
	}

	// bltn(8) txid is:
	// [d4480f924779f408627766462b97976d9e7afdfeca8a3f890c1fefdd1a5d4fa2]
	bltn = fakeUBltn(8)
	m = "This is another message #preflight"
	bltn.Wire.Message = &m
	if err, _ := db.InsertBulletin(bltn); err != nil {
		log.Fatal(err)
	}

	// btln(9) txid is:
	// [8a73289660100036a10d5c11aa8728ff41d599a2cac2cd4b029f0eea099051ac]
	// It is in tst_blk_a
	bltn = fakeUBltn(9)
	bltn.Block = tst_blk_a
	m = "A bulletin is all he brought #lambs"
	bltn.Wire.Message = &m
	if err, _ := db.InsertBulletin(bltn); err != nil {
		log.Fatal(err)
	}
}

func setupTestInsertEndos(db *PublicRecord) {

	// All endorsements contained within point to bltn(4)
	bid_r, _ := hex.DecodeString("c19fbeacb46e865bfee6db89e9b0a41019079efa305b477d14a35945442e9f45")
	bid := &wire.ShaHash{}
	bid.SetBytes(bid_r)

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

	// endo(6) txid is:
	// [c471a636cc5698cda96fdab9caa946e9c741d0bdda8636a292c12879d6620e01]
	net := chaincfg.MainNetParams
	auth, _ := btcutil.DecodeAddress("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", &net)
	endo = fakeUEndoAuth(7, bid, auth)
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

func fakeUEndoAuth(seed int, bid *wire.ShaHash, auth btcutil.Address) *ombutil.Endorsement {
	wireEndo := fakeWireEndo(seed, bid)
	pegBlk := peg.GetStartBlock()
	authStr := ombutil.Author(auth.String())

	endo := &ombutil.Endorsement{
		Tx:     fakeMsgTx(seed),
		Author: authStr,
		Wire:   wireEndo,
		Block:  pegBlk,
	}
	return endo
}
