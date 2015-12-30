package ombwire_test

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/soapboxsys/ombudslib/ombwire"
)

// Test Standard Bulletins
func TestStandardBulletins(t *testing.T) {

	tests := []struct {
		record *ombwire.Bulletin // Expected wire record
		bytes  []byte            // Expected output bytes
	}{
		// SB1 -- Simple Hello Wrald
		{
			ombwire.NewBulletin("Hello world!", 12345678, nil),
			[]byte{
				0x4f, 0x4d, 0x42, 0x55, 0x44, 0x53, 0x01, 0x13, 0x0a, 0x0c, 0x48, 0x65, 0x6c,
				0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x10, 0xce, 0xc2, 0xf1,
				0x05,
			},
		},
		// SB2 -- Hello Wrald with Location
		{
			ombwire.NewBulletin(
				"Hello world!", 12345678,
				ombwire.NewLocation(38.8977, 77.0366, 1000),
			),
			[]byte{
				0x4f, 0x4d, 0x42, 0x55, 0x44, 0x53, 0x01, 0x30, 0x0a, 0x0c, 0x48, 0x65, 0x6c,
				0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x10, 0xce, 0xc2, 0xf1,
				0x05, 0x1a, 0x1b, 0x09, 0x42, 0xcf, 0x66, 0xd5, 0xe7, 0x72, 0x43, 0x40, 0x11,
				0x27, 0xc2, 0x86, 0xa7, 0x57, 0x42, 0x53, 0x40, 0x19, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x40, 0x8f, 0x40,
			},
		},
		// SB3 -- Dec of Independence
		{
			ombwire.NewBulletin(
				"When, in the course of human events, it becomes necessary for one people to",
				12345678, nil,
			),
			[]byte{},
		},
	}

	for i, test := range tests {
		i += 1
		// Check encoding to bytes
		b, err := ombwire.EncodeWireType(test.record)
		if err != nil {
			t.Fatalf("Test SB%d encoding failed with: %s", i, err)
		}

		if !bytes.Equal(b, test.bytes) {
			t.Fatalf("Test SB%d encoding produced incorrect bytes: % x", i, b)
		}

		// Check decoding from bytes
		rbltn, err := ombwire.DecodeWireType(test.bytes)
		bltn := rbltn.(*ombwire.Bulletin)
		if err != nil {
			t.Fatalf("Test SB%d decoding failed with: %s", i, err)
		}

		if !identicalBltn(bltn, test.record) {
			t.Fatal(spew.Sprintf("Test SB%d decoding produced: %v expected: %v",
				i, bltn, test.record))
		}
	}

}

func identicalBltn(a, b *ombwire.Bulletin) bool {
	if a.GetMessage() != b.GetMessage() {
		return false
	} else if a.GetTimestamp() != b.GetTimestamp() {
		return false
	} else if !identicalLoc(a.GetLocation(), b.GetLocation()) {
		return false
	} else {
		return true
	}
}

func identicalLoc(a, b *ombwire.Location) bool {
	if a == nil && b == nil {
		return true
	}
	return a.GetLat() == b.GetLat() && a.GetLon() == b.GetLon() &&
		a.GetH() == b.GetH()
}