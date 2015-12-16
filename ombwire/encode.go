package ombwire

import (
	"bytes"
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/golang/protobuf/proto"
)

func EncodeWireType(m proto.Message) ([]byte, error) {
	empt := []byte{}
	b := make([]byte, 0, MaxRecordLength)

	buf := bytes.NewBuffer(b)

	buf.Write(Magic[:])
	// Write the type byte
	switch m.(type) {
	case *Bulletin:
		buf.WriteByte(BulletinMagic)
	case *Endorsement:
		buf.WriteByte(EndorsementMagic)
	default:
		return empt, errors.New("unsupported type")
	}

	// Write the length of the protocol buf
	s := uint64(proto.Size(m))
	err := writeVarInt(buf, s)
	if err != nil {
		return empt, err
	}
	if 6+4+s > MaxRecordLength {
		return empt, ErrRecordTooBig
	}

	mb, err := proto.Marshal(m)
	if err != nil {
		return empt, err
	}
	// Write the protobuf Message to the buf
	buf.Write(mb)

	return buf.Bytes(), nil
}

// Converts a bulletin into public key scripts for encoding
func (bltn *Bulletin) TxOuts(toBurn int64, net *chaincfg.Params) ([]*wire.TxOut, error) {
	empt := []*wire.TxOut{}

	rawbytes, err := EncodeWireType(bltn)
	if err != nil {
		return empt, err
	}

	numcuts := numOuts(len(rawbytes))

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
			return empt, err
		}
		pkscript, err := txscript.PayToAddrScript(fakeaddr)
		if err != nil {
			return empt, err
		}
		txout := &wire.TxOut{
			PkScript: pkscript,
			Value:    toBurn,
		}

		txouts = append(txouts, txout)
	}
	return txouts, nil
}

// numOuts returns the number of P2PKH outs needed to encode this bulletin
func numOuts(length int) int {
	numouts := length / 20
	if length%20 != 0 {
		numouts += 1
	}
	return numouts
}
