package ombwire

import (
	"bytes"
	"errors"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	MaxRecordLength uint64 = 75000 // Records can be up to 75KB in size.

	// The delimiter that sits at the front of every bulletin.
	Magic = [6]byte{
		0x4f, 0x4d, 0x42, 0x55, 0x44, 0x53, // | OMBUDS |
	}

	// The magic bytes that determine the type of the recorded when it is
	// encoded or decoded.
	BulletinMagic    byte = 0x01
	EndorsementMagic byte = 0x02

	ErrRecordTooBig error = errors.New("record size too big")
	ErrBadWireType  error = errors.New("No such record type")
)

// HasMagic takes the passed TX and determines if it has the magic bytes
// associated with Ombuds. This method only looks for the leading bytes, it
// does not assert anything about the protocol buffers within.
func HasMagic(tx *wire.MsgTx) bool {
	if len(tx.TxOut) == 0 {
		return false
	}
	inFirst := matchFirstOut(tx, Magic[:])

	inLast := matchLastOut(tx, Magic[:])
	return inFirst || inLast
}

// matchFirstOut tests to see if the data within the first txout in the tx
// matches the magic bytes.
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

// matchLastOut determines if the magic prefix is contained within the last
// txOut of the passed transaction.
func matchLastOut(tx *wire.MsgTx, magic []byte) bool {

	lastOutScript := tx.TxOut[len(tx.TxOut)-1].PkScript

	outdata, err := txscript.PushedData(lastOutScript)
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
