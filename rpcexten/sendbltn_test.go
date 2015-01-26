package rpcexten_test

import (
	"testing"

	"github.com/soapboxsys/ombudslib/rpcexten"
)

func TestDustAmount(t *testing.T) {
	target := int64(567)
	dustVal := rpcexten.DustAmnt()
	if dustVal != target {
		t.Fatal("Returned a dust val of %d instead of %d", dustVal, target)
	}
}
