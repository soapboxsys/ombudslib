package pubrecdb

import (
	"database/sql"

	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/protocol/ombproto"
)

var (
	insertBlockHeadSql string = `
		INSERT INTO blocks (hash, prevhash, height, timestamp, version, merkleroot, difficulty, nonce) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	insertBulletinSql string = `
		INSERT INTO bulletins (txid, block, author, message, timestamp, latitude, longitude, height)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	insertTagSql string = `
		INSERT INTO tags (txid, value) VALUES ($1, $2)
	`

	insertEndoSql string = `
		INSERT INTO endorsements (txid, bid, timestamp, author) VALUES ($1, $2, $3, $4)
	`
)

func prepareInserts(db *PublicRecord) (err error) {
	db.insertBlockHead, err = db.conn.Prepare(insertBlockHeadSql)
	if err != nil {
		return err
	}

	db.insertBulletin, err = db.conn.Prepare(insertBulletinSql)
	if err != nil {
		return err
	}

	db.insertTag, err = db.conn.Prepare(insertTagSql)
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

// InsertBulletin takes a bulletins and inserts it into the pubrecord. A
// ombproto.Bulletin is used here instead of a wire.Bulletin because of the
// high level utilities offered by the ombproto type.
func (db *PublicRecord) InsertBulletin(bltn *ombproto.Bulletin) (err error) {

	// Start a sql transaction to insert the bulletin and the relevant tags.
	var tx *sql.Tx
	if tx, err = db.conn.Begin(); err != nil {
		return err
	}

	// Insert the Bulletin
	txid := bltn.Txid.String()
	blkHash := bltn.Block.String()
	ath := bltn.Author
	msg := bltn.Message
	ts := bltn.Timestamp.Unix()
	// TODO add proper agruments
	lt := 45000000
	lg := 75000000
	ht := 250

	_, err = tx.Stmt(db.insertBulletin).Exec(txid, blkHash, ath, msg, ts, lt, lg, ht)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert each tag within the bulletin
	for _, tag := range bltn.Tags() {
		_, err = tx.Stmt(db.insertTag).Exec(txid, tag.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the whole transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// InsertEndorsement commits an endorsement into the public record. It DOES NOT
// enforce foreign key constraints. This allows endorsements to come in out of
// order (or in a staggered fashion) endorsing a bulletin that is yet to be
// mined.
func (db *PublicRecord) InsertEndorsement(endo ombproto.Endorsement) error {

}
