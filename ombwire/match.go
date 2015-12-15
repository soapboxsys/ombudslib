package ombwire

import (
	"bytes"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	// The delimiter that sits at the front of every bulletin.
	Magic = [8]byte{
		0x42, 0x52, 0x45, 0x54, 0x48, 0x52, 0x45, 0x4e, /* | BRETHREN | */
	}
)

// HasMagic takes the passed TX and determines if it has the magic bytes
// associated with Ombuds. This method only looks for the leading bytes, it
// does not assert anything about the protocol buffers within.
func HasMagic(tx *wire.MsgTx) bool {
	return len(tx.TxOut) > 1 && matchFirstOut(tx, Magic[:])
}

// Tests to see if the first txout in the tx matches the magic bytes.
func matchFirstOut(tx *wire.MsgTx, magic []byte) bool {
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
