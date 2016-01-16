package pubrecdb

import (
	"database/sql"

	"github.com/btcsuite/btcd/wire"
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
		INSERT INTO endorsements (txid, block, bid, timestamp, author) 
		VALUES ($1, $2, $3, $4, $5)
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

	db.insertEndorsementStmt, err = db.conn.Prepare(insertEndoSql)
	if err != nil {
		return err
	}

	return nil
}

// InsertUBlock creates a SQL transaction that commits everything in the block
// in one go into the sqlite db. This preserves the consistency of the database
// even in cases where the power fails. If the insert was succesful the
// funciton will return (nil, true). If (anything, false) then the insert
// failed.
func (db *PublicRecord) InsertUBlock(oblk *ombutil.UBlock) (error, bool) {

	// Start a Sql Transaction
	tx, err := db.conn.Begin()

	err = db.insertBlockHead(tx, oblk.Block)
	if err != nil {
		return tx.Rollback(), false
	}

	// Insert every bulletin in the block
	for _, bltn := range oblk.Bulletins {
		err = db.insertBulletin(tx, bltn)
		if err != nil {
			return tx.Rollback(), false
		}
	}

	// Insert every endorsement
	for _, endo := range oblk.Endorsements {
		err = db.insertEndorsement(tx, endo)
		if err != nil {
			return tx.Rollback(), false
		}
	}

	return tx.Commit(), true
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
func (db *PublicRecord) InsertBlockHead(blk *btcutil.Block) (error, bool) {
	tx, err := db.conn.Begin()
	if err != nil {
		return err, false
	}

	err = db.insertBlockHead(tx, blk)
	if err != nil {
		return tx.Rollback(), false
	}

	if err = tx.Commit(); err != nil {
		return err, true
	}

	return nil, true
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

	// Execute the insert sql statement
	_, err = tx.Stmt(db.insertBulletinStmt).Exec(txid, blkHash, ath, msg,
		ts, lt, lg, ht)
	if err != nil {
		return err
	}

	// Insert each tag within the bulletin
	for tag, _ := range bltn.Tags() {
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
func (db *PublicRecord) InsertBulletin(bltn *ombutil.Bulletin) (error, bool) {

	// Start a sql transaction to insert the bulletin and the relevant tags.
	var tx *sql.Tx
	var err error
	if tx, err = db.conn.Begin(); err != nil {
		return err, false
	}

	err = db.insertBulletin(tx, bltn)
	if err != nil {
		return tx.Rollback(), false
	}

	if err = tx.Commit(); err != nil {
		return err, true
	}

	return nil, true
}

func (db *PublicRecord) insertEndorsement(tx *sql.Tx, endo *ombutil.Endorsement) error {

	txid := endo.Tx.TxSha().String()
	blkHash := endo.Block.Sha().String()

	h, err := wire.NewShaHash(endo.Wire.GetBid())
	if err != nil {
		return err
	}

	bid := h.String()
	auth := string(endo.Author)
	time := endo.Wire.GetTimestamp()

	_, err = tx.Stmt(db.insertEndorsementStmt).Exec(txid, blkHash, bid, auth, time)
	if err != nil {
		return err
	}

	return nil
}

// InsertEndorsement commits an endorsement into the public record. It DOES NOT
// enforce foreign key constraints. This allows endorsements to come in out of
// order (or in a staggered fashion) endorsing a bulletin that is yet to be
// mined.
func (db *PublicRecord) InsertEndorsement(endo *ombutil.Endorsement) (error, bool) {
	var tx *sql.Tx
	var err error
	if tx, err = db.conn.Begin(); err != nil {
		return err, false
	}

	err = db.insertEndorsement(tx, endo)
	if err != nil {
		return tx.Rollback(), false
	}

	if err = tx.Commit(); err != nil {
		return err, false
	}
	return nil, true
}
