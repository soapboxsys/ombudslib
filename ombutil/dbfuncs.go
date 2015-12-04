package ombutil

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombwire"
)

// ProcessBlock accepts a btcutil Block and returns the list of bulletins
// stored within it. Each transaction in the block is examined for our
// very clever duck typing scheme. If it is there, the tx is converted into
// a bulletin and added to the list.
func ProcessBlock(blk *btcutil.Block, net *chaincfg.Params) ([]*Bulletin, error) {
	bltns := []*Bulletin{}

	for _, tx := range blk.Transactions() {
		bltn, ok := ConvertTransaction(tx, blk, net)
		if ok {
			bltns = append(bltns, bltn)
		}
	}

	return bltns, nil
}

// ConvertTransaction parses a bitcoin transaction and pulls a bulletin out of it
// if there is a bulletin encoded within. Otherwise it returns (nil, false), which
// indicates that a valid bulletin was not contained within. Due to our duck typing
// rules all older bulletins that do not conform to the current spec are ignored.
func ConvertTransaction(tx *btcutil.Tx, blk *btcutil.Block, net *chaincfg.Params) (*Bulletin, bool) {
	if !ombwire.IsBulletin(tx.MsgTx()) {
		return nil, false
	}

	bltn, err := NewBulletin(tx.MsgTx(), net)
	if err != nil {
		return nil, false
	}
	bltn.AddBlock(blk)

	return bltn, true
}
