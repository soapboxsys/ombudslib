package ombwire

// readVarInt reads a variable length integer from r and returns it as a
import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// readVarInt returns the varint(uint64), the number of bytes read and an
// error if there are any.
func readVarInt(r io.Reader) (uint64, int, error) {
	var b [8]byte
	_, err := io.ReadFull(r, b[0:1])
	if err != nil {
		return 0, 1, err
	}

	var rv uint64
	var n int // number of bytes
	discriminant := uint8(b[0])
	switch discriminant {
	case 0xff:
		n = 9
		_, err := io.ReadFull(r, b[:])
		if err != nil {
			return 0, n, err
		}
		rv = binary.LittleEndian.Uint64(b[:])

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes
		min := uint64(0x100000000)
		if rv < min {
			return 0, n, fmt.Errorf("readVarInt: %d, %d, %d", rv, discriminant, min)
		}

	case 0xfe:
		n = 5
		_, err := io.ReadFull(r, b[0:4])
		if err != nil {
			return 0, n, err
		}
		rv = uint64(binary.LittleEndian.Uint32(b[:]))

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes
		min := uint64(0x10000)
		if rv < min {
			return 0, n, fmt.Errorf("readVarInt: %d, %d, %d", rv, discriminant, min)
		}

	case 0xfd:
		n = 3
		_, err := io.ReadFull(r, b[0:2])
		if err != nil {
			return 0, n, err
		}
		rv = uint64(binary.LittleEndian.Uint16(b[:]))

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes
		min := uint64(0xfd)
		if rv < min {
			return 0, n, fmt.Errorf("readVarInt: %d, %d, %d", rv, discriminant, min)
		}

	default:
		rv = uint64(discriminant)
		n = 1
	}

	return rv, n, nil
}

// writeVarInt serializes val to w using a variable number of bytes depending
// on its value.
func writeVarInt(w io.Writer, val uint64) error {
	if val < 0xfd {
		_, err := w.Write([]byte{uint8(val)})
		return err
	}

	if val <= math.MaxUint16 {
		var buf [3]byte
		buf[0] = 0xfd
		binary.LittleEndian.PutUint16(buf[1:], uint16(val))
		_, err := w.Write(buf[:])
		return err
	}

	if val <= math.MaxUint32 {
		var buf [5]byte
		buf[0] = 0xfe
		binary.LittleEndian.PutUint32(buf[1:], uint32(val))
		_, err := w.Write(buf[:])
		return err
	}

	var buf [9]byte
	buf[0] = 0xff
	binary.LittleEndian.PutUint64(buf[1:], val)
	_, err := w.Write(buf[:])
	return err
}
