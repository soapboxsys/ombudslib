package pubrecdb_test

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/ombwire"
)

// TestBlockHeadInsert tries to insert a <- b and then c which points nowhere
// and should fail.
func TestBlockHeadInsert(t *testing.T) {
	db, _ := setupTestDB(false)

	bogus_h := wire.ShaHash([wire.HashSize]byte{
		0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F,
	})

	// Test the insertion of a genesis block
	a := chaincfg.MainNetParams.GenesisBlock
	blk := btcutil.NewBlock(a)
	blk.SetHeight(0)

	ok, err := db.InsertBlockHead(blk)
	if !ok && err != nil {
		t.Fatalf("Genesis blk header should fail gracefully\n"+
			"Instead we saw: %s", err)
	}

	// Test num rows in blocks
	cnt, err := db.BlockCount()
	if err != nil {
		t.Fatalf("blk cnt failed with: %s", err)
	}
	if cnt != 1 {
		t.Fatalf("After gen insert blk cnt should be 1. It is: %d", cnt)
	}

	// Test the insertion of a linked block
	b := wire.MsgBlock{
		Header: wire.BlockHeader{
			Version:    1,
			PrevBlock:  a.BlockSha(),
			MerkleRoot: bogus_h,
			Timestamp:  time.Unix(1297000000, 0),
			Bits:       0x1d00ffff,
			Nonce:      0x18aea41a,
		},
	}

	blk = btcutil.NewBlock(&b)
	blk.SetHeight(1)

	ok, err = db.InsertBlockHead(blk)
	if !ok {
		t.Fatalf("Blk b header insert failed: %v:", err)
	}

	// Test the insertion of a block that is not in the chain.
	c := wire.MsgBlock{
		Header: wire.BlockHeader{
			Version:    1,
			PrevBlock:  bogus_h,
			MerkleRoot: bogus_h,
			Timestamp:  time.Unix(1299000000, 0),
			Bits:       0x00,
			Nonce:      0x00,
		},
	}

	blk = btcutil.NewBlock(&c)
	blk.SetHeight(99)
	ok, err = db.InsertBlockHead(blk)
	if ok || err != nil {
		// Sqlite should throw a Foreign Key failure with this text:
		expected_err := fmt.Errorf("sqlite: SQL error: foreign key constraint failed")
		t.Fatalf("Blk c header insert should have failed with: %v"+
			" but got: %v", expected_err, err)
	}

}

// TestBulletinInsert asserts that the sql inserts and accompanying logic that
// inserts bulletins into the public records is functioning properly. After
// inserting it examines the state of the test.db to see if the bulletins (and
// tags) are inserted properly.
func TestBulletinInsert(t *testing.T) {
	db, _ := setupTestDB(false)

	wirebltn := fakeWireBltn()
	genBlk := btcutil.NewBlock(chaincfg.MainNetParams.GenesisBlock)
	auth := ombutil.Author("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy")

	gbltn := &ombutil.Bulletin{
		Tx:     fakeMsgTx(),
		Author: auth,
		Wire:   &wirebltn,
		Block:  genBlk,
	}

	if ok, err := db.InsertBulletin(gbltn); err != nil || !ok {
		t.Fatalf("Inserting bltn(g) failed with: %s", err)
	}

	// Assert the bltn is stored in the record.
	cnt, err := db.BulletinCount()
	if err != nil {
		t.Fatal(err)
	}
	if cnt != 1 {
		t.Fatalf("There should be one bltn in the record not: %d", cnt)
	}

	// Remove location from wirebltn
	wirebltn.Location = nil
	// Change wirebltn message
	var m string = "This is #A #Tagged #Message"
	wirebltn.Message = &m
	lbltn := &ombutil.Bulletin{
		Tx:     fakeMsgTx(),
		Author: auth,
		Wire:   &wirebltn,
		Block:  genBlk,
	}

	if ok, err := db.InsertBulletin(lbltn); err != nil || !ok {
		t.Fatalf("Inserting bltn(l) failed with: %s", err)
	}

	cnt, _ = db.BulletinCount()
	if cnt != 2 {
		t.Fatalf("There should be 2 bltns in the record not: %d", cnt)
	}
}

func TestEndorsementInsert(t *testing.T) {
	db, _ := setupTestDB(false)

	bid := []byte("deadbeefdeadbeefdeadbeef")
	ts := uint64(3242232232)

	wendo := ombwire.Endorsement{
		Timestamp: &ts,
		Bid:       bid,
	}

	// TODO fill out fields
	endo := ombutil.Endorsement{
		Wire: &wendo,
	}

	if err := db.InsertEndorsement(&endo); err != nil {
		t.Fatalf("Insert failed with: %s", err)
	}
}

func fakeWireBltn() ombwire.Bulletin {
	var m string = fmt.Sprintf("Climbing is fun: %d", mrand.Int())
	var ts uint64 = uint64(123741234)
	var l float64 = float64(0.01)

	bltn := ombwire.Bulletin{
		Message:   &m,
		Timestamp: &ts,
		Location: &ombwire.Location{
			Lat: &l,
			Lon: &l,
			H:   &l,
		},
	}
	return bltn
}

// fakeMsgTx creates a MsgTx that has a random component so that it hashes the
// returned tx hashes to a different value everytime the funciton returns a new
// tx.
func fakeMsgTx() *wire.MsgTx {
	msgTx := wire.NewMsgTx()
	txIn := wire.TxIn{
		PreviousOutPoint: wire.OutPoint{
			Hash:  wire.ShaHash{},
			Index: 0xffffffff,
		},
		SignatureScript: []byte{0x04, 0x31, 0xdc, 0x00, 0x1b, 0x01, 0x62},
		Sequence:        0xffffffff,
	}
	// Read some random bytes
	b := make([]byte, 20)
	// Ignore errors
	rand.Read(b)
	txOut := wire.TxOut{
		Value:    5000000000,
		PkScript: b,
	}
	msgTx.AddTxIn(&txIn)
	msgTx.AddTxOut(&txOut)
	return msgTx
}
