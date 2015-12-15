package ombutil

import (
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombwire"
)

// A helper struct that contains all of the data relevant to ombuds in a
// Bitcoin block.
type UBlock struct {
	Block        *btcutil.Block
	Bulletins    []*Bulletin
	Endorsements []*Endorsement
}

// CreateUBlock parses a btcutil block and parses out the relevant records.
func CreateUBlock(blk *btcutil.Block) *UBlock {
	ublk := &UBlock{
		Block:        blk,
		Bulletins:    []*Bulletin{},
		Endorsements: []*Endorsement{},
	}

	for _, tx := range blk.Transactions() {
		if !ombwire.HasMagic(tx.MsgTx()) {
			continue
		}

		w := ombwire.CreateWireType(tx.MsgTx())
		if w == nil {
			continue
		}

		switch w := w.(type) {
		case *ombwire.Bulletin:
			bltn, err := NewBltn(w, tx, blk)
			if err != nil {
				continue
			}
			ublk.Bulletins = append(ublk.Bulletins, bltn)
		case *ombwire.Endorsement:
			endo, err := NewEndo(w, tx, blk)
			if err != nil {
				continue
			}
			ublk.Endorsements = append(ublk.Endorsements, endo)
		default:
			continue
		}
	}

	return ublk
}
