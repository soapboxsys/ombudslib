package pubrecdb

import (
	"database/sql"

	"github.com/btcsuite/btcd/wire"
	"github.com/soapboxsys/ombudslib/ombjson"
)

var (
	selectEndoSql string = `
		SELECT e.txid, e.author, e.bid, e.timestamp, e.block, 
			   blocks.height, blocks.timestamp, bulletins.txid
		FROM endorsements as e
		LEFT JOIN blocks ON blocks.hash = e.block
		LEFT JOIN bulletins ON bulletins.txid = e.bid
		WHERE e.txid = $1
	`
)

// GetEndorsement returns a single json Endorsement. If the record does not
// exist the method throws sql.ErrNoRows
func (db *PublicRecord) GetEndorsement(txid *wire.ShaHash) (*ombjson.Endorsement, error) {
	row := db.selectEndo.QueryRow(txid.String())
	return scanEndo(row)
}

func scanEndo(cursor scannable) (*ombjson.Endorsement, error) {

	var txid, blkHash, bid, author string
	var bltnTxid sql.NullString
	var endoTs, blkHeight, blkTs int64

	err := cursor.Scan(&txid, &author, &bid, &endoTs,
		&blkHash, &blkHeight, &blkTs, &bltnTxid)
	if err != nil {
		return nil, err
	}

	endo := &ombjson.Endorsement{
		Txid:       txid,
		Author:     author,
		Bid:        bid,
		BltnExists: false,
		Timestamp:  endoTs,
		BlockRef: &ombjson.BlockRef{
			Hash:      blkHash,
			Timestamp: blkTs,
			Height:    int32(blkHeight),
		},
	}
	if bltnTxid.Valid {
		endo.BltnExists = true
	}
	return endo, nil
}
