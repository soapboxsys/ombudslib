// A script that retrieves a given txid from the insight api and attempts to
// decode it as an Ombuds record. It throws an error if it is unsuccessful.
// Usage is simple: build the program and then run it with:
// > ./decode -testnet 65d5a50a90255447973f5b32966c0c80192e80817e38f9b45ff7ddaea16cbbe5
package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/btcsuite/btcd/wire"
	"github.com/davecgh/go-spew/spew"
	"github.com/soapboxsys/ombudslib/ombwire"
)

var t = flag.String("txid", "", "The txid of the transaction you want to decode")
var n = flag.Bool("testnet", false, "The testnet flag")

// A representation of the json returned by the insight api
type bitPayRawTx struct {
	Rawtx string `json:"rawtx"`
}

func main() {
	flag.Parse()

	if *t == "" && flag.Arg(0) == "" {
		log.Fatal("No txid provided")
	}
	if *t == "" {
		*t = flag.Arg(0)
	}

	_, err := wire.NewShaHashFromStr(*t)
	if err != nil {
		log.Fatalf("There was a problem with the txid you provided: %s", err)
	}

	// Get the transaction from bitpay's api
	net := "insight.bitpay.com"
	if *n {
		net = "test-insight.bitpay.com"
	}

	raw_url := fmt.Sprintf("https://%s/api/rawtx/%s", net, *t)
	resp, err := http.Get(raw_url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Problem retrieving txid[%s] got: %s",
			*t, raw_url, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	txJson := &bitPayRawTx{}
	if err = json.Unmarshal(b, txJson); err != nil {
		log.Fatalf("Problem decoding json: %s", err)
	}

	tx_bytes, err := hex.DecodeString(txJson.Rawtx)
	if err != nil {
		log.Fatal(err)
	}

	// Deserialize the bytes into a bitcoin tx
	tx := wire.NewMsgTx()
	if err := tx.Deserialize(bytes.NewBuffer(tx_bytes)); err != nil {
		log.Fatal(err)
	}

	// Decode the Bitcoin transaction into a Ombuds record.
	record, err := ombwire.ParseTx(tx)
	if err != nil {
		log.Fatalf("Decoding wire type failed with: %s", err)
	}

	spew.Printf("SUCCESS! Decoded Ombuds Record:\n%s\n", record)
}
