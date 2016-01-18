package pubrecdb

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
)

// BlockCount returns the number of blocks stored in the DB
func (db *PublicRecord) BlockCount() (int, error) {
	return db.countRows("blocks")
}

func (db *PublicRecord) BulletinCount() (int, error) {
	return db.countRows("bulletins")
}

func (db *PublicRecord) EndoCount() (int, error) {
	return db.countRows("endorsements")
}

func (db *PublicRecord) countRows(table string) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT count(*) FROM %s;`, table)
	err := db.conn.QueryRow(query).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func scanBlockHead(cursor scannable) (*ombjson.Block, error) {
	var hash, prevhash string
	var height, endo_cnt, bltn_cnt int32
	var ts int64

	err := cursor.Scan(&hash, &prevhash, &height, &ts, &endo_cnt, &bltn_cnt)
	if err != nil {
		return nil, err
	}

	blk := &ombjson.Block{
		Head: &ombjson.BlockHead{
			Hash:      hash,
			PrevHash:  prevhash,
			Height:    height,
			Timestamp: ts,
			NumBltns:  bltn_cnt,
			NumEndos:  endo_cnt,
		},
	}

	return blk, nil
}

func (db *PublicRecord) GetBlockTip() (*ombjson.Block, error) {
	row := db.selectBlockTip.QueryRow()
	blk, err := scanBlockHead(row)
	if err != nil {
		return nil, err
	}

	return blk, nil
}

// FindHeight returns the height of the block. It returns -1, sql.ErrNoRows if
// block hash is not in the record. If the passed hash is the prevHash of the
// peg block, FindHeight returns the height from memory, it does not make a
// round trip to the db.
func (db *PublicRecord) FindHeight(hash *wire.ShaHash) (int32, error) {
	firstHash := peg.GetStartBlock().MsgBlock().Header.PrevBlock.Bytes()
	if bytes.Equal(hash.Bytes(), firstHash) {
		return int32(peg.StartHeight - 1), nil
	}

	row := db.findHeight.QueryRow(hash.String())

	var height int32
	err := row.Scan(&height)
	if err != nil {
		return -1, err
	}
	return height, nil
}
