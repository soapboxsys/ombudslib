// This script sends a bulletin via wallet RPC commands
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/davecgh/go-spew/spew"
	"github.com/soapboxsys/ombudslib/ombwire"
)

var rpcuser string = "rpcuser"
var rpcpass string = "14PkTcQcKj8AsTpljxsExs3Idb2CQGRfHcNcUKjgUTKdM0994ce0eN15M6MVFdli"
var verbose bool = true

var activeNet chaincfg.Params = chaincfg.TestNet3Params

var fee btcutil.Amount = 50000

func panicWithMsg(msg string) {
	log.Fatalf("Killed execution because: %s\n", msg)
}

func panicOn(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s\n", msg, err)
	}
}

func main() {

	//pipe := false
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is being piped to stdin")
		//pipe = true
	} else {
		fmt.Println("stdin is from a terminal")
	}

	client := createClient()
	defer client.Shutdown()

	ulst, err := client.ListUnspent()
	panicOn(err, "ListUnspent threw")
	if len(ulst) < 1 {
		panicWithMsg("No unspent outputs")
	}

	// Deduce amount from msg size.
	requiredSats := btcutil.Amount(100000)
	// Find the total spendable coins
	spendAmnt := btcutil.Amount(0)
	// The final list of unspents to use
	unspentsToUse := []btcjson.ListUnspentResult{}
	for _, unspent := range ulst {
		if unspent.Spendable && spendAmnt < requiredSats {
			spendAmnt += btcutil.Amount(unspent.Amount)
			unspentsToUse = append(unspentsToUse, unspent)
		}
	}
	// Take those unspents create the input side
	msgtx := wire.NewMsgTx()

	for i, unspent := range unspentsToUse {
		empt := []byte{}
		txid, _ := wire.NewShaHashFromStr(unspent.TxID)
		outpoint := wire.OutPoint{
			Hash:  *txid,
			Index: uint32(i),
		}
		txIn := wire.NewTxIn(&outpoint, empt)
		msgtx.AddTxIn(txIn)
	}

	wireBltn := ombwire.NewBulletin("This is a sample message", 12345, nil)
	dustAmnt := int64(567)
	txOuts, err := wireBltn.TxOuts(dustAmnt, &activeNet)
	panicOn(err, "Encoding bltn failed.")
	// Form the output side (dust+txouts and change)
	dustSum := btcutil.Amount(567 * len(txOuts))
	for _, txOut := range txOuts {
		msgtx.AddTxOut(txOut)
	}

	// Create change Addr
	addrStr := unspentsToUse[0].Address
	changeAddr, err := btcutil.DecodeAddress(addrStr, &activeNet)
	panicOn(err, "Change address could not be decoded")

	// Determine change to send
	change := int64(spendAmnt - (dustSum + fee))

	// Add Change TxOut to tx
	txOut := wire.NewTxOut(change, changeAddr.ScriptAddress())
	msgtx.AddTxOut(txOut)

	if verbose {
		spew.Printf("MsgTx Pre-Sig: %s\n", msgtx)
	}

	// Get the TX signed
	var ok bool
	msgtx, ok, err = client.SignRawTransaction(msgtx)
	panicOn(err, "Signing failed with")
	if !ok {
		log.Fatalf("Signing tx was not ok: %s, %s", ok, err)
	}

	log.Println("Signed the tx!")
	if verbose {
		spew.Printf("MsgTx Post-Sig: %s\n", msgtx)
	}

	// Submit it to the network.
	// TODO NOTICE
	if !verbose {
		txid, err := client.SendRawTransaction(msgtx, false)
		panicOn(err, "Sending Tx failed")
		log.Println("Broadcast tx[%s]", txid.String())
	}
}

func createClient() *btcrpcclient.Client {
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}

	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:18332",
		User:         rpcuser,
		Pass:         rpcpass,
		HTTPPostMode: true,
		Certificates: certs,
	}

	client, err := btcrpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}

	return client
}
