package ombutil

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombjson"
	"github.com/soapboxsys/ombudslib/ombwire"
)

type Author string
type Tag string
type Tags map[Tag]struct{}

// Bulletin is a utility type that holds data and references. The unexported fields can be
// nil.
type Bulletin struct {
	// pulled from the enclosing tx
	Author Author

	// The containing transaction
	Tx *wire.MsgTx

	Block *btcutil.Block

	// Derived types
	Json *ombjson.Bulletin
	Wire *ombwire.Bulletin
}

// NewBltn takes what looks like a bulletin and prepares it for insertion into
// the public record. If there any problems NewBltn throws an error.
func NewBltn(w *ombwire.Bulletin, tx *btcutil.Tx, blk *btcutil.Block, net *chaincfg.Params) (*Bulletin, error) {
	// Validate wire tx msg
	if len(w.GetMessage()) < 1 {
		return nil, fmt.Errorf("Wire msg has no content!")
	}

	// Parse author
	author, err := ParseAuthor(tx.MsgTx(), net)
	if err != nil {
		return nil, err
	}

	// return type
	bltn := &Bulletin{
		Tx:     tx.MsgTx(),
		Block:  blk,
		Wire:   w,
		Author: author,
	}
	return bltn, nil
}

func (bltn *Bulletin) AddBlock(blk *btcutil.Block) {
	bltn.Block = blk
}

// Returns the "Author" who signed the first txin of the transaction
func ParseAuthor(tx *wire.MsgTx, net *chaincfg.Params) (Author, error) {
	if len(tx.TxIn) < 1 {
		return "", errors.New("No TxIns, malformed Bitcoin transaction")
	}
	sigScript := tx.TxIn[0].SignatureScript[:]

	// Setup a script executor to parse the raw bytes of the signature script.
	script, err := txscript.NewEngine(sigScript, tx, 0, txscript.ScriptBip16, nil)
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

// ParseTags returns up to maxNum tags in the passed string. Tags are pulled
// out in iterative order and they are started with a '#' and concluded with a
// tag break character.
func ParseTags(m string) Tags {
	maxNum := 5
	tags := make(Tags)
	var r rune
	var i_s int = 0

	for i := 0; i < len(m); i += i_s {
		r, i_s = utf8.DecodeRuneInString(m[i:])
		if r == '#' {
			var j_s int
			var v rune
			var j int
			for j = i; j < len(m); j += j_s {
				v, j_s = utf8.DecodeRuneInString(m[j:])
				if isTagBreak(v) {
					break
				}
			}
			tag := Tag(m[i:j])
			//fmt.Printf("i: %d, j: %d, tag: %s\n", i, j, tag)
			if len(tags) >= maxNum {
				break
			}
			tags[tag] = struct{}{}
			i = j
			i_s = 0
		}
	}
	return tags
}

// Tags returns all of the tags encoded within the message body of the
// bulletin. Only the first 5 tags are counted and returned.
func (bltn *Bulletin) Tags() Tags {
	m := bltn.Wire.GetMessage()
	return ParseTags(m)
}

var tagBreaks []rune = []rune{
	' ', '\f', '\n', '\r', '\t', '\v', '\u00a0', '\u1680',
	'\u180e', '\u2000', '\u200a', '\u2028', '\u2029',
	'\u202f', '\u205f', '\u3000', '\ufeff',
}

// isTagBreak returns true if the passed rune is one of the acknowledge tag
// seperator characters.
func isTagBreak(r rune) bool {
	for i := range tagBreaks {
		if r == tagBreaks[i] {
			return true
		}
	}
	return false
}
