package ombproto

import (
	"time"

	"github.com/btcsuite/btcd/wire"
)

type Endorsement struct {
	txid      *wire.ShaHash
	bid       *wire.ShaHash
	Timestamp *time.Time
}
