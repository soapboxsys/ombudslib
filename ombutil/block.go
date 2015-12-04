package ombutil

import "github.com/btcsuite/btcutil"

// A helper datastruct that contains all of the data relevant to ombuds in a
// Bitcoin block.
type UBlock struct {
	Block        *btcutil.Block
	Bulletins    []*Bulletin
	Endorsements []*Endorsement
}

// NewUBlock parses a btcutil block and parses out the relevant records.
func NewUBlock(blk *btcutil.Block) *UBlock {
	ublk := &UBlock{
		Block: blk,
	}
	return ublk
}
