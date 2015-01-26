package ahimsa_test

import (
	"bytes"
	"testing"

	"github.com/NSkelsey/protocol/ahimsa"

	"code.google.com/p/goprotobuf/proto"
)

func TestBulletinCreate(t *testing.T) {
	topic := "King Arthur Arrives in Camelot"
	msg := "What knight live in that castle over there?"

	bltn, err := ahimsa.NewBulletinFromStr("nick", topic, msg)
	if err != nil {
		t.Errorf("New failed with: %v", err)
		return
	}

	if bltn.Message != msg {
		t.Errorf("Msgs do not match: %v", err)
		return
	}
}

func TestWireCreate(t *testing.T) {

	topic := "King Arthur Arrives in Camelot"
	msg := "What knight live in that castle over there?"

	wireb := &ahimsa.WireBulletin{
		Version: proto.Uint32(ahimsa.ProtocolVersion),
		Topic:   proto.String(topic),
		Message: proto.String(msg),
	}

	pbytes, err := proto.Marshal(wireb)
	if err != nil {
		t.Errorf("Could not marshal WireBulletin: %v", err)
		return
	}

	buf := []byte{
		0x08, 0x01, 0x12, 0x1e, 0x4b, 0x69, 0x6e, 0x67,
		0x20, 0x41, 0x72, 0x74, 0x68, 0x75, 0x72, 0x20,
		0x41, 0x72, 0x72, 0x69, 0x76, 0x65, 0x73, 0x20,
		0x69, 0x6e, 0x20, 0x43, 0x61, 0x6d, 0x65, 0x6c,
		0x6f, 0x74, 0x1a, 0x2b, 0x57, 0x68, 0x61, 0x74,
		0x20, 0x6b, 0x6e, 0x69, 0x67, 0x68, 0x74, 0x20,
		0x6c, 0x69, 0x76, 0x65, 0x20, 0x69, 0x6e, 0x20,
		0x74, 0x68, 0x61, 0x74, 0x20, 0x63, 0x61, 0x73,
		0x74, 0x6c, 0x65, 0x20, 0x6f, 0x76, 0x65, 0x72,
		0x20, 0x74, 0x68, 0x65, 0x72, 0x65, 0x3f,
	}

	if !bytes.Equal(buf, pbytes) {
		t.Errorf("The protocol has changed and this test was not updated.")
	}

	_msg := wireb.GetMessage()
	if msg != _msg {
		t.Errorf("Serialized + deserialized msgs do not match.")
	}

	_topic := wireb.GetTopic()
	if topic != _topic {
		t.Errorf("Serialized + deserialized topics do not match.")
	}

}
