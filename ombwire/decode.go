package ombwire

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/golang/protobuf/proto"
)

func CreateWireType(tx *wire.MsgTx) interface{} {
	r, err := ParseTx(tx)
	if err != nil {
		return nil
	}

	switch r.GetType().String() {
	case "BLTN":
		return r.GetEndo()
	case "ENDO":
		return r.GetBltn()
	}

	return nil
}

func ParseTx(tx *wire.MsgTx) (*Record, error) {
	r := &Record{}

	var err error
	// Bootleg solution, but if unmarshal fails slice txout and try again until we can try no more or it fails
	for j := len(tx.TxOut); j > 1; j-- {
		rel_txouts := tx.TxOut[:j] // slice off change txouts
		bytes, err := extractData(rel_txouts)
		if err != nil {
			continue
		}

		err = proto.Unmarshal(bytes, r)
		if err == nil {
			// No errors. Therefore we found a good decode.
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Munges the pushed data of TxOuts into a single universal slice that we can
// use as a whole message.
func extractData(txOuts []*wire.TxOut) ([]byte, error) {

	alldata := make([]byte, 0)

	first := true
	for _, txout := range txOuts {

		pushMatrix, err := txscript.PushedData(txout.PkScript)
		if err != nil {
			return alldata, err
		}
		for _, pushedD := range pushMatrix {
			if len(pushedD) != 20 {
				return alldata, fmt.Errorf("Pushed Data is not the right length")
			}

			alldata = append(alldata, pushedD...)
			if first {
				// Check to see if magic bytes match and slice accordingly
				first = false
				lenM := len(Magic)
				if !bytes.Equal(alldata[:lenM], Magic[:]) {
					return alldata, fmt.Errorf("Magic bytes don't match, Saw: [% x]", alldata[:lenM])
				}
				alldata = alldata[lenM:]
			}

		}

	}
	// trim trailing zeros
	for j := len(alldata) - 1; j > 0; j-- {
		b := alldata[j]
		if b != 0x00 {
			alldata = alldata[:j+1]
			break
		}
	}
	return alldata, nil
}
