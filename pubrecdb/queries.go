package pubrecdb

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/soapboxsys/ombudslib/ombjson"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

var (
	ErrBltnCensored error = errors.New("Bulletin is withheld for some reason")

	// Used by GetJsonBltn
	selectTxidSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE bulletins.txid = $1
	`
	// Used by GetJsonBlock
	selectBlockHeadSql string = `
		SELECT hash, prevhash, height, blocks.timestamp, count(bulletins.txid) 
		FROM blocks LEFT JOIN bulletins ON blocks.hash = bulletins.block
		WHERE blocks.hash = $1
	`
	selectBlockBltnsSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE blocks.hash = $1
	`

	// Used by GetJsonAuthor
	selectAuthorSql string = `
		SELECT author, count(*), blocks.timestamp
		FROM bulletins LEFT JOIN blocks on bulletins.block = blocks.hash
		WHERE author = $1
		ORDER BY blocks.timestamp ASC
	`

	// Used by GetJsonAuthor
	selectAuthorBltnsSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE author = $1
	`

	// Used by GetJsonBlacklist
	selectBlacklistSql string = `
		SELECT txid, reason from blacklist
	`

	// Used by GetWholeBoard
	selectBoardSumSql string = `
		SELECT board, count(*), last_bltn.bltn_ts, first_bltn.blk_ts, author 
		FROM bulletins, 
			(SELECT max(bulletins.timestamp) AS bltn_ts FROM bulletins WHERE board = $1) AS last_bltn,
			(SELECT min(blocks.timestamp)  blk_ts FROM bulletins JOIN blocks on bulletins.block = blocks.hash
				WHERE board = $1	
			) AS first_bltn
		LEFT JOIN blocks ON bulletins.block = blocks.hash
		WHERE board = $1
		ORDER BY bulletins.timestamp ASC
		LIMIT 1
	`

	// Used by GetWholeBoard
	selectBoardBltnsSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE board = $1
		ORDER BY blocks.timestamp, bulletins.timestamp
	`

	// Used by GetAllBoards
	selectAllBoardsSql string = `
		SELECT board, count(*), max(bulletins.timestamp), blocks.timestamp, author
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		GROUP BY board
		ORDER BY blocks.timestamp ASC
	`

	// Used by GetRecentBltns
	selectRecentConfSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins, (
			SELECT max(blocks.height) AS height FROM blocks	
		) AS tip		
		JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE blocks.height > (tip.height - $1)
		ORDER BY blocks.timestamp DESC
	`

	// Used by GetUnconfirmed
	selectUnconfirmedSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, NULL, NULL, blacklist.reason
		FROM bulletins
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE block IS NULL
		ORDER BY bulletins.timestamp
	`

	// Used by GetBlocksByDay
	selectBlksByDaySql string = `
		SELECT hash, prevhash, height, blocks.timestamp, count(bulletins.txid) 
		FROM blocks LEFT JOIN bulletins ON bulletins.block = blocks.hash
		WHERE blocks.timestamp > $1 AND blocks.timestamp < $2
		GROUP BY blocks.hash
		ORDER BY height
	`

	// Used by LatestBlkAndBltn
	selectDBStatusSql string = `
		SELECT l_blk.timestamp, l_bltn.timestamp
		FROM (SELECT max(blocks.timestamp) AS timestamp FROM blocks) as l_blk,
			 (SELECT max(bulletins.timestamp) AS timestamp FROM bulletins) as l_bltn
	`

	// Used by GetAllAuthors
	selectAllAuthors string = `
	SELECT author, count(*), min(blocks.timestamp)
		FROM bulletins LEFT JOIN blocks on bulletins.block = blocks.hash
		GROUP BY author
	`
)

// Returns information about a single author
func (db *PublicRecord) GetJsonAuthor(address string) (*ombjson.AuthorResp, error) {

	var numBltns uint64
	var addrstr sql.NullString
	var firstBlockTs sql.NullInt64

	row := db.selectAuthor.QueryRow(address)
	err := row.Scan(&addrstr, &numBltns, &firstBlockTs)
	if err != nil {
		return nil, err
	}

	// Check to see if query returned a real row indicating that this author
	// acutally exists.
	if !addrstr.Valid {
		return nil, sql.ErrNoRows
	}

	authorSum := &ombjson.AuthorSummary{
		Address:  address,
		NumBltns: numBltns,
	}

	if firstBlockTs.Valid {
		authorSum.FirstBlkTs = firstBlockTs.Int64
	}

	rows, err := db.selectAuthorBltns.Query(address)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return nil, err
	}

	authorResp := &ombjson.AuthorResp{
		Author:    authorSum,
		Bulletins: bltns,
	}
	return authorResp, nil
}

// Returns the single bulletin in json format, that is identified by txid.
// If the bltn does not exist GetJsonBltn returns sql.ErrNoRows. If the bltn
// is within the blacklist then we will throw an ErrBltnCensored trying to scan
// it in.
func (db *PublicRecord) GetJsonBltn(txid string) (*ombjson.JsonBltn, error) {

	txid = strings.ToLower(txid)

	row := db.selectTxid.QueryRow(txid)
	// If the bulletin is banned withold the bulletin
	withhold := true
	return scanJsonBltn(row, withhold)
}

