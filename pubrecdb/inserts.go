package pubrecdb

import "github.com/btcsuite/btcutil"

var (
	insertBlockHeadSql string = `
		INSERT INTO blocks (hash, prevhash, height, timestamp, version, merkleroot, difficulty, nonce) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
)

func prepareInserts(db *PublicRecord) error {
	var err error

	db.insertBlockHead, err = db.conn.Prepare(insertBlockHeadSql)
	if err != nil {
		return err
	}

	return nil
}

// InsertBlock does exactly what you expect it does. If the block is already in
// the record then an error is thrown.
func (db *PublicRecord) InsertBlockHead(blk *btcutil.Block) error {
	h := blk.MsgBlock().Header

	hash := h.BlockSha().String()
	prevhash := h.PrevBlock.String()
	height := blk.Height()
	timestamp := uint32(h.Timestamp.Unix())
	v := h.Version
	m := h.MerkleRoot.String()
	d := h.Bits
	n := h.Nonce

	_, err := db.insertBlockHead.Exec(hash, prevhash, height, timestamp, v, m, d, n)
	if err != nil {
		return err
	}

	return nil
}
