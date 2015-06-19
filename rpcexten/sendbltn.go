package rpcexten

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	newjson "github.com/btcsuite/btcd/btcjson/v2/btcjson"
)

var (
	sendbltnMeth    = "sendbulletin"
	composebltnMeth = "composebulletin"
)

type BulletinCmd interface {
	GetAddress() string
	GetMessage() string
	GetBoard() string
}

// NOTE any changes to the sendbulletin api must be reflected in compose as well.
type SendBulletinCmd struct {
	id      interface{}
	Address string
	Board   string
	Message string
}

// This was added to handle any interfaces that were switched in v0.10.0
type SendBulletinCmdv2 struct {
	Address string
	Board   string
	Message string
}

func NewSendBulletinCmd(id interface{}, address, board, message string) *SendBulletinCmd {
	return &SendBulletinCmd{
		id:      id,
		Address: address,
		Message: message,
		Board:   board,
	}
}

func (cmd SendBulletinCmd) GetAddress() string {
	return cmd.Address
}

func (cmd SendBulletinCmd) GetBoard() string {
	return cmd.Board
}

func (cmd SendBulletinCmd) GetMessage() string {
	return cmd.Message
}

func (cmd SendBulletinCmd) Id() interface{} {
	return cmd.id
}

func (cmd SendBulletinCmd) Method() string {
	return sendbltnMeth
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd SendBulletinCmd) MarshalJSON() ([]byte, error) {
	params := []interface{}{
		cmd.Address,
		cmd.Board,
		cmd.Message,
	}

	// Fill and marshal a RawCmd.
	raw, err := btcjson.NewRawCmd(cmd.id, cmd.Method(), params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

func (cmd SendBulletinCmd) UnmarshalJSON(b []byte) error {
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	if len(r.Params) != 3 {
		return btcjson.ErrWrongNumberOfParams
	}

	address, board, message, err := extractParams(r.Params)
	if err != nil {
		return err
	}

	newCmd := NewSendBulletinCmd(r.Id, address, board, message)

	cmd = *newCmd
	return nil
}

type ComposeBulletinCmd struct {
	id      interface{}
	Address string
	Board   string
	Message string
}

type ComposeBulletinCmdv2 struct {
	Address string
	Board   string
	Message string
}

func NewComposeBulletinCmd(id interface{}, address, board, message string) *ComposeBulletinCmd {
	return &ComposeBulletinCmd{
		id:      id,
		Address: address,
		Message: message,
		Board:   board,
	}
}

func (cmd ComposeBulletinCmd) GetAddress() string {
	return cmd.Address
}

func (cmd ComposeBulletinCmd) GetBoard() string {
	return cmd.Board
}

func (cmd ComposeBulletinCmd) GetMessage() string {
	return cmd.Message
}

func (cmd ComposeBulletinCmd) Id() interface{} {
	return cmd.id
}

func (cmd ComposeBulletinCmd) Method() string {
	return composebltnMeth
}

func extractParams(params []json.RawMessage) (string, string, string, error) {
	throw := func(err error) (string, string, string, error) {
		return "", "", "", err
	}

	var address string
	if err := json.Unmarshal(params[0], &address); err != nil {
		return throw(fmt.Errorf("first parameter 'address' must be a string: %v", err))
	}

	var board string
	if err := json.Unmarshal(params[1], &board); err != nil {
		return throw(fmt.Errorf("second parameter 'board' must be a string: %v", err))
	}

	var message string
	if err := json.Unmarshal(params[2], &message); err != nil {
		return throw(fmt.Errorf("third parameter 'board' must be a string: %v", err))
	}

	return address, board, message, nil
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd ComposeBulletinCmd) MarshalJSON() ([]byte, error) {
	params := []interface{}{
		cmd.Address,
		cmd.Board,
		cmd.Message,
	}

	// Fill and marshal a RawCmd.
	raw, err := btcjson.NewRawCmd(cmd.id, cmd.Method(), params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

func (cmd ComposeBulletinCmd) UnmarshalJSON(b []byte) error {
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	if len(r.Params) != 3 {
		return btcjson.ErrWrongNumberOfParams
	}

	address, board, message, err := extractParams(r.Params)
	if err != nil {
		return err
	}

	newCmd := NewComposeBulletinCmd(r.Id, address, board, message)

	cmd = *newCmd
	return nil
}

// rawCmdParser works for both sendbulletin and composebulletin commands by
// using the method of the RawCmd.
func rawCmdParser(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) != 3 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var address string
	if err := json.Unmarshal(r.Params[0], &address); err != nil {
		return nil, fmt.Errorf("first parameter 'address' must be a string: %v", err)
	}

	var board string
	if err := json.Unmarshal(r.Params[1], &board); err != nil {
		return nil, fmt.Errorf("second parameter 'board' must be a string: %v", err)
	}

	var message string
	if err := json.Unmarshal(r.Params[2], &message); err != nil {
		return nil, fmt.Errorf("third parameter 'board' must be a string: %v", err)
	}

	var cmd btcjson.Cmd

	switch {
	case r.Method == sendbltnMeth:
		cmd = NewSendBulletinCmd(r.Id, address, board, message)
	case r.Method == composebltnMeth:
		cmd = NewComposeBulletinCmd(r.Id, address, board, message)
	default:
		return nil, fmt.Errorf("Wrong method for json cmd: %s", r.Method)
	}

	return cmd, nil
}

func sendReplyParser(rawJ json.RawMessage) (interface{}, error) {
	var txSha string
	err := json.Unmarshal(rawJ, &txSha)
	if err != nil {
		return nil, err
	}
	return txSha, nil
}

func composeReplyParser(rawJ json.RawMessage) (interface{}, error) {
	var rawHex string
	err := json.Unmarshal(rawJ, &rawHex)
	if err != nil {
		return nil, err
	}
	return rawHex, nil
}

func registerJsonCmds() {
	sendHelpStr := sendbltnMeth + " <address> <board> <message>"
	btcjson.RegisterCustomCmd(sendbltnMeth, rawCmdParser, sendReplyParser, sendHelpStr)
	newjson.MustRegisterCmd(sendbltnMeth, (*SendBulletinCmdv2)(nil), newjson.UFWalletOnly)

	composeHelpStr := composebltnMeth + " <address> <board> <message>"
	btcjson.RegisterCustomCmd(composebltnMeth, rawCmdParser, composeReplyParser, composeHelpStr)
	newjson.MustRegisterCmd(composebltnMeth, (*ComposeBulletinCmdv2)(nil), newjson.UFWalletOnly)
}

func init() {
	registerJsonCmds()
}
