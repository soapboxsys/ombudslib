package pubrecdb

import (
	"database/sql"
	"sort"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
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

	selectBlockSql string = blockHeadSql + `
		FROM blocks LEFT JOIN endorsements ON endorsements.block = blocks.hash
		LEFT JOIN bulletins ON bulletins.block = blocks.hash
		WHERE blocks.hash = $1
		GROUP BY blocks.hash HAVING blocks.hash NOT null
	`

	selectBlockTipSql string = blockHeadSql + `
		FROM blocks LEFT JOIN endorsements ON endorsements.block = blocks.hash
		LEFT JOIN bulletins ON bulletins.block = blocks.hash
		GROUP BY blocks.hash HAVING blocks.hash NOT null
		ORDER BY blocks.height DESC
		LIMIT 1
	`

	selectBestTagsSql string = `
		SELECT tags.value, count(*), bulletins.timestamp 
		FROM tags LEFT JOIN bulletins on tags.txid = bulletins.txid
		GROUP BY tags.value
		ORDER BY count(tags.value) DESC, bulletins.timestamp DESC
		LIMIT $1
	`

	selectAuthorBltnsSql string = bltnSql + `
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN endorsements ON bulletins.txid = endorsements.bid
		WHERE bulletins.author = $1 
		GROUP BY bulletins.txid HAVING bulletins.txid NOT null
		ORDER BY blocks.timestamp DESC
	`

	selectMostEndoBltnsSql string = bltnSql + `
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		INNER JOIN endorsements ON bulletins.txid = endorsements.bid
		GROUP BY bulletins.txid
		ORDER BY count(endorsements.txid) DESC
		LIMIT $1
	`
)

func prepareQueries(db *PublicRecord) error {
	var err error
	db.selectEndosByHeight, err = db.conn.Prepare(selectEndosByHeightSql)

	db.selectMostEndoBltns, err = db.conn.Prepare(selectMostEndoBltnsSql)
	if err != nil {
		return err
	}

	db.selectNearbyBltns, err = db.conn.Prepare(selectNearbyBltns)
	if err != nil {
		return err
	}

	db.selectAuthorEndos, err = db.conn.Prepare(selectAuthorEndosSql)
	if err != nil {
		return err
	}

	db.selectAuthorBltns, err = db.conn.Prepare(selectAuthorBltnsSql)
	if err != nil {
		return err
	}

	db.selectBestTags, err = db.conn.Prepare(selectBestTagsSql)
	if err != nil {
		return err
	}

	db.selectEndosByBid, err = db.conn.Prepare(selectEndosByBidSql)
	if err != nil {
		return err
	}

	db.selectBlock, err = db.conn.Prepare(selectBlockSql)
	if err != nil {
		return err
	}

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

	startH := tipBlk.Head.Height
	stopH := pegBlk.Height() - 1
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

	endos, err := db.GetEndosByHeight(startH, stopH)
	if err != nil {
		return nil, err
	}

	page := &ombjson.Page{
		Start:        tipBlk.Head.Hash,
		Stop:         pegBlk.Sha().String(),
		Bulletins:    bltns,
		Endorsements: endos,
	}

	return page, nil
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

	// TODO Query for endos between heights
	endos, err := db.GetEndosByHeight(startH, stopH)
	if err != nil {
		return nil, err
	}

	page := &ombjson.Page{
		Start:        start.String(),
		Stop:         stop.String(),
		Bulletins:    bltns,
		Endorsements: endos,
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

// GetBestTag returns the 'best' tags as determined by a simple geometric
// formula that takes into account the number of times the tag has been used
// and the first time the tag was seen in the record.
func (db *PublicRecord) GetBestTags() ([]*ombjson.Tag, error) {
	tags := []*ombjson.Tag{}

	rows, err := db.selectBestTags.Query(50)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ts, count int64
		var val string
		err := rows.Scan(&val, &count, &ts)
		if err != nil {
			return tags, err
		}
		tag := ombjson.NewTag(val, count, ts)
		tags = append(tags, &tag)
	}
	// Sort tags by score.
	sort.Sort(ombjson.ByScore(tags))

	return tags, nil
}

// GetBulletin returns a single bulletin as json that is identified by txid.
// If the bltn does not exist the functions returns sql.ErrNoRows. The function
// assumes that the passed txid string is correctly formed (all lower case hex
// string).
func (db *PublicRecord) GetBulletin(txid *wire.ShaHash) (*ombjson.Bulletin, error) {
	row := db.selectBltn.QueryRow(txid.String())
	bltn, err := scanBltn(row)
	if err != nil {
		return nil, err
	}
	// Scan endorsements related to the bulletin
	endos, err := db.GetEndosByBid(txid)
	if err != nil {
		return nil, err
	}
	bltn.Endorsements = endos

	return bltn, nil
}

// getAuthorBltns just returns the bulletins that have been signed by the
// passed bitcoin address.
func (db *PublicRecord) getAuthorBltns(author btcutil.Address) ([]*ombjson.Bulletin, error) {
	rows, err := db.selectAuthorBltns.Query(author.String())
	if err != nil {
		return []*ombjson.Bulletin{}, err
	}
	return scanBltns(rows)
}

// GetAuthor returns the bulletins and the endorsements a bitcoin address has
// sent.
func (db *PublicRecord) GetAuthor(author btcutil.Address) (*ombjson.AuthorResp, error) {

	bltns, err := db.getAuthorBltns(author)
	if err != nil {
		return nil, err
	}

	endos, err := db.getAuthorEndos(author)
	if err != nil {
		return nil, err
	}

	auth := &ombjson.AuthorResp{
		Bulletins:    bltns,
		Endorsements: endos,
	}

	if len(bltns) > 0 {
		auth.Summary = &ombjson.AuthorSummary{
			Address:    author.String(),
			LastBlkTs:  bltns[0].BlockRef.Timestamp,
			FirstBlkTs: bltns[len(bltns)-1].BlockRef.Timestamp,
		}
	}

	return auth, nil
}

// GetMostEndorsedBltns returns a list of bltns sorted by number of
// endorsements received. It does not return bltns with 0 endorsements.
func (db *PublicRecord) GetMostEndorsedBltns(lim int) ([]*ombjson.Bulletin, error) {
	rows, err := db.selectMostEndoBltns.Query(lim)
	defer rows.Close()
	if err != nil {
		return []*ombjson.Bulletin{}, err
	}

	bltns, err := scanBltns(rows)
	if err != nil {
		return []*ombjson.Bulletin{}, err
	}

	return bltns, nil
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

type scannable interface {
	Scan(dest ...interface{}) error
}
