package rpcexten

import "github.com/btcsuite/btcd/btcjson"

var (
	sendbltnMeth    = "sendbulletin"
	composebltnMeth = "composebulletin"
)

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

func registerJsonSendCmds() {
	btcjson.MustRegisterCmd(sendbltnMeth, (*SendBulletinCmdv2)(nil), btcjson.UFWalletOnly)

	btcjson.MustRegisterCmd(composebltnMeth, (*ComposeBulletinCmdv2)(nil), btcjson.UFWalletOnly)
}
