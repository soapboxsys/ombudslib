package pubrecdb

import (
	"bytes"
	"fmt"
	"time"

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

func (db *PublicRecord) getBltnsByHeight(startH, stopH int32) ([]*ombjson.Bulletin, error) {
	// Query for bltns between heights
	rows, err := db.selectBltnsHeight.Query(startH, stopH)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := scanBltns(rows)
	if err != nil {
		return nil, err
	}
	return bltns, nil
}

func (db *PublicRecord) getBlockBltns(height int32) ([]*ombjson.Bulletin, error) {
	return db.getBltnsByHeight(height, height-1)
}

func (db *PublicRecord) getBlockEndos(height int32) ([]*ombjson.Endorsement, error) {
	return db.GetEndosByHeight(height, height-1)
}

// GetBlock returns the block in the record specified by 'hash'. If it is not
// present then sql.ErrNoRows is returned.
func (db *PublicRecord) GetBlock(hash *wire.ShaHash) (*ombjson.Block, error) {
	row := db.selectBlock.QueryRow(hash.String())
	block, err := scanBlockHead(row)
	if err != nil {
		return &ombjson.Block{}, err
	}

	bltns, err := db.getBlockBltns(block.Head.Height)
	if err != nil {
		return &ombjson.Block{}, err
	}

	endos, err := db.getBlockEndos(block.Head.Height)
	if err != nil {
		return &ombjson.Block{}, err
	}

	block.Bulletins = bltns
	block.Endorsements = endos
	return block, nil
}

// GetBlockTip works exactly the same as GetBlock except that the query always
// returns the block at the tip of the chain. This is the block that has the
// greatest height in the record.
func (db *PublicRecord) GetBlockTip() (*ombjson.Block, error) {
	row := db.selectBlockTip.QueryRow()
	return scanBlockHead(row)
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

var computeStatisticsSql = `
SELECT 
	( 
		SELECT COUNT(*) FROM bulletins 
		WHERE bulletins.timestamp > $1 AND bulletins.timestamp <= $2
	) as bltn_cnt,
	( 
		SELECT COUNT(*) FROM endorsements 
		WHERE endorsements.timestamp > $1 AND endorsements.timestamp <= $2
	) as endo_cnt,
	( 
		SELECT COUNT(*) FROM blocks 
		WHERE blocks.timestamp > $1 AND blocks.timestamp <= $2
	) as blk_cnt
`

// GetStatistics returns information about some continous period of time over
// which records and blocks were stored in the public record. Currently it only
// produces counts, but one day it will produce interesting facts for the
// scientists to measure.
func (db *PublicRecord) GetStatistics(start, fin time.Time) (*ombjson.Statistics, error) {

	empt := &ombjson.Statistics{}
	row := db.computeStatistics.QueryRow(start.Unix(), fin.Unix())

	var nBltns, nEndos, nBlks int64
	err := row.Scan(&nBltns, &nEndos, &nBlks)
	if err != nil {
		return empt, err
	}
	stat := &ombjson.Statistics{
		StartTs:  start.Unix(),
		StopTs:   fin.Unix(),
		NumBltns: nBltns,
		NumEndos: nEndos,
		NumBlks:  nBlks,
	}
	return stat, nil
}
