package ombwire

import (
	"bytes"
	"testing"
)

func TestReadVarInt(t *testing.T) {
	tests := []struct {
		in  []byte
		out uint64
	}{
		{[]byte{0x32}, uint64(50)},
		{[]byte{0xfd, 0x32, 0x01}, uint64(306)},
		{[]byte{0xfd, 0xff, 0xff}, uint64(65535)},
		{[]byte{0xfe, 0xff, 0xff, 0xff, 0xff}, uint64(4294967295)},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff}, uint64(18446744073709551615)},
	}

	for _, tst := range tests {
		var buf *bytes.Buffer = bytes.NewBuffer(tst.in)

		c, _, err := readVarInt(buf)
		if err != nil {
			t.Fatalf("Read of %x failed with err: %s", tst.in, err)
		}
		if c != tst.out {
			t.Fatalf("VarInt should equal: %d not %d", tst.out, c)
		}
	}

}
