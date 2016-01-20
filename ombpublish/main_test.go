package ombpublish_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func TestPublishBulletin(t *testing.T) {
	// NOTE THESE CREDS ARE LEAKED INTO GITHUB
	var rpcuser string = "rpcuser"
	var rpcpass string = "14PkTcQcKj8AsTpljxsExs3Idb2CQGRfHcNcUKjgUTKdM0994ce0eN15M6MVFdli"
	var verbose bool = true

	var activeNet chaincfg.Params = chaincfg.TestNet3Params

	var fee btcutil.Amount = 50000

}
