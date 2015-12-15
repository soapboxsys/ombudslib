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
	return decodeWireType(b)
}

func decodeWireType(b []byte) (proto.Message, error) {
	buf := bytes.NewBuffer(b)
	if len(b) < 8 {
		return nil, fmt.Errorf("Malformated tx")
	}

	// Check the magic byte.
	m := make([]byte, 6)
	buf.Read(m)
	if !bytes.Equal(m, Magic[:]) {
		return nil, fmt.Errorf("TxOut does not start with magic prefix")
	}

	// Extract the record type
	t, _ := buf.ReadByte()

	// Read the length and n bytes read
	raw_l, n, err := readVarInt(buf)
	if err != nil {
		return nil, fmt.Errorf("Parse failed: %s", err)
	}

	h_len := 6 + 1 + n // The total length of the header.

	// Assert that the length provided is reasonable
	if raw_l > MaxRecordLength || raw_l > uint64(len(b)-h_len) {
		return nil, ErrRecordTooBig
	}

	// slice the byte array to the appropriate length
	r := b[h_len : int(raw_l)+h_len]

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
