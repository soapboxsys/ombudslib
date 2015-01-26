package ombproto_test

import (
	"testing"

	"github.com/soapboxsys/ombudslib/protocol/ombproto"
)

func TestFullPath(t *testing.T) {
	cmd, err := ombproto.GetCreateSql()

	if err != nil {
		t.Fatal(err)
	}

	if len(cmd) < 20 {
		t.Fatalf("Returned create is broken: [%s]", cmd)
	}
}
