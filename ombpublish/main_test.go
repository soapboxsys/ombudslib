package ombpublish_test

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcrpcclient"
	"github.com/soapboxsys/ombudslib/ombpublish"
	"github.com/soapboxsys/ombudslib/ombwire"
)

func TestPublishBulletin(t *testing.T) {
	// NOTE THESE CREDS ARE LEAKED INTO GITHUB
	var rpcuser string = "rpcuser"
	var rpcpass string = "14PkTcQcKj8AsTpljxsExs3Idb2CQGRfHcNcUKjgUTKdM0994ce0eN15M6MVFdli"

	var net *chaincfg.Params = &chaincfg.MainNetParams
	//var fee btcutil.Amount = 50000
	var passphrase string = "passphrase"

	certs, err := ioutil.ReadFile("/home/ubuntu/.btcd/rpc.cert")
	if err != nil {
		t.Fatal(err)
	}

	p := ombpublish.NormalParams(net, passphrase)
	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:18332",
		User:         rpcuser,
		Pass:         rpcpass,
		HTTPPostMode: true,
		Certificates: certs,
	}

	client, err := btcrpcclient.New(connCfg, nil)
	defer client.Shutdown()
	if err != nil {
		t.Fatal(err)
	}

	bltn := ombwire.NewBulletin("Are you off something reciting incantations?", uint64(time.Now().Unix()), nil)

	txid, err := ombpublish.PublishBulletin(client, bltn, p)
	if err != nil || txid == nil {
		t.Fatal(err)
	}

	log.Println("Successfully sent bltn:", txid.String())
}
