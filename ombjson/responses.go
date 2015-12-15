package ombjson

// Single Items

type BlkRef struct {
	Hash      string `json:"hash"`
	Height    int64  `json:"h"`
	Timestamp int64  `json:"ts"`
}

// Holds all the information available about a given Bulletin
type Bulletin struct {
	Txid      string   `json:"txid"`
	Author    string   `json:"author"`
	Message   string   `json:"msg"`
	Timestamp int64    `json:"timestamp,omitempty"`
	Tags      []string `json:"tags",omitempty"`
	NumEndos  int32    `json:"numEndos"`
	BlkRef    *BlkRef  `json:"blkref",omitempty"`
}

type Endorsement struct {
	Txid      string  `json:"txid"`      // txid of the endorsements transaction
	Bid       string  `json:"bid"`       // txid of the endorsed bulletin
	Timestamp string  `json:"timestamp"` // User generated timestamp
	BlkRef    *BlkRef `json:"blkref",omitempty`
}

// Holds meta information about a single unique block
type BlockHead struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prevHash"`
	Timestamp int64  `json:"timestamp"`
	Height    uint64 `json:"height"`
	NumBltns  uint64 `json:"numBltns"`
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
