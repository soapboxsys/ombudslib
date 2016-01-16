package pubrecdb

import (
	"database/sql"
	"strings"

	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
)

var (
	bltnSql string = `
		SELECT bulletins.txid, bulletins.author, message, bulletins.timestamp, bulletins.block, blocks.timestamp, blocks.height, count(endorsements.txid), latitude, longitude, bulletins.height
	`

	selectBltnSql string = bltnSql +
		`
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
)

func prepareQueries(db *PublicRecord) error {

	var err error
	db.selectBltn, err = db.conn.Prepare(selectBltnSql)
	if err != nil {
		return err
	}

	db.selectTag, err = db.conn.Prepare(selectTagSql)
	if err != nil {
		return err
	}
	return nil
}

// GetTag returns a blk cursor with all of the bulletins in a tag ordered by the
// bulletins timestamp. If no bulletins exist in the record with that tag, an empty
// list is returned.
func (db *PublicRecord) GetTag(tag ombutil.Tag) (*ombjson.Page, error) {
	rows, err := db.selectTag.Query(string(tag), db.maxQueryLimit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns := []*ombjson.Bulletin{}
	for rows.Next() {
		bltn, err := scanBltn(rows)
		if err != nil {
			return nil, err
		}
		bltns = append(bltns, bltn)
	}

	page := &ombjson.Page{
		Bulletins: bltns,
	}

	if len(bltns) < 1 {
		page.Start = peg.GetStartBlock().Sha().String()
		hash, err := db.CurrentTip()
		if err != nil {
			return nil, err
		}
		page.Start = hash
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
func (db *PublicRecord) GetBulletin(txid string) (*ombjson.Bulletin, error) {
	txid = strings.ToLower(txid)

	row := db.selectBltn.QueryRow(txid)
	return scanBltn(row)
}

func scanBltn(cursor scannable) (*ombjson.Bulletin, error) {

	var txid, author, blkHash, msg string
	var bltnTs, blkTs, blkHeight, numEndos int64
	var lat, lon, h sql.NullFloat64

	// TODO properly scan the tx.
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
