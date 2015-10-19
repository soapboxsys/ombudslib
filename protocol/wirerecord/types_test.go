package wirerecord

import "testing"

func TestWireEndorsement(t *testing.T) {

	endo := Endorsement{
		Bltnid:    []byte{34, 65, 34, 64, 23},
		Timestamp: 543223,
	}

	t.Logf("%v", endo)
	t.Errorf("Test failed! %v", endo)
}
