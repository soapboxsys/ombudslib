package pubrecdb_test

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

// Determines if the delete works and fails when it is supposed to.
func TestDeleteBlockTip(t *testing.T) {
	db, _ := setupTestDB(false)

	bogus_h := wire.ShaHash([wire.HashSize]byte{
		0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F,
	})

	pegSha := peg.GetStartBlock().Sha()

	blk := wire.MsgBlock{
		Header: wire.BlockHeader{
			Version:    1,
			PrevBlock:  *pegSha,
			MerkleRoot: bogus_h,
			Timestamp:  time.Unix(1297000000, 0),
			Bits:       0x1d00ffff,
			Nonce:      0x18aea41a,
		},
	}

	a := btcutil.NewBlock(&blk)
	a.SetHeight(2)

	err, ok := db.InsertBlockHead(a)
	if !ok || err != nil {
		t.Fatalf("Blk(a) header insert failed: %v, %v", ok, err)
	}

	// Try to delete the genesis block.
	err, ok = db.DeleteBlockTip(pegSha)
	if ok || err != pubrecdb.ErrBlockNotTip {
		t.Fatalf("Blk(pegBlk) delete needs to fail: %v, %v", ok, err)
	}

	// Delete the current block tip
	err, ok = db.DeleteBlockTip(a.Sha())
	if !ok || err != nil {
		t.Fatalf("Blk(a) delete failed: %v, %v", ok, err)
	}
}
