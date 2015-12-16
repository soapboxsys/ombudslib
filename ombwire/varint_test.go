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

func TestVarIntWrite(t *testing.T) {
	tests := []struct {
		num int
		out []byte
	}{
		{34, []byte{0x22}},
		{252, []byte{0xfc}},
		{255, []byte{0xfd, 0xff, 0x00}},
		{532423, []byte{0xfe, 0xc7, 0x1f, 0x08, 0x00}},
		{23423422, []byte{0xfe, 0xbe, 0x69, 0x65, 0x01}},
		{4294967295, []byte{0xfe, 0xff, 0xff, 0xff, 0xff}},
	}

	for _, test := range tests {
		buf := bytes.NewBuffer([]byte{})

		err := writeVarInt(buf, uint64(test.num))
		if err != nil {
			t.Fatalf("VarIntWrite for %d failed: ", test.num, err)
		}

		o := buf.Bytes()
		if !bytes.Equal(o, test.out) {
			t.Fatalf("VarIntWrite(%d) bytes were incorrect: % x want: % x",
				test.num, o, test.out)
		}
	}
}
