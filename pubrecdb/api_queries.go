package pubrecdb

import (
	"database/sql"
	"strings"

	"github.com/soapboxsys/ombudslib/ombjson"
)

var (
	bltnSql string = `
		SELECT bulletins.txid, bulletins.author, message, bulletins.timestamp, bulletins.block, blocks.timestamp, blocks.height, count(endorsements.txid), latitude, longitude, bulletins.height
	`

	// TODO seems alright
	selectBltnSql string = bltnSql +
		`
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN endorsements ON bulletins.txid = endorsements.bid
		WHERE bulletins.txid = $1 
		GROUP BY bulletins.txid HAVING bulletins.txid NOT null
	`
)

func prepareQueries(db *PublicRecord) error {

	var err error
	db.selectBltn, err = db.conn.Prepare(selectBltnSql)
	if err != nil {
		return err
	}
	return nil
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
