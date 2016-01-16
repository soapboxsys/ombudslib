package ombutil

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombwire"
)

type Endorsement struct {
	Block  *btcutil.Block
	Tx     *wire.MsgTx
	Author Author

	Wire *ombwire.Endorsement
	Json *ombjson.Endorsement
}

// NewEndo functions very similarly to NewBltn. It bails out if there are any
// problems with the passed wire, tx, or blk.
func NewEndo(w *ombwire.Endorsement, tx *btcutil.Tx, blk *btcutil.Block) (*Endorsement, error) {
	// Check Bid is correct length
	if len(w.GetBid()) != 32 {
		return nil, fmt.Errorf("Endo's bid is wrong len")
	}

	author, err := ParseAuthor(tx.MsgTx(), &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	endo := &Endorsement{
		Block:  blk,
		Tx:     tx.MsgTx(),
		Wire:   w,
		Author: author,
	}

	return endo, nil
}
