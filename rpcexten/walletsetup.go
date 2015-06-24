package rpcexten

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	newjson "github.com/btcsuite/btcd/btcjson/v2/btcjson"
)

var (
	walletSetupMeth      = "walletsetup"
	walletStateCheckMeth = "walletstatecheck"
)

type WalletSetupCmd struct {
	id         interface{}
	Passphrase string `json:"passphrase"`
}

// For ombctl
type WalletSetupCmdv2 struct {
	Passphrase string
}

func NewWalletSetupCmd(id interface{}, phrase string) *WalletSetupCmd {
	return &WalletSetupCmd{
		id:         id,
		Passphrase: phrase,
	}
}

func (cmd WalletSetupCmd) Id() interface{} {
	return cmd.id
}

func (cmd WalletSetupCmd) Method() string {
	return walletSetupMeth
}

func (cmd WalletSetupCmd) MarshalJSON() ([]byte, error) {
	params := []interface{}{
		cmd.Passphrase,
	}
	raw, err := btcjson.NewRawCmd(cmd.id, cmd.Method(), params)
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(raw)
}

func (cmd WalletSetupCmd) UnmarshalJSON(b []byte) error {
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	if len(r.Params) != 1 {
		return btcjson.ErrWrongNumberOfParams
	}

	var passphrase string

	if err := json.Unmarshal(r.Params[0], &passphrase); err != nil {
		return err
	}

	newCmd := NewWalletSetupCmd(r.Id, passphrase)
	cmd = *newCmd
	return nil
}

func walletSetupReplyParser(rawJ json.RawMessage) (interface{}, error) {
	var msg string
	err := json.Unmarshal(rawJ, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

type WalletStateCheck struct {
	id          interface{}
	hasChainSvr bool `json:"hasChainSvr"`
	hasWallet   bool `json:"hasWallet"`
}

// using the method of the RawCmd.
func rawSetupCmdParser(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 1 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var passphrase string
	if err := json.Unmarshal(r.Params[0], &passphrase); err != nil {
		return nil, fmt.Errorf("first parameter 'passphrase' must be a string: %v", err)
	}

	var cmd btcjson.Cmd
	cmd = NewWalletSetupCmd(r.Id, passphrase)

	return cmd, nil
}

func registerWalletSetupCmds() {
	walletSetupHelpStr := walletSetupMeth + " <passphrase>"

	btcjson.RegisterCustomCmd(walletSetupMeth, rawSetupCmdParser, walletSetupReplyParser, walletSetupHelpStr)
	newjson.MustRegisterCmd(walletSetupMeth, (*WalletSetupCmdv2)(nil), newjson.UFWalletOnly)
}
