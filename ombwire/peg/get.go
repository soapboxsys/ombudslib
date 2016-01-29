package peg

import (
	"bytes"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// The first block's hash after the new year is:
// [0000000000000000036f69604b2f9074571814702400dbb5d5cf6a78fd1dad40]
const rawfilename = "new-year-blk.dat"

const StartHeight = int32(391182)

// The testnet's first block hash is from 01/26/2016:
// [00000000002251dfc4286569caff12b36c8aeff081498364910c50c28bda1d02]
const TestStartHeight = int32(658493)

func GetTestStartBlock() *btcutil.Block {
	peg_b, _ := Asset("testnet-peg-blk.dat")
	buf := bytes.NewBuffer(peg_b)
	msg, _, _ := wire.ReadMessage(buf, wire.ProtocolVersion, wire.TestNet3)
	blk := msg.(*wire.MsgBlock)
	ublk := btcutil.NewBlock(blk)
	ublk.SetHeight(TestStartHeight)
	return ublk
}

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
