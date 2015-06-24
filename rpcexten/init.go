package rpcexten

import "log"

func init() {
	log.Println("I ran damnit")
	registerJsonSendCmds()
	registerWalletSetupCmds()
}
