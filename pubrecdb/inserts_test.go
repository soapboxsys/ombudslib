package pubrecdb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombutil"
)

// TestBlockHeadInsert tries to insert a <- b and then c which points nowhere
// and should fail.
func TestBlockHeadInserts(t *testing.T) {
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
	if ok {
		t.Fatalf("Genesis blk header should fail\n"+
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
	if ok {
		// Sqlite should throw a Foreign Key failure with this text:
		expected_err := fmt.Errorf("sqlite: SQL error: foreign key constraint failed")
		t.Fatalf("Blk c header insert should have failed with: %v"+
			" but got: %v", expected_err, err)
	}

}

// TestBulletinInserts asserts that the sql inserts and accompanying logic that
// inserts bulletins into the public records is functioning properly. After
// inserting it examines the state of the test.db to see if the bulletins (and
// tags) are inserted properly.
func TestBulletinInserts(t *testing.T) {
	db, _ := setupTestDB(false)

	tx := wire.NewMsgTx()

	bltn, err := ombutil.NewBulletin(tx, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("Creating the bltn failed: %s", err)
	}

	// TODO add bulletin with more than five tags

	if err := db.InsertBulletin(bltn); err != nil {
		t.Fatalf("Inserting bltn a failed with: %s", err)
	}

	// TODO assert bltn is stored in ledger.
}
