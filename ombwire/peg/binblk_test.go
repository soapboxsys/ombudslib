package peg_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

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

	if blk.BlockSha() != peg.StartSha {
		t.Fatalf("Returned blk sha differs from peg: blk:[%s] startSha:[%s]\n",
			blk.BlockSha(), peg.StartSha)
	}
}
