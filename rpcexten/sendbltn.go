package rpcexten

import (
	"encoding/json"
	"fmt"

	"github.com/conformal/btcjson"
)

type SendBulletinCmd struct {
	id      interface{}
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

func (cmd SendBulletinCmd) Id() interface{} {
	return cmd.id
}

func (cmd SendBulletinCmd) Method() string {
	return "sendbulletin"
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

	var address string
	if err := json.Unmarshal(r.Params[0], &address); err != nil {
		return fmt.Errorf("first parameter 'address' must be a string: %v", err)
	}

	var board string
	if err := json.Unmarshal(r.Params[0], &address); err != nil {
		return fmt.Errorf("second parameter 'board' must be a string: %v", err)
	}

	var message string
	if err := json.Unmarshal(r.Params[0], &address); err != nil {
		return fmt.Errorf("third parameter 'board' must be a string: %v", err)
	}

	newCmd := NewSendBulletinCmd(r.Id, address, board, message)

	cmd = *newCmd
	return nil
}

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

	cmd := NewSendBulletinCmd(r.Id, address, board, message)
	return *cmd, nil
}

func replyParser(rawJ json.RawMessage) (interface{}, error) {
	var txSha string
	err := json.Unmarshal(rawJ, &txSha)
	if err != nil {
		return nil, err
	}
	return txSha, nil
}

func registerJsonCmd() {
	helpStr := "sendbulletin <address> <board> <message>"
	btcjson.RegisterCustomCmd("sendbulletin", rawCmdParser, replyParser, helpStr)
}

func init() {
	fmt.Println("blah")
	registerJsonCmd()
}

/*func pain() {
	addr, board, content := "thisisanaddr", "ahimsa-dev", "Derp derp derp"
	cmd := NewSendBulletinCmd(float64(1), addr, board, content)

	helpStr := "sendbulletin <address> <board> <message>"

	btcjson.RegisterCustomCmd("sendbulletin", rawCmdParser, replyParser, helpStr)

	msg, err := json.Marshal(cmd)
	if err != nil {
		log.Fatal(err)
	}

	cmd2, err := btcjson.ParseMarshaledCmd(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", cmd2)

	resp := []byte(`{"id":1, "result":"1HB5XMLmzFVj8ALj6mfBsbifRoD4miY36v", "error":null}`)
	_, err = btcjson.ReadResultCmd("getnewaddress", resp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("getnew worked")
	reply, err := btcjson.ReadResultCmd("sendbulletin", resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", reply)
}*/

// The minimum dust value for a PayToPubKey tx accepted by the network
func DustAmnt() int64 {
	return 567
}
