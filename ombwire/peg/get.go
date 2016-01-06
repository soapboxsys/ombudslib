package peg

import (
	"bytes"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const rawfilename = "new-year-blk.dat"
const StartHeight = int32(391182)

// GetStartBlock returns the pegged block which is assembled from the
// underlying binary data.
func GetStartBlock() *btcutil.Block {
	peg_b, _ := Asset(rawfilename)

	buf := bytes.NewBuffer(peg_b)

	msg, _, _ := wire.ReadMessage(buf, wire.ProtocolVersion, wire.MainNet)

	blk := msg.(*wire.MsgBlock)
	ublk := btcutil.NewBlock(blk)
	ublk.SetHeight(StartHeight)
	return ublk
}
