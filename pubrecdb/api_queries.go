package pubrecdb

import (
	"database/sql"

	"github.com/btcsuite/btcd/wire"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
)

var (
	bltnSql string = `
		SELECT bulletins.txid, bulletins.author, message, bulletins.timestamp, 
		bulletins.block, blocks.timestamp, blocks.height, count(endorsements.txid), 
		latitude, longitude, bulletins.height
	`

	selectBltnSql string = bltnSql + `
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN endorsements ON bulletins.txid = endorsements.bid
		WHERE bulletins.txid = $1 
		GROUP BY bulletins.txid HAVING bulletins.txid NOT null
	`

	selectTagSql string = bltnSql + `
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN endorsements ON bulletins.txid = endorsements.bid
		LEFT JOIN tags ON bulletins.txid = tags.txid
		WHERE tags.value = $1 COLLATE NOCASE
		GROUP BY bulletins.txid HAVING bulletins.txid NOT null
		ORDER BY blocks.height DESC, bulletins.timestamp DESC
		LIMIT $2
	`

	selectBltnsHeightSql string = bltnSql + `
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN endorsements ON bulletins.txid = endorsements.bid
		WHERE blocks.height <= $1 AND blocks.height > $2
		GROUP BY bulletins.txid HAVING bulletins.txid NOT null
	`

	findHeightSql string = `
		SELECT height FROM blocks WHERE hash = $1
	`

	blockHeadSql string = `
		SELECT blocks.hash, prevhash, blocks.height, blocks.timestamp, 
		       count(endorsements.txid), count(bulletins.txid)
	`

	selectBlockTipSql string = blockHeadSql + `
		FROM blocks LEFT JOIN endorsements ON endorsements.block = blocks.hash
		LEFT JOIN bulletins ON bulletins.block = blocks.hash
		GROUP BY blocks.hash HAVING blocks.hash NOT null
		ORDER BY blocks.height DESC
		LIMIT 1
	`

	// TODO implement selectBlock
	selectBlockSql string = blockHeadSql + `
		FROM blocks LEFT JOIN endorsements ON endorsements.block = blocks.hash
		LEFT JOIN bulletins ON bulletins.block = blocks.hash
		WHERE blocks.hash = $1
		GROUP BY blocks.hash HAVING blocks.hash NOT null
	`
)

func prepareQueries(db *PublicRecord) error {

	var err error
	db.selectBlockTip, err = db.conn.Prepare(selectBlockTipSql)
	if err != nil {
		return err
	}

	db.selectBltn, err = db.conn.Prepare(selectBltnSql)
	if err != nil {
		return err
	}

	db.selectTag, err = db.conn.Prepare(selectTagSql)
	if err != nil {
		return err
	}

	db.selectEndo, err = db.conn.Prepare(selectEndoSql)
	if err != nil {
		return err
	}

	db.findHeight, err = db.conn.Prepare(findHeightSql)
	if err != nil {
		return err
	}

	db.selectBltnsHeight, err = db.conn.Prepare(selectBltnsHeightSql)
	if err != nil {
		return err
	}

	return nil
}

// GetLatestPage returns all of the bulletins and endorsements under the max query
// limit.
func (db *PublicRecord) GetLatestPage() (*ombjson.Page, error) {
	tipBlk, err := db.GetBlockTip()
	if err != nil {
		return nil, err
	}

	pegBlk := peg.GetStartBlock()

	// Query for bltns between heights
	rows, err := db.selectBltnsHeight.Query(tipBlk.Head.Height, pegBlk.Height()-1)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := scanBltns(rows)
	if err != nil {
		return nil, err
	}

	page := &ombjson.Page{
		Start:     tipBlk.Head.Hash,
		Stop:      pegBlk.Sha().String(),
		Bulletins: bltns,
	}

	return page, nil
}

func scanBltns(rows *sql.Rows) ([]*ombjson.Bulletin, error) {
	bltns := []*ombjson.Bulletin{}
	for rows.Next() {
		bltn, err := scanBltn(rows)
		if err != nil {
			return []*ombjson.Bulletin{}, err
		}
		bltns = append(bltns, bltn)
	}
	return bltns, nil
}

// QueryRange returns all of the bulletins and endorsements within the
// selected start and stop block
func (db *PublicRecord) QueryRange(start, stop *wire.ShaHash) (*ombjson.Page, error) {
	// Find the heights of the block hashes.
	startH, err := db.FindHeight(start)
	if err != nil {
		return nil, err
	}
	stopH, err := db.FindHeight(stop)
	if err != nil {
		return nil, err
	}

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

	page := &ombjson.Page{
		Start:     start.String(),
		Stop:      stop.String(),
		Bulletins: bltns,
	}

	return page, nil
}

// GetTag returns a blk cursor with all of the bulletins in a tag ordered by the
// bulletins timestamp. If no bulletins exist in the record with that tag, an empty
// list is returned.
func (db *PublicRecord) GetTag(tag ombutil.Tag) (*ombjson.BltnPage, error) {
	rows, err := db.selectTag.Query(string(tag), db.maxQueryLimit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := scanBltns(rows)
	if err != nil {
		return nil, err
	}

	page := &ombjson.BltnPage{
		Bulletins: bltns,
	}

	if len(bltns) < 1 {
		page.Start = peg.GetStartBlock().Sha().String()
		blk, err := db.GetBlockTip()
		if err != nil {
			return nil, err
		}
		page.Start = blk.Head.Hash
	} else {
		page.Start = bltns[0].BlockRef.Hash
		page.Stop = bltns[len(bltns)-1].BlockRef.Hash
	}

	return page, nil
}

// GetBulletin returns a single bulletin as json that is identified by txid.
// If the bltn does not exist the functions returns sql.ErrNoRows. The function
// assumes that the passed txid string is correctly formed (all lower case hex
// string).
func (db *PublicRecord) GetBulletin(txid *wire.ShaHash) (*ombjson.Bulletin, error) {
	row := db.selectBltn.QueryRow(txid.String())
	return scanBltn(row)
}

func scanBltn(cursor scannable) (*ombjson.Bulletin, error) {

	var txid, author, blkHash, msg string
	var bltnTs, blkTs, blkHeight, numEndos int64
	var lat, lon, h sql.NullFloat64

	err := cursor.Scan(&txid, &author, &msg, &bltnTs,
		&blkHash, &blkTs, &blkHeight, &numEndos, &lat, &lon, &h)
	if err != nil {
		return nil, err
	}

	bltn := &ombjson.Bulletin{
		Txid:      txid,
		Author:    author,
		Message:   msg,
		Timestamp: bltnTs,
		BlockRef: &ombjson.BlockRef{
			Hash:      blkHash,
			Timestamp: blkTs,
			Height:    int32(blkHeight),
		},
		NumEndos: int32(numEndos),
	}

	if lat.Valid && lon.Valid && h.Valid {
		bltn.Location = &ombjson.Location{
			Lat: lat.Float64,
			Lon: lon.Float64,
			H:   h.Float64,
		}
	}

	return bltn, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}
