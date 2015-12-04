package ombproto_test

import (
	"testing"

	"github.com/soapboxsys/ombudslib/protocol/ombproto"
)

func TestBulletinCreate(t *testing.T) {
	topic := "King Arthur Arrives in Camelot"
	msg := "What knight live in that castle over there?"

	bltn, err := ombproto.NewBulletinFromStr("nick", topic, msg)
	if err != nil {
		t.Errorf("New failed with: %v", err)
		return
	}

	if bltn.Message != msg {
		t.Errorf("Msgs do not match: %v", err)
		return
	}
}
