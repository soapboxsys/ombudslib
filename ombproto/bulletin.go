package ombproto

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/golang/protobuf/proto"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombproto/ombwire"
	"github.com/soapboxsys/ombudslib/protocol/ombproto/wirebulletin"
)

const (
	MaxBoardLen int = 30
)

var (
	// The delimiter that sits at the front of every bulletin.
	Magic = [8]byte{
		0x42, 0x52, 0x45, 0x54, 0x48, 0x52, 0x45, 0x4e, /* | BRETHREN | */
	}
	ErrMaxBoardLen = errors.New("Board length is too long")
	ErrNoMsg       = errors.New("Message has no content")
)

type author string

// A utility type that holds data and references. The unexported fields can be
// nil.
type Bulletin struct {
	// pulled from the enclosing tx
	Author author
	txid   *wire.ShaHash

	// The containing transaction
	tx *wire.MsgTx

	block *btcutil.Block

	// Derived types
	json     *ombjson.Bulletin
	wireBltn *ombwire.Bulletin
}

// Creates a new bulletin from the containing Tx, supplied author and optional blockhash
// by unpacking txOuts that are considered data. It ignores extra junk behind the protobuffer.
// NewBulletin also asserts aspects of valid bulletins by throwing errors when msg len
// is zero or board len is greater than MaxBoardLen.
func NewBulletin(tx *wire.MsgTx, blkhash *wire.ShaHash, net *chaincfg.Params) (*Bulletin, error) {
	wireBltn := &wirebulletin.WireBulletin{}

	author, err := getAuthor(tx, net)
	if err != nil {
		return nil, err
	}

	// TODO. scrutinize.

	// Bootleg solution, but if unmarshal fails slice txout and try again until we can try no more or it fails
	for j := len(tx.TxOut); j > 1; j-- {
		rel_txouts := tx.TxOut[:j] // slice off change txouts
		bytes, err := extractData(rel_txouts)
		if err != nil {
			continue
		}

		err = proto.Unmarshal(bytes, wireBltn)
		if err == nil {
			// No errors. Therefore we found a good decode.
			break
		}
	}
	if err != nil {
		return nil, err
	}

	board := wireBltn.GetBoard()
	// assert that the length of the board is within its max size!
	if len(board) > MaxBoardLen {
		return nil, ErrMaxBoardLen
	}

	msg := wireBltn.GetMessage()
	// assert that the bulletin has a non zero message length.
	if len(msg) < 1 {
		return nil, ErrNoMsg
	}

	// TODO assert that msg and board are valid UTF-8 strings.
	hash := tx.TxSha()

	bltn := &Bulletin{
		Txid:      &hash,
		Block:     blkhash,
		Author:    author,
		Board:     board,
		Message:   msg,
		Timestamp: time.Unix(wireBltn.GetTimestamp(), 0),
	}

	return bltn, nil
}

// Converts a bulletin into public key scripts for encoding
func (bltn *Bulletin) TxOuts(toBurn int64, net *chaincfg.Params) ([]*wire.TxOut, error) {

	rawbytes, err := bltn.ToBytes()
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

// Returns the "Author" who signed the first txin of the transaction
func getAuthor(tx *wire.MsgTx, net *chaincfg.Params) (string, error) {
	sigScript := tx.TxIn[0].SignatureScript

	dummyTx := wire.NewMsgTx()

	// Setup a script executor to parse the raw bytes of the signature script.
	script, err := txscript.NewEngine(sigScript, dummyTx, 0, txscript.ScriptBip16, nil)
	if err != nil {
		return "", err
	}
	// Step twice due to <sig> <pubkey> format of pay 2pubkeyhash
	script.Step()
	script.Step()
	// Pull off the <pubkey>
	pkBytes := script.GetStack()[1]

	addrPubKey, err := btcutil.NewAddressPubKey(pkBytes, net)
	if err != nil {
		return "", err
	}

	return addrPubKey.EncodeAddress(), nil
}

// Tags returns all of the tags encoded within the message body of the
// bulletin. Only the first 5 tags are counted and returned.
func (bltn *Bulletin) Tags() []Tag {
	t := []Tag{
		NewTag("#foo", bltn),
		NewTag("#bar", bltn),
		NewTag("#baz", bltn),
	}
	return t
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
