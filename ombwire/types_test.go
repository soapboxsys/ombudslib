package ombwire_test

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/soapboxsys/ombudslib/ombproto/ombwire"
)

func TestWireEndorsement(t *testing.T) {

	bid := []byte{
		0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F,
	}
	var ts uint64 = uint64(123741234)

	endo := &ombwire.Endorsement{
		Bid:       bid,
		Timestamp: &ts,
	}

	b, err := proto.Marshal(endo)
	if err != nil {
		t.Fatalf("Endorsement marshal failed: %s", err)
	}

	rendo := &ombwire.Endorsement{}
	err = proto.Unmarshal(b, rendo)
	if err != nil {
		t.Fatalf("Endorsement unmarshal failed: %s", err)
	}

	if !bytes.Equal(rendo.GetBid(), bid) || ts != rendo.GetTimestamp() {
		t.Fatalf(spew.Sprintf("Orignial and final values differ. Something went"+
			" seriously wrong, the raw endo: %b : %b", rendo, rendo.GetBid()))
	}
}

// TestWireBulletin takes a bulletin through encoding and decoding and tests to
// see if the values we expect are indeed there. To do this a wirebulletin is
// marshaled and unmarshaled via the protobuf library.
func TestWireBulletin(t *testing.T) {

	var m string = "That is a little too high for me to climb"
	var ts uint64 = uint64(123741234)
	var l float64 = float64(322323223)

	bltn := &ombwire.Bulletin{
		Message:   &m,
		Timestamp: &ts,
		Location: &ombwire.Location{
			Lat: &l,
			Lon: &l,
			H:   &l,
		},
	}

	bytes, err := proto.Marshal(bltn)
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}

	rbltn := &ombwire.Bulletin{}

	err = proto.Unmarshal(bytes, rbltn)
	if err != nil {
		t.Fatalf("Unmarshal failed: %s", err)
	}

	// Test to see if unmarshaled fields match
	if m != rbltn.GetMessage() || ts != rbltn.GetTimestamp() ||
		rbltn.GetLocation().GetLat() != l {
		t.Fatalf(spew.Sprintf("Orignial and final values differ. Something went"+
			"seriously wrong, the raw bltn: %s", rbltn))
	}
}

// Utilties
func u(t int) *uint64 {
	i := uint64(t)
	return &i
}

func f(i int) *float64 {
	f := float64(i)
	return &f
}
