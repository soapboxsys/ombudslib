// This package handles sending bulletins via a standard bitcoin wallet
package ombpublish

import (
	"errors"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcrpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/davecgh/go-spew/spew"
	"github.com/soapboxsys/ombudslib/ombwire"
)

var defaultDustAmnt = btcutil.Amount(600)
var defaultSatPerByte = btcutil.Amount(175)

func NormalParams(net *chaincfg.Params, passphrase string) Params {
	return Params{
		MinSatToSpend: btcutil.Amount(100000),
		DustAmnt:      defaultDustAmnt,
		SatPerByte:    defaultSatPerByte,
		activeNet:     net,
		passphrase:    passphrase,
		verbose:       true,
	}
}

// Params to pass to Publish* for parameterizing the send
type Params struct {
	MinSatToSpend btcutil.Amount // The floor for the cost of sending a single msg
	DustAmnt      btcutil.Amount
	SatPerByte    btcutil.Amount // The fee to pay per byte
	passphrase    string         // The wallets passphrase
	activeNet     *chaincfg.Params
	verbose       bool
}

// PublishBulletin uses the passed client
func PublishBulletin(client *btcrpcclient.Client, bltn *ombwire.Bulletin, params Params) (*wire.ShaHash, error) {

	ulst, err := client.ListUnspent()
	if err != nil {
		return nil, err
	}
	if len(ulst) < 1 {
		return nil, errors.New("No unspent outputs")
	}

	// Find the total spendable coins that fnc will send with Tx
	sendAmnt := btcutil.Amount(0)
	// The final list of unspents to use
	unspentsToUse := []btcjson.ListUnspentResult{}
	for _, unspent := range ulst {
		if unspent.Spendable && sendAmnt < params.MinSatToSpend {
			sendAmnt += btcutil.Amount(unspent.Amount)
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

	// Use wire helper func to build bltns TxOuts
	txOuts, err := bltn.TxOuts(int64(params.DustAmnt), params.activeNet)
	if err != nil {
		return nil, fmt.Errorf("Encoding bltn failed: %s", err)
	}

	// Form the output side (dust+txouts and change)
	for _, txOut := range txOuts {
		msgtx.AddTxOut(txOut)
	}

	// Create change Addr
	addrStr := unspentsToUse[0].Address
	changeAddr, err := btcutil.DecodeAddress(addrStr, params.activeNet)
	if err != nil {
		return nil, err
	}

	// Determine change to send from amount being sent and tx cost
	change := sendAmnt - determineCost(msgtx.SerializeSize(), len(txOuts), params)

	// Add Change TxOut to tx
	txOut := wire.NewTxOut(int64(change), changeAddr.ScriptAddress())
	msgtx.AddTxOut(txOut)

	if params.verbose {
		spew.Printf("MsgTx Pre-Sig: %s\n", msgtx)
	}

	// Unlock the wallet for 15 seconds
	err = client.WalletPassphrase(params.passphrase, 15)
	if err != nil {
		return nil, fmt.Errorf("Unlocking wallet threw: %s", err)
	}

	// Get the TX signed
	var ok bool
	msgtx, ok, err = client.SignRawTransaction(msgtx)
	if err != nil || !ok {
		log.Println("Signing tx was not ok: %s, %s", ok, err)
		return nil, err
	}

	log.Println("Signed the tx!")
	if params.verbose {
		spew.Printf("MsgTx Post-Sig: %s\n", msgtx)
	}

	// Submit it to the network.
	txid, err := client.SendRawTransaction(msgtx, false)
	if err != nil {
		return nil, fmt.Errorf("Sending Tx failed: %s", err)
	}
	log.Println("Succesfully broadcast Tx[%s]", txid.String())

	return txid, nil
}
func amnt(a int) btcutil.Amount {
	return btcutil.Amount(a)
}

func determineCost(txSizeEst, numTxOuts int, params Params) btcutil.Amount {
	dustSum := amnt(numTxOuts) * params.DustAmnt
	fee := amnt(txSizeEst) * params.SatPerByte
	return dustSum + fee
}
