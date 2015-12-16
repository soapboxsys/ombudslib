package ombwire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestEncodeWireType(t *testing.T) {

	var m string = "That is a little too high for me to climb"
	var ts uint64 = uint64(123741234)
	var l float64 = float64(322323223)

	orig_bltn := &Bulletin{
		Message:   &m,
		Timestamp: &ts,
		Location: &Location{
			Lat: &l,
			Lon: &l,
			H:   &l,
		},
	}

	b, err := EncodeWireType(orig_bltn)
	if err != nil {
		t.Fatal(err)
	}

	pm, err := DecodeWireType(b)
	if err != nil {
		t.Fatal(err)
	}

	n_bltn := pm.(*Bulletin)
	if n_bltn.GetMessage() != orig_bltn.GetMessage() {
		t.Fatalf("Original and final messages do not match!")
	}
}

func TestEncodeEndo(t *testing.T) {
	bid := []byte{
		0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F,
	}
	var ts uint64 = uint64(123741234)

	i_endo := &Endorsement{
		Bid:       bid,
		Timestamp: &ts,
	}

	b, err := EncodeWireType(i_endo)
	if err != nil {
		t.Fatalf("Endorsement encode failed: %s", err)
	}

	out, err := DecodeWireType(b)
	if err != nil {
		t.Fatalf("Endorsement decode failed: %s", err)
	}
	o_endo := out.(*Endorsement)

	if !bytes.Equal(o_endo.GetBid(), bid) || ts != o_endo.GetTimestamp() {
		t.Fatalf(spew.Sprintf("Orignial and final values differ. Something went"+
			" seriously wrong, the raw endo: %b : %b", o_endo, i_endo))
	}
}

func TestExample_GenerateBltnBytes(t *testing.T) {
	var m string = "Hello, world!"
	var ts uint64 = uint64(12345678)

	bltn := &Bulletin{
		Message:   &m,
		Timestamp: &ts,
	}

	b, _ := EncodeWireType(bltn)
	j, _ := json.Marshal(bltn)

	fmt.Printf("JSON\t\t: %s\n", j)
	fmt.Printf("Header\t\t: % x\n", b[:8])
	fmt.Printf("Wirerecord\t: % x\n", b[8:])
}
