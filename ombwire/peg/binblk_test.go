package peg_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/soapboxsys/ombudslib/ombwire/peg"
)

func TestDecode(t *testing.T) {
	peg_b, err := peg.Asset("new-year-blk.dat")
	if err != nil {
		t.Fatal(err)
	}

	// Read the whole raw file and do a byte for bytes comparision
	f, err := os.Open("new-year-blk.dat")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	raw_b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(raw_b, peg_b) {
		t.Fatal("Bytes in peg.Asset and ...blk.dat are different!")
	}
}

func TestGetBlk(t *testing.T) {
	blk := peg.GetStartBlock()

	startSha := wire.ShaHash([wire.HashSize]byte{
		0x40, 0xad, 0x1d, 0xfd, 0x78, 0x6a, 0xcf, 0xd5,
		0xb5, 0xdb, 0x00, 0x24, 0x70, 0x14, 0x18, 0x57,
		0x74, 0x90, 0x2f, 0x4b, 0x60, 0x69, 0x6f, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	if !startSha.IsEqual(blk.Sha()) {
		t.Fatalf("Returned blk sha differs from peg: blk:\n[%s] startSha:\n[%s]\n",
			blk.Sha(), startSha)
	}

	startHeight := peg.StartHeight
	if blk.Height() != startHeight {
		t.Fatalf("Block height is wrong. [%d] & [%d]\n",
			blk.Height(), startHeight)
	}
}
