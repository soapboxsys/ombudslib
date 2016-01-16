package ombutil_test

import (
	"testing"

	. "github.com/soapboxsys/ombudslib/ombutil"
)

// Challenge IsBulletin with several blks of transactions that are not
// bulletins.
func TestIsNotBulletin(t *testing.T) {

}

func TestIsBulletin(t *testing.T) {

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
