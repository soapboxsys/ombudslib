package rpcexten

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
)

var (
	walletSetupMeth      = "walletsetup"
	walletStateCheckMeth = "walletstatecheck"
)

type WalletSetupCmd struct {
	id         interface{}
	Passphrase string `json:"passphrase"`
	UseSeed    bool   `json:"useSeed"`
	RandSeed   string `json:"randseed"`
}

func NewWalletSetupCmd(id interface{}, phrase string, us bool, seed string) *WalletSetupCmd {
	return &WalletSetupCmd{
		id:         id,
		Passphrase: phrase,
		UseSeed:    us,
		RandSeed:   seed,
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
		cmd.UseSeed,
		cmd.RandSeed,
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

	if len(r.Params) != 3 {
		return btcjson.ErrWrongNumberOfParams
	}

	var passphrase, randSeed string
	var useSeed bool

	if err := json.Unmarshal(r.Params[0], &passphrase); err != nil {
		return err
	}

	if err := json.Unmarshal(r.Params[1], &useSeed); err != nil {
		return err
	}

	if err := json.Unmarshal(r.Params[2], &randSeed); err != nil {
		return err
	}

	newCmd := NewWalletSetupCmd(r.Id, passphrase, useSeed, randSeed)
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
	if len(r.Params) != 3 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var passphrase string
	if err := json.Unmarshal(r.Params[0], &passphrase); err != nil {
		return nil, fmt.Errorf("first parameter 'passphrase' must be a string: %v", err)
	}

	var useSeed bool
	if err := json.Unmarshal(r.Params[1], &useSeed); err != nil {
		return nil, fmt.Errorf("second parameter 'useSeed' must be a bool: %v", err)
	}

	var randSeed string
	if err := json.Unmarshal(r.Params[2], &randSeed); err != nil {
		return nil, fmt.Errorf("third parameter 'randSeed' must be a string: %v", err)
	}

	var cmd btcjson.Cmd
	cmd = NewWalletSetupCmd(r.Id, passphrase, useSeed, randSeed)

	return cmd, nil
}

func registerWalletSetupCmds() {
	walletSetupHelpStr := walletSetupMeth + " <passphrase> [<useSeed> <randseed>]"

	btcjson.RegisterCustomCmd(walletSetupMeth, rawSetupCmdParser, walletSetupReplyParser, walletSetupHelpStr)
}
