package ombutil

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btclog"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombwire"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
)

// A helper struct that contains all of the data relevant to ombuds in a
// Bitcoin block.
type UBlock struct {
	Block        *btcutil.Block
	Bulletins    []*Bulletin
	Endorsements []*Endorsement
}

// CreateUBlock parses a btcutil block and parses out the relevant records. If
// logger is not nil, it is used to report strange problems as the functions
// parses through a block.
func CreateUBlock(blk *btcutil.Block, log btclog.Logger, net *chaincfg.Params) *UBlock {
	ublk := &UBlock{
		Block:        blk,
		Bulletins:    []*Bulletin{},
		Endorsements: []*Endorsement{},
	}

	wLog := func(s string, args ...interface{}) {
		if log != nil {
			log.Warnf(s, args)
		}
	}

	for _, tx := range blk.Transactions() {
		if !ombwire.HasMagic(tx.MsgTx()) {
			continue
		}

		w, err := ombwire.ParseTx(tx.MsgTx())
		if w == nil || err != nil {
			wLog("Parsing wire threw: %s", err)
			continue
		}

		switch w := w.(type) {
		case *ombwire.Bulletin:
			bltn, err := NewBltn(w, tx, blk, net)
			if err != nil {
				wLog("Creating bltn threw: %s", err)
				continue
			}
			ublk.Bulletins = append(ublk.Bulletins, bltn)
		case *ombwire.Endorsement:
			endo, err := NewEndo(w, tx, blk, net)
			if err != nil {
				wLog("Creating endo threw: %s", err)
				continue
			}
			ublk.Endorsements = append(ublk.Endorsements, endo)
		default:
			continue
		}
	}

	return ublk
}

// PastPegDate determines if the passed block was created after the target peg
// date after which entries can be added to the public record.
func PastPegDate(blk *btcutil.Block, net *chaincfg.Params) bool {
	if net.Net == wire.MainNet {
		return blk.Height() > peg.StartHeight
	} else if net.Net == wire.TestNet3 {
		return blk.Height() > peg.TestStartHeight
	} else {
		return false
	}
}
