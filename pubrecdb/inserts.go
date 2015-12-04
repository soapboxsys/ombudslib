package pubrecdb

import (
	"database/sql"

	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombutil"
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
	db.insertBlockHeadStmt, err = db.conn.Prepare(insertBlockHeadSql)
	if err != nil {
		return err
	}

	db.insertBulletinStmt, err = db.conn.Prepare(insertBulletinSql)
	if err != nil {
		return err
	}

	db.insertTagStmt, err = db.conn.Prepare(insertTagSql)
	if err != nil {
		return err
	}

	return nil
}

// InsertOmbBlock creates a SQL transaction that commits everything in the block
// in one go into the sqlite db. This preserves the consistency of the database
// even in cases where the power fails.
func (db *PublicRecord) InsertOmbBlock(oblk *ombutil.UBlock) error {

	// Start a Sql Transaction
	tx, err := db.conn.Begin()

	err = db.insertBlockHead(tx, oblk.Block)
	if err != nil {
		return tx.Rollback()
	}

	// Insert every bulletin in the block
	for _, bltn := range oblk.Bulletins {
		err = db.insertBulletin(tx, bltn)
		if err != nil {
			return tx.Rollback()
		}
	}

	return tx.Commit()
}

func (db *PublicRecord) insertBlockHead(tx *sql.Tx, blk *btcutil.Block) error {
	h := blk.MsgBlock().Header

	hash := h.BlockSha().String()
	prevhash := h.PrevBlock.String()
	height := blk.Height()
	timestamp := uint32(h.Timestamp.Unix())

	v := h.Version
	m := h.MerkleRoot.String()
	d := h.Bits
	n := h.Nonce

	_, err := tx.Stmt(db.insertBlockHeadStmt).Exec(hash, prevhash, height, timestamp,
		v, m, d, n)
	if err != nil {
		return err
	}

	return nil
}

// InsertBlockHead only inserts the headers of the block into the DB.
// If the block is already in the record then an error is thrown. It
// makes no effort to insert bulletins or endorsements contained within it.
func (db *PublicRecord) InsertBlockHead(blk *btcutil.Block) (bool, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return false, err
	}

	err = db.insertBlockHead(tx, blk)
	if err != nil {
		return false, tx.Rollback()
	}

	if err = tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

// insertBulletin ombutil.Bulletin used here requires references to a MsgTx, a
// Block and a ombwire.Bulletin. It will throw a nil pointer if any of these
// are missing.
func (db *PublicRecord) insertBulletin(tx *sql.Tx, bltn *ombutil.Bulletin) (err error) {

	txid := bltn.Tx.TxSha().String()
	blkHash := bltn.Block.Sha().String()
	ath := string(bltn.Author)

	msg := bltn.Wire.GetMessage()
	ts := bltn.Wire.GetTimestamp()

	loc := bltn.Wire.GetLocation()
	lt := loc.GetLat()
	lg := loc.GetLon()
	ht := loc.GetH()

	// Insert the Bulletin
	_, err = tx.Stmt(db.insertBulletinStmt).Exec(txid, blkHash, ath, msg, ts, lt, lg, ht)
	if err != nil {
		return err
	}

	// Insert each tag within the bulletin
	for _, tag := range bltn.Tags() {
		_, err = tx.Stmt(db.insertTagStmt).Exec(txid, string(tag))
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertBulletin takes a bulletins and inserts it into the pubrecord. A
// ombproto.Bulletin is used here instead of a wire.Bulletin because of the
// high level utilities offered by the ombproto type.
func (db *PublicRecord) InsertBulletin(bltn *ombutil.Bulletin) (err error) {

	// Start a sql transaction to insert the bulletin and the relevant tags.
	var tx *sql.Tx
	if tx, err = db.conn.Begin(); err != nil {
		return err
	}

	err = db.insertBulletin(tx, bltn)
	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

// InsertEndorsement commits an endorsement into the public record. It DOES NOT
// enforce foreign key constraints. This allows endorsements to come in out of
// order (or in a staggered fashion) endorsing a bulletin that is yet to be
// mined.
func (db *PublicRecord) InsertEndorsement(endo ombutil.Endorsement) error {
	return nil
}
