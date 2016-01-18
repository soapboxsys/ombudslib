package ombjson

// Single Items

type BlockRef struct {
	Hash      string `json:"hash"`
	Height    int32  `json:"h"`
	Timestamp int64  `json:"ts"`
}

// Either all of the fields in a location are present or the location is not
// created.
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	H   float64 `json:"h"`
}

// Holds all the information available about a given Bulletin
type Bulletin struct {
	Txid      string    `json:"txid"`
	Author    string    `json:"author"`
	Message   string    `json:"msg"`
	Timestamp int64     `json:"timestamp", omitempty`
	NumEndos  int32     `json:"numEndos"`
	BlockRef  *BlockRef `json:"blkref", omitempty`
	Location  *Location `json:"loc", omitempty`
}

type Endorsement struct {
	Txid       string    `json:"txid"`       // txid of the endorsements transaction
	Author     string    `json:"author"`     // the creator of the endorsement
	Bid        string    `json:"bid"`        // txid of the endorsed bulletin
	Timestamp  int64     `json:"timestamp"`  // User generated timestamp
	BltnExists bool      `json:"bltnExists"` // Indicates existence of Bid in the record
	BlockRef   *BlockRef `json:"blkref",omitempty`
}

// Holds meta information about a single unique block
type BlockHead struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prevHash"`
	Timestamp int64  `json:"timestamp"`
	Height    int32  `json:"height"`
	NumBltns  int32  `json:"numBltns"`
	NumEndos  int32  `json:"numEndos"`
}

// Contains a block head and the contained bulletins in that block
type Block struct {
	Head         *BlockHead     `json:"head"`
	Bulletins    []*Bulletin    `json:"bltns"`
	Endorsements []*Endorsement `json:"endos"`
}

// Holds meta information about a single author
type AuthorSummary struct {
	Address    string `json:"addr"`
	NumBltns   uint64 `json:"numBltns"`
	FirstBlkTs int64  `json:"firstBlkTs,omitempty"`
}

// Contains info about an author and posts by that author
type AuthorResp struct {
	Author       *AuthorSummary `json:"author"`
	Bulletins    []*Bulletin    `json:"bltns",omitempty`
	Endorsements []*Endorsement `json:"endos",omitempty`
}

// Holds meta information about the server
type Status struct {
	Version    string `json:"version"`
	AppStart   int64  `json:"appStart"`
	LatestBlk  int64  `json:"latestBlock"`
	LatestBltn int64  `json:"latestBltn"`
	BlkCount   uint64 `json:"blkCount"`
}

// Holds summary information about a given board
type BoardSummary struct {
	Name       string `json:"name"`
	NumBltns   uint64 `json:"numBltns"`
	CreatedAt  int64  `json:"createdAt"`  // The block timestamp of when this board was started.
	LastActive int64  `json:"lastActive"` // The bulletin timestamp of the latest post.
	CreatedBy  string `json:"createdBy"`
}

// An entire bulletin board that is sorted by default in ascending order
type WholeBoard struct {
	Summary   *BoardSummary `json:"summary"`
	Bulletins []*Bulletin   `json:"bltns"`
}
