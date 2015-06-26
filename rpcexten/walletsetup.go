package rpcexten

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	newjson "github.com/btcsuite/btcd/btcjson/v2/btcjson"
)

var (
	walletSetupMeth    = "walletsetup"
	getWalletStateMeth = "getwalletstate"
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

func (cmd *WalletSetupCmd) UnmarshalJSON(b []byte) error {
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
	*cmd = *newCmd
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

type GetWalletStateCmd struct {
	id interface{}
}

type GetWalletStateCmdv2 struct{}

func NewGetWalletStateCmd(id interface{}) *GetWalletStateCmd {
	return &GetWalletStateCmd{id: id}
}

func (cmd GetWalletStateCmd) Method() string {
	return getWalletStateMeth
}

func (cmd GetWalletStateCmd) Id() interface{} {
	return cmd.id
}

func (cmd GetWalletStateCmd) MarshalJSON() ([]byte, error) {
	raw, err := btcjson.NewRawCmd(cmd.id, cmd.Method(), []interface{}{})
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(raw)

}

func (cmd *GetWalletStateCmd) UnmarshalJSON(b []byte) error {
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	if len(r.Params) != 0 {
		return btcjson.ErrWrongNumberOfParams
	}

	cmd.id = r.Id
	return nil
}

// using the method of the RawCmd.
func rawGetWalletStateParser(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 0 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var cmd btcjson.Cmd
	cmd = NewGetWalletStateCmd(r.Id)
	return cmd, nil
}

func getStateReplyParser(rawJ json.RawMessage) (interface{}, error) {
	var res GetWalletStateResult
	err := json.Unmarshal(rawJ, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type GetWalletStateResult struct {
	HasWallet   bool `json:"hasWallet"`
	HasChainSvr bool `json:"hasChainSvr"`
	ChainSynced bool `json:"chainSynced"`
}

func registerWalletSetupCmds() {
	walletSetupHelpStr := walletSetupMeth + " <passphrase>"

	btcjson.RegisterCustomCmd(walletSetupMeth, rawSetupCmdParser, walletSetupReplyParser, walletSetupHelpStr)
	newjson.MustRegisterCmd(walletSetupMeth, (*WalletSetupCmdv2)(nil), newjson.UFWalletOnly)

	btcjson.RegisterCustomCmd(getWalletStateMeth, rawGetWalletStateParser, getStateReplyParser, getWalletStateMeth)
	newjson.MustRegisterCmd(getWalletStateMeth, (*GetWalletStateCmdv2)(nil), newjson.UFWalletOnly)
}
