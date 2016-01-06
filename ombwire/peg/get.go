package peg

import (
	"bytes"

	"github.com/btcsuite/btcd/wire"
)

var rawfilename = "new-year-blk.dat"

var StartSha = wire.ShaHash([wire.HashSize]byte{
	0x40, 0xad, 0x1d, 0xfd, 0x78, 0x6a, 0xcf, 0xd5,
	0xb5, 0xdb, 0x00, 0x24, 0x70, 0x14, 0x18, 0x57,
	0x74, 0x90, 0x2f, 0x4b, 0x60, 0x69, 0x6f, 0x03,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
})

// GetStartBlock returns the pegged block which is assembled from the
// underlying binary data.
func GetStartBlock() *wire.MsgBlock {
	peg_b, _ := Asset(rawfilename)

	buf := bytes.NewBuffer(peg_b)

	msg, _, _ := wire.ReadMessage(buf, wire.ProtocolVersion, wire.MainNet)

	blk := msg.(*wire.MsgBlock)
	return blk
}
