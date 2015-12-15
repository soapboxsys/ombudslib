package ombwire

// readVarInt reads a variable length integer from r and returns it as a
import (
	"encoding/binary"
	"fmt"
	"io"
)

// uint64.
func readVarInt(r io.Reader) (uint64, error) {
	var b [8]byte
	_, err := io.ReadFull(r, b[0:1])
	if err != nil {
		return 0, err
	}

	var rv uint64
	discriminant := uint8(b[0])
	switch discriminant {
	case 0xff:
		_, err := io.ReadFull(r, b[:])
		if err != nil {
			return 0, err
		}
		rv = binary.LittleEndian.Uint64(b[:])

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes
		min := uint64(0x100000000)
		if rv < min {
			return 0, fmt.Errorf("readVarInt: %d, %d, %d", rv, discriminant, min)
		}

	case 0xfe:
		_, err := io.ReadFull(r, b[0:4])
		if err != nil {
			return 0, err
		}
		rv = uint64(binary.LittleEndian.Uint32(b[:]))

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes
		min := uint64(0x10000)
		if rv < min {
			return 0, fmt.Errorf("readVarInt: %d, %d, %d", rv, discriminant, min)
		}

	case 0xfd:
		_, err := io.ReadFull(r, b[0:2])
		if err != nil {
			return 0, err
		}
		rv = uint64(binary.LittleEndian.Uint16(b[:]))

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes
		min := uint64(0xfd)
		if rv < min {

			return 0, fmt.Errorf("readVarInt: %d, %d, %d", rv, discriminant, min)
		}

	default:
		rv = uint64(discriminant)
	}

	return rv, nil
}
