package ombutil

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombwire"
)

type Endorsement struct {
	Block *btcutil.Block
	Tx    *wire.MsgTx

	Wire *ombwire.Endorsement
	Json *ombjson.Endorsement
}

func NewEndo(w *ombwire.Endorsement, tx *btcutil.Tx, blk *btcutil.Block) (*Endorsement, error) {
	// validate w
	endo := &Endorsement{
		Block: blk,
		Tx:    tx.MsgTx(),
		Wire:  w, // Check bid to see if it is a real hash
	}

	return endo, nil
}
