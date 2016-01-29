// This script connects to a Bitcoin peer requests a block based on the
// "target" block hash and attempts to write it to a file.
package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/davecgh/go-spew/spew"
)

var node = flag.String("node", "", "The node to connect to")
var target = flag.String("target", "", "The bitcoin block to get")
var testnet = flag.Bool("testnet", false, "The Testnet3 flag")

var bnet = wire.MainNet
var pver = wire.ProtocolVersion

func main() {
	flag.Parse()

	if *testnet {
		bnet = wire.TestNet3
	}

	conn, err := net.DialTimeout("tcp", *node, 500*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	// Deal with the Bitcoin handshake
	performHandShake(conn)

	writer := composeWrite(conn, wire.ProtocolVersion, bnet)

	err = writer(makeBlkDataMsg(*target))
	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, _, err := wire.ReadMessage(conn, pver, bnet)
		if err != nil {
			log.Fatal(err)
		}

		switch resp := resp.(type) {
		case *wire.MsgBlock:
			spew.Printf("Got block: %s\n", resp.BlockSha())
			n := "blk-" + *node
			f, err := os.Create(n)
			if err != nil {
				log.Fatal(err)
			}
			i, err := wire.WriteMessageN(f, resp, pver, bnet)
			if err != nil {
				log.Fatal(err)
			}
			log.Fatalf("Wrote %s bytes to: %s", i, n)

		default:
			spew.Printf("Got: %s: %v\n", resp.Command(), resp)
		}

	}
}

func composeWrite(conn net.Conn, pver uint32, net wire.BitcoinNet) func(wire.Message) error {
	return func(msg wire.Message) error {
		return wire.WriteMessage(conn, msg, pver, net)
	}
}

func genNonce() uint64 {
	n, _ := wire.RandomUint64()
	return n
}

func performHandShake(conn net.Conn) {
	writer := composeWrite(conn, wire.ProtocolVersion, bnet)

	nonce := genNonce()
	ver_m, _ := wire.NewMsgVersionFromConn(conn, nonce, 0)
	ver_m.AddUserAgent("blockget", "0.0.1")

	err := writer(ver_m)
	if err != nil {
		log.Fatal(err)
	}

	resp, _, err := wire.ReadMessage(conn, wire.ProtocolVersion, bnet)
	if err != nil {
		log.Fatal(err)
	}

	verack := wire.NewMsgVerAck()
	if err = writer(verack); err != nil {
		log.Fatal(err)
	}

	resp, _, err = wire.ReadMessage(conn, wire.ProtocolVersion, bnet)
	if err != nil {
		log.Fatal(err)
	}
	spew.Printf("Resp: %s: %v\n", resp.Command(), resp)
}

func makeBlkDataMsg(h string) *wire.MsgGetData {
	hash, err := wire.NewShaHashFromStr(h)
	if err != nil {
		log.Fatal(err)
	}

	gd := wire.NewMsgGetData()
	iv := wire.NewInvVect(wire.InvTypeBlock, hash)
	gd.AddInvVect(iv)

	return gd
}
