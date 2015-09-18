package ombutil

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/protocol/ombproto"
)

// ProcessBlock accepts a btcutil Block and returns the list of bulletins
// stored within it. Each transaction in the block is examined for our
// very clever duck typing scheme. If it is there, the tx is converted into
// a bulletin and added to the list.
func ProcessBlock(blk *btcutil.Block, net *chaincfg.Params) ([]*ombproto.Bulletin, error) {
	bltns := []*ombproto.Bulletin{}

	sha, _ := blk.Sha()
	for _, tx := range blk.Transactions() {
		bltn, ok := ConvertTransaction(tx, sha, net)
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
func ConvertTransaction(tx *btcutil.Tx, blkSha *wire.ShaHash, net *chaincfg.Params) (*ombproto.Bulletin, bool) {
	if !ombproto.IsBulletin(tx.MsgTx()) {
		return nil, false
	}

	bltn, err := ombproto.NewBulletin(tx.MsgTx(), blkSha, net)
	if err != nil {
		return nil, false
	}

	return bltn, true
}
