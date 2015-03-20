package pubrecdb

import (
	"database/sql"
	"log"

	"github.com/soapboxsys/ombudslib/ombjson"
)

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanJsonBltn(cursor scannable, withhold bool) (*ombjson.JsonBltn, error) {

	var txid, author, msg string
	var board, blockH, bannedReason sql.NullString
	var blkTs, bltnTs sql.NullInt64

	err := cursor.Scan(&txid, &author, &board, &msg, &bltnTs, &blockH, &blkTs, &bannedReason)
	if err != nil {
		return nil, err
	}

	bltn := &ombjson.JsonBltn{
		Txid:    txid,
		Author:  author,
		Message: msg,
	}

	if bltnTs.Valid {
		bltn.Timestamp = bltnTs.Int64
	}

	if board.Valid {
		// // escape board string to be url encoded
		//u, _ := url.Parse("")
		//u.Path = board.String
		//bltn.Board = u.String()
		bltn.Board = board.String
	}

	// If the response contained a block, fill the optional params
	if blockH.Valid {
		bltn.Block = blockH.String
		bltn.BlkTimestamp = blkTs.Int64
	}

	// If the bulletin was banned and withold is flagged then throw ErrBltnCensored
	if bannedReason.Valid {
		if withhold {
			return nil, ErrBltnCensored
		} else {
			// Otherwise, scrub fields and return a censored bltn
			bltn.BannedReason = bannedReason.String
			bltn.Message = ""
		}
	}

	return bltn, nil
}

// Returns a JsonBlk scanned from the cursor
func scanJsonBlk(cursor scannable) (*ombjson.JsonBlkHead, error) {

	var hash, prevhash sql.NullString
	var timestamp, height, numbltns sql.NullInt64

	err := cursor.Scan(&hash, &prevhash, &height, &timestamp, &numbltns)
	if err != nil {
		log.Println("scan failed")
		return nil, err
	}

	if !hash.Valid {
		log.Println("Hello world")
		return nil, sql.ErrNoRows
	}

	blkHead := &ombjson.JsonBlkHead{
		Hash:      hash.String,
		PrevHash:  prevhash.String,
		Height:    uint64(height.Int64),
		Timestamp: timestamp.Int64,
		NumBltns:  uint64(numbltns.Int64),
	}

	return blkHead, nil
}

// Returns the bulletins returned by the sql query.
func getRelevantBltns(rows *sql.Rows) ([]*ombjson.JsonBltn, error) {
	bltns := []*ombjson.JsonBltn{}
	empt := []*ombjson.JsonBltn{}

	for rows.Next() {
		// Include banned bulletins with msg scrubbed.
		bltn, err := scanJsonBltn(rows, false)
		if err != nil {
			return empt, err
		}
		bltns = append(bltns, bltn)
	}

	return bltns, nil
}

func scanBoardSummary(cursor scannable) (*ombjson.BoardSummary, error) {

	var numposts uint64
	var latestact, createdat sql.NullInt64
	var boardstr, createdby sql.NullString

	err := cursor.Scan(&boardstr, &numposts, &latestact, &createdat, &createdby)
	if err != nil {
		return nil, err
	}

	if !createdby.Valid {
		return nil, sql.ErrNoRows
	}

	boardSum := &ombjson.BoardSummary{
		NumBltns:   uint64(numposts),
		CreatedAt:  createdat.Int64,
		LastActive: latestact.Int64,
		CreatedBy:  createdby.String,
	}

	if boardstr.Valid {
		boardSum.Name = boardstr.String
	}

	return boardSum, nil
}
