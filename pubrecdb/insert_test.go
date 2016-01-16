package pubrecdb_test

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/ombwire"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
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

	// Test the insertion of the peg block
	blk := peg.GetStartBlock()
	a := blk.MsgBlock()

	err, ok := db.InsertBlockHead(blk)
	if !ok && err != nil {
		t.Fatalf("Peg blk header should fail gracefully\n"+
			"Instead we saw: %s", err)
	}

	// Test num rows in blocks
	cnt, err := db.BlockCount()
	if err != nil {
		t.Fatalf("blk cnt failed with: %s", err)
	}
	if cnt != 1 {
		t.Fatalf("After peg insert blk cnt should be 1. It is: %d", cnt)
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

	err, ok = db.InsertBlockHead(blk)
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
	err, ok = db.InsertBlockHead(blk)
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

	wirebltn := fakeWireBltn(1)
	pegBlk := peg.GetStartBlock()
	auth := ombutil.Author("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy")

	gbltn := &ombutil.Bulletin{
		Tx:     fakeMsgTx(1),
		Author: auth,
		Wire:   &wirebltn,
		Block:  pegBlk,
	}

	if err, ok := db.InsertBulletin(gbltn); err != nil || !ok {
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
		Tx:     fakeMsgTx(2),
		Author: auth,
		Wire:   &wirebltn,
		Block:  pegBlk,
	}

	if err, ok := db.InsertBulletin(lbltn); err != nil || !ok {
		t.Fatalf("Inserting bltn(l) failed with: %s", err)
	}

	cnt, _ = db.BulletinCount()
	if cnt != 2 {
		t.Fatalf("There should be 2 bltns in the record not: %d", cnt)
	}
}

func TestEndorsementInsert(t *testing.T) {
	db, _ := setupTestDB(false)

	bid := []byte("adfdsawerjklgroastbeefgroastbeef")
	ts := uint64(3242232232)

	wendo := ombwire.Endorsement{
		Timestamp: &ts,
		Bid:       bid,
	}

	endo := ombutil.Endorsement{
		Wire:   &wendo,
		Block:  peg.GetStartBlock(),
		Author: ombutil.Author("1asfde238jfha32hydsa"),
		Tx:     fakeMsgTx(2),
	}

	if err, ok := db.InsertEndorsement(&endo); err != nil || !ok {
		t.Fatalf("Insert failed with: %s: %s", err, ok)
	}

	if c, _ := db.EndoCount(); c != 1 {
		t.Fatalf("Insert did not add record: %d", c)
	}

	// Run the insert again
	if err, ok := db.InsertEndorsement(&endo); err == nil && ok {
		t.Fatalf("Insert with duplicate txid should have failed")
	}

	// Try insert with bad Bid hash value
	endo.Tx = fakeMsgTx(3)
	endo.Wire.Bid = []byte("roastbeef")

	if err, ok := db.InsertEndorsement(&endo); err == nil && ok {
		t.Fatalf("Insert w/ bad bid should have failed")
	}
}

// fakeWireBltn lets us create a random or deterministic bltn as needed for
// various test cases.
func fakeWireBltn(nonce int) ombwire.Bulletin {
	var m string = fmt.Sprintf("This is a unique bltn[%d]", nonce)
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

func fakeWireEndo(seed int, bid *wire.ShaHash) *ombwire.Endorsement {
	ts := uint64(1234567890 + seed)

	endo := &ombwire.Endorsement{
		Bid:       bid.Bytes(),
		Timestamp: &ts,
	}

	return endo
}

// fakeMsgTx creates a MsgTx that has a random seed so that it hashes the
// returned tx hashes to a different value when a differnt nonce is provided.
func fakeMsgTx(nonce int) *wire.MsgTx {
	msgTx := wire.NewMsgTx()
	txIn := wire.TxIn{
		PreviousOutPoint: wire.OutPoint{
			Hash:  wire.ShaHash{},
			Index: 0xffffffff,
		},
		SignatureScript: []byte{0x04, 0x31, 0xdc, 0x00, 0x1b, 0x01, 0x62},
		Sequence:        0xffffffff,
	}
	// Place the nonce
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(nonce))
	txOut := wire.TxOut{
		Value:    5000000000,
		PkScript: b,
	}
	msgTx.AddTxIn(&txIn)
	msgTx.AddTxOut(&txOut)
	return msgTx
}
