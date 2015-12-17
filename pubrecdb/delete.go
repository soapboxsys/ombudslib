package pubrecdb

import (
	"database/sql"
	"errors"

	"github.com/btcsuite/btcd/wire"
)

var (
	ErrNoSuchBlock error = errors.New("no such block in db")
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
func (db *PublicRecord) DeleteBlockTip(sha wire.ShaHash) (bool, error) {

	var tx *sql.Tx
	var err error
	if tx, err = db.conn.Begin(); err != nil {
		return false, err
	}

	// Check to see if the block is the chain tip
	err = db.blockIsTip(tx, &sha)
	if err != nil {
		tx.Rollback()
		return false, err // Would rather pass back error that causes the rollback
	}

	// Delete the block head which cascades deletes through the db
	_, err = tx.Stmt(db.deleteBlockStmt).Exec(sha.String())
	if err != nil {
		return false, tx.Rollback()
	}

	// May have to check if commit had no error.
	return true, tx.Commit()
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
