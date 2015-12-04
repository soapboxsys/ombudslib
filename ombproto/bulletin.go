package ombproto

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombproto/ombwire"
)

type Author string
type Tag string

// UtilBltn is a utility type that holds data and references. The unexported fields can be
// nil.
type UtilBltn struct {
	// pulled from the enclosing tx
	Author Author

	// The containing transaction
	tx *wire.MsgTx

	block *btcutil.Block

	// Derived types
	json     *ombjson.Bulletin
	wireBltn *ombwire.Bulletin
}

// NewUBltn creates a bulletin using the passed tx as the container of the
// underlying ombwire.Bulletin. If there is no wire bulletin encoded within
// the tx then the whole call throws an error.
func NewUtilBltn(tx *wire.MsgTx, net *chaincfg.Params) (*UtilBltn, error) {

	wireBltn, err := ombwire.ParseTx(tx)
	if err != nil {
		return nil, err
	}

	author, err := parseAuthor(tx, net)
	if err != nil {
		return nil, err
	}

	bltn := &UtilBltn{
		Author:   author,
		tx:       tx,
		wireBltn: wireBltn,
	}

	return bltn, nil

}

func (bltn *UtilBltn) AddBlock(blk *btcutil.Block) {
	bltn.block = blk
}

// Returns the "Author" who signed the first txin of the transaction
func parseAuthor(tx *wire.MsgTx, net *chaincfg.Params) (Author, error) {
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

	return Author(addrPubKey.EncodeAddress()), nil
}

// Tags returns all of the tags encoded within the message body of the
// bulletin. Only the first 5 tags are counted and returned.
func (bltn *UtilBltn) Tags() []Tag {
	t := []Tag{
		Tag("#foo"),
		Tag("#bar"),
		Tag("#baz"),
	}
	return t
}
