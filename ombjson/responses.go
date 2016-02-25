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
	Txid         string         `json:"txid"`
	Author       string         `json:"author"`
	Message      string         `json:"msg"`
	Timestamp    int64          `json:"timestamp",omitempty`
	NumEndos     int32          `json:"numEndos"`
	BlockRef     *BlockRef      `json:"blkref",omitempty`
	Location     *Location      `json:"loc",omitempty`
	Endorsements []*Endorsement `json:"endos",omitempty`
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
	FirstBlkTs int64  `json:"firstBlkTs,omitempty"`
	LastBlkTs  int64  `json:"lastBlkTs,omitempty"`
}

// Contains info about an author and posts by that author
type AuthorResp struct {
	Summary      *AuthorSummary `json:"summary"`
	Bulletins    []*Bulletin    `json:"bltns",omitempty`
	Endorsements []*Endorsement `json:"endos",omitempty`
}

// Holds meta information about the server
type Status struct {
	BlockTip   *BlockHead `json:"blkTip"`
	Timestamp  int64      `json:"timestamp"`  // The unix time of "now"
	TimestampH string     `json:"timestamph"` // The same as above but human readable
	Uptime     int64      `json:"uptime"`     // The seconds since this object was created
	UptimeH    string     `json:"uptimeh"`    // The same as above but human readable
}

// Holds summary information about a given board
type BoardSummary struct {
	Name       string `json:"name"`
	NumBltns   uint64 `json:"numBltns"`
	CreatedAt  int64  `json:"createdAt"`  // The block timestamp of when this board was started.
	LastActive int64  `json:"lastActive"` // The bulletin timestamp of the latest post.
	CreatedBy  string `json:"createdBy"`
}

// Holds statistics about the public record
type Statistics struct {
	StartTs  int64 `json:"startTs"`
	StopTs   int64 `json:"stopTs"`
	NumBlks  int64 `json:"numBlks"`
	NumBltns int64 `json:"numBltns"`
	NumEndos int64 `json:"numEndos"`
}
