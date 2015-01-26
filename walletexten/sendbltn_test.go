package walletexten

import (
	"testing"

	"github.com/NSkelsey/derp"
)

func TestDustAmount(t *testing.T) {
	target := int64(567)
	dustVal := derp.DustAmnt()
	if dustVal != target {
		t.Fatal("Returned a dust val of %d instead of %d", dustVal, target)
	}
}
