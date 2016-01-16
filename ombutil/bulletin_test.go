package ombutil_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	. "github.com/soapboxsys/ombudslib/ombutil"
)

// Challenge IsBulletin with several blks of transactions that are not
// bulletins.
func TestRejectBulletin(t *testing.T) {

}

func TestDetectBulletin(t *testing.T) {

}

// Ensure that parse author can extract the records author out of a raw bitcoin
// tx.
func TestParseAuthor(t *testing.T) {
	var hexTx string = `0100000001d275627c84029b6c46155bb55423d5610936a1680bd2c900bda78fa0730f5178000000006a4730440220537a3bba833876116d55dedf09b696552a83cd0acd0fd8b7d30e5c67e339c66d02202eeb72ba5fe251e243c121bc377471de29335a386b1b672dd76dfb980b947d24012103ee01b63fde1a69fd75d8714ce02010fbc1d025a1c1e72ecee78805c8316092f5ffffffff022202000000000000296a274f4d42554453011f0a1754686520427269746973682061726520636f6d696e67211092b097b405354f1000000000001976a9148eace5df09b9a54f03661dd2cc6646f1b353329088ac00000000`

	b, _ := hex.DecodeString(hexTx)

	tx := wire.NewMsgTx()
	err := tx.BtcDecode(bytes.NewBuffer(b), wire.ProtocolVersion)
	if err != nil {
		t.Fatal(err)
	}

	a, err := ParseAuthor(tx, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("Parsing valid author failed with: %s", err)
	}

	want := "18QxSqzxsWL9Fc6jpBM9vXZy6wZ7XFLCFy"
	if string(a) != want {
		t.Fatalf("Parsed: %s Wanted: %s", string(a), want)
	}
}

func TestTagParse(t *testing.T) {
	m1 := "#This #is #a #default #Message"
	t1 := ParseTags(m1)
	e := struct{}{}

	ts := Tags{
		Tag("#This"):    e,
		Tag("#is"):      e,
		Tag("#a"):       e,
		Tag("#default"): e,
		Tag("#Message"): e,
	}
	if !sameTags(t1, ts) {
		t.Fatalf("Parsed :%b, Wanted :%b", t1, ts)
	}

	m2 := "#AVeryLongTag"
	t2 := ParseTags(m2)

	if len(t2) != 1 {
		ts = Tags{Tag("#AVeryLongTag"): e}
		t.Fatalf("Parsed: %b, Wanted: %b", t2, ts)
	}

	m3 := "#More #than #five #tags #are #here also#bad#tags #today"
	t3 := ParseTags(m3)

	ts = Tags{
		Tag("#More"): e,
		Tag("#than"): e,
		Tag("#five"): e,
		Tag("#tags"): e,
		Tag("#are"):  e,
	}
	if !sameTags(t3, ts) {
		t.Fatalf("Parsed: %b, Wanted: %b", t3, ts)
	}

	// Duplicate tags are ignored
	ts = Tags{
		Tag("#more"): e,
		Tag("#than"): e,
	}
	m4 := "#more #more #more #more #more #than #than #than"
	t4 := ParseTags(m4)
	if !sameTags(t4, ts) {
		t.Fatalf("Parsed: %b, Wanted: %b", t4, ts)
	}
}

// Utility function to see if the array of tags are exactly equal.
func sameTags(a, b Tags) bool {
	if len(a) != len(b) {
		return false
	}
	for a_i := range a {
		if _, ok := b[a_i]; !ok {
			return false
		}
	}
	return true
}
