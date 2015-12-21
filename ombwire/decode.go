package ombwire

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/golang/protobuf/proto"
)

func ParseTx(tx *wire.MsgTx) (interface{}, error) {
	b, err := extractData(tx.TxOut[:])
	if err != nil {
		return nil, err
	}
	return DecodeWireType(b)
}

func DecodeWireType(b []byte) (proto.Message, error) {
	buf := bytes.NewBuffer(b)
	if len(b) < 8 {
		return nil, fmt.Errorf("Malformated tx")
	}

	// Find the magic bytes and read past them
	i := bytes.Index(b, Magic[:])
	if i < 0 {
		return nil, fmt.Errorf("No magic prefix")
	}
	buf.Next(i + len(Magic))

	// Extract the record type
	t, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	// Read the length and return the _ bytes used by the var int
	raw_l, _, err := readVarInt(buf)
	if err != nil {
		return nil, fmt.Errorf("Parse failed: %s", err)
	}

	//h_len := len(Magic) + 1 + n // The total length of the header.

	// Assert that the length provided is reasonable
	if raw_l > MaxRecordLength || raw_l > uint64(buf.Len()) {
		return nil, ErrRecordTooBig
	}

	// slice the byte array to the appropriate length
	r := buf.Next(int(raw_l))

	var pm proto.Message
	// Switch on the provided type to unmarshal the record
	switch t {
	case BulletinMagic:
		pm = &Bulletin{}
		err = proto.Unmarshal(r, pm)
		if err != nil {
			return nil, err
		}
	case EndorsementMagic:
		pm = &Endorsement{}
		err = proto.Unmarshal(r, pm)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrBadWireType
	}

	return pm, nil
}

// Munges the pushed data of TxOuts into a single universal slice that we can
// use as a whole message.
func extractData(txOuts []*wire.TxOut) ([]byte, error) {
	alldata := make([]byte, 0)
	empt := []byte{}

	for _, txout := range txOuts {
		pushMatrix, err := txscript.PushedData(txout.PkScript)
		if err != nil {
			return empt, err
		}

		for _, pushedD := range pushMatrix {
			if len(pushedD) != 20 {
				return empt, fmt.Errorf("Pushed Data is not the right length")
			}
			alldata = append(alldata, pushedD...)
		}

	}

	return alldata, nil
}
