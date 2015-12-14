package ombutil_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	. "github.com/soapboxsys/ombudslib/ombutil"
)

// Challenge IsBulletin with several blks of transactions that are not
// bulletins.
func TestIsNotBulletin(t *testing.T) {

}

func TestIsBulletin(t *testing.T) {

}

func TestBulletinParse(t *testing.T) {

	// Not enough TxIns
	tx := wire.NewMsgTx()

	_, err := NewBulletin(tx, &chaincfg.MainNetParams)
	if err == nil {
		t.Fatalf("Creating the bltn should have failed: %s", err)
	}

	// Not a Bulletin

	tx = chaincfg.MainNetParams.GenesisBlock.Transactions[0]
	_, err = NewBulletin(tx, &chaincfg.MainNetParams)
	if err == nil {
		t.Fatalf("Creating bltn from coinbase should fail")
	}
}

func TestTagParse(t *testing.T) {
	m1 := "#This #is #a #default #Message"
	t1 := ParseTags(m1)

	ts := []Tag{Tag("#This"), Tag("#is"), Tag("#a"), Tag("#default"), Tag("#Message")}
	if !sameTags(t1, ts) {
		t.Fatalf("Parsed :%b, Wanted :%b", t1, ts)
	}

	m2 := "#AVeryLongTag"
	t2 := ParseTags(m2)

	if len(t2) != 1 {
		ts = []Tag{Tag("#AVeryLongTag")}
		t.Fatalf("Parsed: %b, Wanted: %b", t2, ts)
	}

	m3 := "#More #than #five #tags #are #here also#bad#tags #today"
	t3 := ParseTags(m3)

	ts = []Tag{Tag("#More"), Tag("#than"), Tag("#five"), Tag("#tags"), Tag("#are")}
	if !sameTags(t3, ts) {
		t.Fatalf("Parsed: %b, Wanted: %b", t3, ts)
	}
}

// Utility function to see if the array of tags are exactly equal.
func sameTags(a, b []Tag) bool {
	if len(a) != len(b) {
		return false
	}
	for i, a_i := range a {
		if string(a_i) != string(b[i]) {
			return false
		}
	}
	return true
}
