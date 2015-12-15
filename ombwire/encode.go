package ombwire

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/golang/protobuf/proto"
)

// Converts a bulletin into public key scripts for encoding
func (bltn *Bulletin) TxOuts(toBurn int64, net *chaincfg.Params) ([]*wire.TxOut, error) {

	rawbytes, err := proto.Marshal(bltn)
	if err != nil {
		return []*wire.TxOut{}, err
	}

	numcuts, _ := bltn.NumOuts()

	cuts := make([][]byte, numcuts, numcuts)
	for i := 0; i < numcuts; i++ {
		sliceb := make([]byte, 20, 20)
		copy(sliceb, rawbytes)
		cuts[i] = sliceb
		if len(rawbytes) >= 20 {
			rawbytes = rawbytes[20:]
		}
	}

	// Convert raw data into txouts
	txouts := make([]*wire.TxOut, 0)
	for _, cut := range cuts {

		fakeaddr, err := btcutil.NewAddressPubKeyHash(cut, net)
		if err != nil {
			return []*wire.TxOut{}, err
		}
		pkscript, err := txscript.PayToAddrScript(fakeaddr)
		if err != nil {
			return []*wire.TxOut{}, err
		}
		txout := &wire.TxOut{
			PkScript: pkscript,
			Value:    toBurn,
		}

		txouts = append(txouts, txout)
	}
	return txouts, nil
}

// Takes a bulletin and converts into a byte array. A bulletin has two
// components. The leading 8 magic bytes and then the serialized protocol
// buffer that contains the real message 'payload'.
func (bltn *Bulletin) Bytes() ([]byte, error) {
	payload := make([]byte, 0)

	pbytes, err := proto.Marshal(bltn)
	if err != nil {
		return payload, err
	}

	payload = append(payload, Magic[:]...)
	payload = append(payload, pbytes...)
	return payload, nil
}

// Returns the number of txouts needed to encode this bulletin
func (bltn *Bulletin) NumOuts() (int, error) {

	rawbytes, err := bltn.Bytes()
	if err != nil {
		return 0, err
	}

	numouts := len(rawbytes) / 20
	if len(rawbytes)%20 != 0 {
		numouts += 1
	}

	return numouts, nil
}
