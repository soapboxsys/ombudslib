package ombwire

import "testing"

func TestExtractHardErrors(t *testing.T) {
	// Test Failure modes of extraction
	tests := []struct {
		in       []byte // the magic bytes are added at the start of the test
		failWith error
	}{
		{[]byte{0x23, 0x00, 0x00, 0x00, 0x00}, ErrBadWireType},
		{[]byte{0x01, 0xfd, 0xff, 0xff, 0xff, 0xff}, ErrRecordTooBig},
		{[]byte{0x01, 0x04, 0x00, 0x00, 0x00}, ErrRecordTooBig},
	}

	for i, test := range tests {
		b := append(Magic[:], test.in...)
		_, err := extractWireType(b)

		if err != test.failWith {
			t.Fatalf("Expected test(%d) to fail with: %s, instead go %s",
				i, test.failWith, err)
		}
	}
}

func TestExtractSoftErrors(t *testing.T) {
	tests := []struct {
		in []byte // the magic bytes are added at the start of the test
	}{
		{[]byte{0xff, 0x00}},
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{0x02, 0xfd, 0xff, 0xff, 0xff, 0xff}},
		{[]byte{0x03, 0xd3, 0x22, 0x32, 0x22, 0x22}},
		{[]byte{0x01, 0x06, 0x22, 0x32, 0x22, 0x22, 0x31}},
		{[]byte{0x01, 0x04, 0x00, 0x00, 0x00, 0x00}},
	}

	for i, test := range tests {
		b := append(Magic[:], test.in...)
		_, err := extractWireType(b)

		if err == nil {
			t.Fatalf("Expected test(%d) to fail.", i)
		}
	}
}