// Returns a JsonBlkResp which contains a block head and all of the bulletins in
// in the block.
func (db *PublicRecord) GetJsonBlock(hash string) (*ombjson.JsonBlkResp, error) {

	hash = strings.ToLower(hash)

	blkHead, err := db.GetJsonBlockHead(hash)
	if err != nil {
		return nil, err
	}

	rows, err := db.selectBlockBltns.Query(hash)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return nil, err
	}

	blkResp := &ombjson.JsonBlkResp{
		Head:      blkHead,
		Bulletins: bltns,
	}
	return blkResp, nil
}

func (db *PublicRecord) GetJsonBlockHead(hash string) (*ombjson.JsonBlkHead, error) {

	hash = strings.ToLower(hash)

	row := db.selectBlockHead.QueryRow(hash)
	blkHead, err := scanJsonBlk(row)
	if err != nil {
		return nil, err
	}

	return blkHead, nil
}

func (db *PublicRecord) GetJsonBlacklist() ([]*ombjson.BannedBltn, error) {

	blacklist := []*ombjson.BannedBltn{}
	empt := []*ombjson.BannedBltn{}
	rows, err := db.selectBlacklist.Query()
	defer rows.Close()
	if err != nil {
		return empt, err
	}
	for rows.Next() {
		var txid, reason string
		if err := rows.Scan(&txid, &reason); err != nil {
			return empt, err
		}
		bannedBltn := &ombjson.BannedBltn{
			Txid:   txid,
			Reason: reason,
		}
		blacklist = append(blacklist, bannedBltn)

	}

	return blacklist, nil
}

// Returns a board summary and the bulletins posted to that board. This works on
// the null board as well!
func (db *PublicRecord) GetWholeBoard(boardstr string) (*ombjson.WholeBoard, error) {

	// Unescape boardstr and consider the string utf-8. After this unescape we
	// must use unescapedboard because that *IS* the value stored in the DB.
	unescapedboard, err := url.QueryUnescape(boardstr)
	if err != nil {
		return nil, err
	}

	row := db.selectBoardSum.QueryRow(unescapedboard)

	boardSum, err := scanBoardSummary(row)
	if err != nil {
		return nil, err
	}

	rows, err := db.selectBoardBltns.Query(unescapedboard)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return nil, err
	}

	wholeboard := &ombjson.WholeBoard{
		Summary:   boardSum,
		Bulletins: bltns,
	}

	return wholeboard, nil
}

// Returns a board summary for every board in the database.
func (db *PublicRecord) GetAllBoards() ([]*ombjson.BoardSummary, error) {
	boards := []*ombjson.BoardSummary{}

	rows, err := db.selectAllBoards.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		boardSum, err := scanBoardSummary(rows)
		if err != nil {
			return nil, err
		}
		boards = append(boards, boardSum)
	}

	return boards, nil
}

// Returns the last num of confirmed bulletins in the order they were mined
func (db *PublicRecord) GetRecentConf(num int) ([]*ombjson.JsonBltn, error) {

	empt := make([]*ombjson.JsonBltn, 0, num)

	rows, err := db.selectRecentConf.Query(num)
	if err != nil {
		return empt, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return empt, err
	}

	return bltns, nil
}

func (db *PublicRecord) GetUnconfirmed() ([]*ombjson.JsonBltn, error) {
	empt := []*ombjson.JsonBltn{}

	rows, err := db.selectUnconfirmed.Query()
	if err != nil {
		return empt, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return empt, err
	}

	return bltns, nil
}

func (db *PublicRecord) GetBlocksByDay(day time.Time) ([]*ombjson.JsonBlkHead, error) {
	blocks := []*ombjson.JsonBlkHead{}

	start := day.Unix()
	fin := day.AddDate(0, 0, 1).Unix()

	rows, err := db.selectBlksByDay.Query(start, fin)
	defer rows.Close()
	if err != nil {
		return blocks, err
	}

	for rows.Next() {
		blk, err := scanJsonBlk(rows)
		if err != nil {
			return blocks, err
		}

		blocks = append(blocks, blk)
	}

	// Catch case where rows.Next was never true. Caused by the GROUP BY
	if len(blocks) < 1 {
		return blocks, sql.ErrNoRows
	}

	return blocks, nil
}

// Returns the timestamps of the latest block and bulletin by their self
// reported timesetamps. This is entirely gameable by someone who plays
// with their bltn's timestamp, but for now it is a good hueristic to see
// if the db is actively getting written to.
func (db *PublicRecord) LatestBlkAndBltn() (int64, int64, error) {

	var latestBlk, latestBltn int64

	row := db.selectDBStatus.QueryRow()

	err := row.Scan(&latestBlk, &latestBltn)
	if err != nil {
		return -1, -1, err
	}

	return latestBlk, latestBltn, nil
}

func (db *PublicRecord) GetAllAuthors() ([]*ombjson.AuthorSummary, error) {

	authors := []*ombjson.AuthorSummary{}

	rows, err := db.selectAllAuthors.Query()
	if err != nil {
		return authors, err
	}

	for rows.Next() {
		var address string
		var numbltns uint64
		var firstblkts sql.NullInt64

		err := rows.Scan(&address, &numbltns, &firstblkts)
		if err != nil {
			return authors, err
		}

		author := &ombjson.AuthorSummary{
			Address:  address,
			NumBltns: numbltns,
		}

		if firstblkts.Valid {
			author.FirstBlkTs = firstblkts.Int64
		}

		authors = append(authors, author)
	}

	return authors, nil
}
