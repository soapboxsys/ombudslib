package pubrecdb

import (
	"database/sql"
	"errors"

	"github.com/btcsuite/btcd/wire"
)

var (
	ErrBlockNotTip error = errors.New("block is not chain tip")

	// This stmt causes a foreign key cascade.
	deleteBlockSql string = `
	DELETE FROM blocks WHERE hash == $1;
	`

	// Utility query
	blockIsTipSql string = `
	SELECT EXISTS(SELECT hash FROM 
		(SELECT hash, MAX(height) FROM blocks) 
	WHERE hash == $1) AND 
	NOT EXISTS(SELECT hash FROM blocks where prevhash == $1);
	`
)

func prepareDeletes(db *PublicRecord) (err error) {
	db.deleteBlockStmt, err = db.conn.Prepare(deleteBlockSql)
	if err != nil {
		return err
	}
	db.blockIsTipStmt, err = db.conn.Prepare(blockIsTipSql)
	if err != nil {
		return err
	}
	return
}

// DeleteBlockTip deletes the block header and any dependent data. It will throw
// a ErrBlockNotTip  error if you try to delete any block that is not the tip.
func (db *PublicRecord) DeleteBlockTip(sha *wire.ShaHash) (error, bool) {

	var tx *sql.Tx
	var err error
	if tx, err = db.conn.Begin(); err != nil {
		return err, false
	}

	// Check to see if the block is the chain tip
	err = db.blockIsTip(tx, sha)
	if err != nil {
		tx.Rollback()
		return err, false // Would rather pass back error that causes the rollback
	}

	// Delete the block head which cascades deletes through the db
	_, err = tx.Stmt(db.deleteBlockStmt).Exec(sha.String())
	if err != nil {
		return tx.Rollback(), false
	}

	// Have to check that commit did not fail
	if err = tx.Commit(); err != nil {
		return err, true
	}

	return nil, true
}

// blockIsTip determines if the block is in the db and if it is the current
// chain tip. Returning a nil indicates that the block is the Chain Tip.
func (db *PublicRecord) blockIsTip(tx *sql.Tx, sha *wire.ShaHash) error {
	r, err := tx.Stmt(db.blockIsTipStmt).Query(sha.String())
	defer r.Close()
	if err != nil {
		return err
	}

	var isTip bool
	r.Next()
	if err = r.Scan(&isTip); err != nil {
		return err
	}

	if !isTip {
		return ErrBlockNotTip
	}
	return nil
}

// DropAfterBlockBySha deletes all of the blocks after the passed sha in the
// database.
func (db *PublicRecord) DropAfterBlockBySha(sha *wire.ShaHash) error {
	blk, err := db.GetBlock(sha)
	if err != nil {
		return err
	}

	query := "DELETE FROM blocks WHERE height > $1"
	_, err = db.conn.Exec(query, blk.Head.Height)
	return err
}
