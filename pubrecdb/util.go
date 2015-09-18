package pubrecdb

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/protocol/ombproto"
)

// Writes a bulletin into the sqlite db, runs an insert or update depending on whether
// block hash exists.
func (db *PublicRecord) StoreBulletin(bltn *ombproto.Bulletin) error {

	var err error
	if bltn.Block == nil {
		cmd := `
		INSERT OR REPLACE INTO bulletins 
		(txid, author, board, message, timestamp) VALUES($1, $2, $3, $4, $5)
		`
		_, err = db.conn.Exec(cmd,
			bltn.Txid.String(),
			bltn.Author,
			bltn.Board,
			bltn.Message,
			bltn.Timestamp,
		)
	} else {
		blockstr := bltn.Block.String()
		cmd := `
		INSERT OR REPLACE INTO bulletins 
		(txid, block, author, board, message, timestamp) VALUES($1, $2, $3, $4, $5, $6)
		`
		_, err = db.conn.Exec(cmd,
			bltn.Txid.String(),
			blockstr,
			bltn.Author,
			bltn.Board,
			bltn.Message,
			bltn.Timestamp,
		)
	}
	if err != nil {
		return err
	}

	return nil
}

func makeBlockRecord(blk *btcutil.Block) *BlockRecord {
	sha, _ := blk.Sha()
	head := blk.MsgBlock().Header
	return &BlockRecord{
		Hash:      sha,
		PrevHash:  &head.PrevBlock,
		Height:    uint64(blk.Height()),
		Timestamp: head.Timestamp.Unix(),
	}
}

// Returns a getblocks msg that requests the best chain.
func (db *PublicRecord) MakeBlockMsg() (wire.Message, error) {

	chaintip, err := db.GetChainTip()
	if err != nil {
		return wire.NewMsgGetBlocks(nil), err
	}

	var curblk *BlockRecord = chaintip
	msg := wire.NewMsgGetBlocks(curblk.Hash)

	heights := []int{}
	step, start := 1, 0
	for i := int(chaintip.Height); i > 0; i -= step {
		// Push last 10 indices first
		if start >= 10 {
			step *= 2
		}
		heights = append(heights, i)
		start++
	}
	heights = append(heights, 0)

	for _, h := range heights {

		var err error
		curblk, err := db.getBlkAtHeight(h)
		if err != nil {
			return nil, err
		}
		msg.AddBlockLocatorHash(curblk.Hash)
	}

	return msg, nil
}

// Randomly returns a BlockBecord at height
func (db *PublicRecord) getBlkAtHeight(height int) (*BlockRecord, error) {
	cmd := `
	SELECT hash, prevhash, height, timestamp FROM blocks WHERE height = $1
	ORDER BY RANDOM()
	LIMIT 1
	`

	row := db.conn.QueryRow(cmd, height)

	blkrec, err := scanBlkRec(row)
	if err != nil {
		return nil, err
	}
	return blkrec, nil
}
