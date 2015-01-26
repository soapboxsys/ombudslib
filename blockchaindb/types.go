package blockchaindb

// A struct that exposes the amount of fuel, the number of confirmations, and
// the outpoint of the txout that exposes that fuel. This struct is annotated
// for json encoding.
type FuelOut struct {
	Depth    uint64   `json:"depth"`
	Fuel     uint32   `json:"fuel"`
	OutPoint OutPoint `json:"outpoint"`
}

func NewFuelOut() FuelOut {
	return FuelOut{Depth: 10, Fuel: 25, OutPoint: OutPoint{2, "deadbeefdeadbeefdeadbeef"}}
}

// A struct holding the position of the txout in the blockchain.
type OutPoint struct {
	Index uint32 `json:"idx"`
	Hash  string `json:"hash"`
}
