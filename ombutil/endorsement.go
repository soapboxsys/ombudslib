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
