package pubrecdb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/protocol/ombproto"
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

	err := db.InsertBlockHead(blk)
	if err == nil {
		expected_err := fmt.Errorf("sqlite3: column hash is not unique [2067]")
		t.Fatalf("Genesis blk header insert did not fail with: %s\n"+
			"Instead we saw: %s", expected_err, err)
	}
	// TODO test num rows in blocks

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

	err = db.InsertBlockHead(blk)
	if err != nil {
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
	err = db.InsertBlockHead(blk)
	if err == nil {
		// Sqlite will throw a Foreign Key failure with this text:
		expected_err := fmt.Errorf("sqlite: SQL error: foreign key constraint failed")
		t.Fatalf("Blk c header insert should have failed with: %v"+
			" but got: %v", expected_err, err)
	}

	// TODO test num rows in blocks
}

// TestBulletinInserts asserts that the sql inserts and accompanying logic that
// inserts bulletins into the public records is functioning properly. After
// inserting it examines the state of the test.db to see if the bulletins (and
// tags) are inserted properly.
func TestBulletinInserts(t *testing.T) {
	db, _ := setupTestDB(false)

	fakeTxid, _ := wire.NewShaHash([]byte{
		0xF0, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F,
	})
	genHash := chaincfg.MainNetParams.GenesisBlock.BlockSha()

	// TODO use the proper constructor
	bltn := &ombproto.Bulletin{
		// TODO add some tags
		Message:   "This is a test. There are many like it, but it is mine.",
		Author:    "1F1tAaz5x1HUXrCNLbtMDqcw6o5GNn4xqX",
		Timestamp: time.Now(),
		Txid:      fakeTxid,
		Block:     &genHash,
	}

	// TODO add bulletin with more than five tags

	if err := db.InsertBulletin(bltn); err != nil {
		t.Fatalf("Inserting bltn a failed with: %s", err)
	}

	// TODO assert bltn is stored in ledger.
}
