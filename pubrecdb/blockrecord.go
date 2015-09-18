package pubrecdb

import (
	"database/sql"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

var (
	insertBlock = `
		INSERT OR REPLACE INTO blocks (hash, prevhash, height, timestamp) VALUES($1, $2, $3, $4)
	`
)

// A BlockRecord maps to a single block stored in the db.
type BlockRecord struct {
	Hash      *wire.ShaHash
	PrevHash  *wire.ShaHash
	Height    uint64
	Timestamp int64
}

// Writes a btcutil.Block to the db. If the block is already in the db, overwrite
// it with the new parameters. Throws an error if there is a problem writing. Does
// not check to see if the block hashes to the proper value.
func (db *PublicRecord) StoreBlock(blk *btcutil.Block) error {
	blkRec := makeBlockRecord(blk)
	return db.StoreBlockRecord(blkRec)
}

// Writes a block to the sqlite db
func (db *PublicRecord) StoreBlockRecord(blkrec *BlockRecord) error {

	cmd := `INSERT INTO blocks (hash, prevhash, height, timestamp) VALUES($1, $2, $3, $4)`

	_, err := db.conn.Exec(cmd,
		blkrec.Hash.String(),
		blkrec.PrevHash.String(),
		blkrec.Height,
		blkrec.Timestamp,
	)

	if err != nil {
		return err
	}
	return nil
}

// Returns a block record specified by target hash. If the block does not exists
// the function returns a sql.ErrNoRows error.
func (db *PublicRecord) GetBlkRecord(target *wire.ShaHash) (*BlockRecord, error) {
	cmd := `SELECT hash, prevhash, height, timestamp FROM blocks WHERE hash = $1`
	row := db.conn.QueryRow(cmd, target.String())

	blkrec, err := scanBlkRec(row)
	if err != nil {
		return nil, err
	}
	return blkrec, nil
}

// Returns the block that has the greatest height according to the db.
func (db *PublicRecord) GetChainTip() (*BlockRecord, error) {
	cmd := `SELECT hash, prevhash, max(height), timestamp FROM blocks`
	row := db.conn.QueryRow(cmd)

	blkrec, err := scanBlkRec(row)
	if err != nil {
		return nil, err
	}
	return blkrec, nil
}

// Creates a Block record from a single row.
func scanBlkRec(row *sql.Row) (*BlockRecord, error) {

	var hash, prevhash string
	var height uint64
	var timestamp int64

	if err := row.Scan(&hash, &prevhash, &height, &timestamp); err != nil {
		return nil, err
	}

	btchash, err := wire.NewShaHashFromStr(hash)
	if err != nil {
		return nil, err
	}

	btcprevhash, err := wire.NewShaHashFromStr(prevhash)
	if err != nil {
		return nil, err
	}

	blkrec := &BlockRecord{
		Hash:      btchash,
		PrevHash:  btcprevhash,
		Height:    height,
		Timestamp: timestamp,
	}

	return blkrec, nil
}

// Returns the current height of the blocks in the db, if db is not initialized
// return 0.
func (db *PublicRecord) GetBlockCount() int64 {
	cmd := `SELECT max(height) FROM blocks`
	row := db.conn.QueryRow(cmd)

	var height uint64
	err := row.Scan(&height)
	if err != nil {
		return 0
	}
	return int64(height)
}
