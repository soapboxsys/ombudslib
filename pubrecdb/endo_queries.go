package pubrecdb

import (
	"database/sql"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombjson"
)

var (
	selectEndosByBidSql string = `
		SELECT e.txid, e.author, e.bid, e.timestamp, e.block, 
			   blocks.height, blocks.timestamp, NULL
		From endorsements as e
		LEFT JOIN blocks ON blocks.hash = e.block
		WHERE e.bid = $1
	`

	selectEndoSql string = `
		SELECT e.txid, e.author, e.bid, e.timestamp, e.block, 
			   blocks.height, blocks.timestamp, bulletins.txid
		FROM endorsements as e
		LEFT JOIN blocks ON blocks.hash = e.block
		LEFT JOIN bulletins ON bulletins.txid = e.bid
		WHERE e.txid = $1
	`

	selectAuthorEndosSql string = `
		SELECT e.txid, e.author, e.bid, e.timestamp, e.block, 
			   blocks.height, blocks.timestamp, bulletins.txid
		FROM endorsements as e
		LEFT JOIN blocks ON blocks.hash = e.block
		LEFT JOIN bulletins ON bulletins.txid = e.bid
		WHERE e.author = $1
	`
)

// GetEndorsement returns a single json Endorsement. If the record does not
// exist the method throws sql.ErrNoRows
func (db *PublicRecord) GetEndorsement(txid *wire.ShaHash) (*ombjson.Endorsement, error) {
	row := db.selectEndo.QueryRow(txid.String())
	return scanEndo(row)
}

// GetEndosByBid returns all of the endorsements for a specific bulletin. This
// is used by GetBulletin to fill out the endorsements a specific bulletin has
// received.
func (db *PublicRecord) GetEndosByBid(bid *wire.ShaHash) ([]*ombjson.Endorsement, error) {
	rows, err := db.selectEndosByBid.Query(bid.String())
	defer rows.Close()
	if err != nil {
		return []*ombjson.Endorsement{}, err
	}
	return scanEndos(rows)
}

func (db *PublicRecord) getAuthorEndos(author btcutil.Address) ([]*ombjson.Endorsement, error) {
	rows, err := db.selectAuthorEndos.Query(author.String())
	defer rows.Close()
	if err != nil {
		return []*ombjson.Endorsement{}, err
	}
	return scanEndos(rows)
}

func scanEndos(rows *sql.Rows) ([]*ombjson.Endorsement, error) {
	endos := []*ombjson.Endorsement{}
	for rows.Next() {
		endo, err := scanEndo(rows)
		if err != nil {
			return []*ombjson.Endorsement{}, err
		}
		endos = append(endos, endo)
	}
	return endos, nil
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
