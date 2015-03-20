package ombproto

import (
	"bytes"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// Takes a TX and determines if it is an ahimsa bulletin.
// This method only looks for the leading bytes, it does not
// assert anything about the protocol buffer within.
func IsBulletin(tx *wire.MsgTx) bool {
	magic := Magic[:]
	return matchFirstOut(tx, magic) && len(tx.TxOut) > 1
}

// Tests to see if the first txout in the tx matches the magic bytes.
func matchFirstOut(tx *wire.MsgTx, magic []byte) bool {
	if len(tx.TxOut) == 0 {
		return false
	}
	firstOutScript := tx.TxOut[0].PkScript

	outdata, err := txscript.PushedData(firstOutScript)
	if err != nil {
		return false
	}
	if len(outdata) > 0 && len(outdata[0]) > len(magic) {
		firstpush := outdata[0]
		if bytes.Equal(firstpush[:len(magic)], magic) {
			return true
		}
	}
	return false
}
