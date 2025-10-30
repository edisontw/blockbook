package pepepow

import (
	"encoding/json"

	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/bchain/coins/dash"
)

// --------------------------------------------------------------------------------
// PepePow RPC driver – positional params only (Dash v0.12.2 compatibility)
// --------------------------------------------------------------------------------
type PepepowRPC struct {
	*btc.BitcoinRPC
}

// NewPepepowRPC constructs the RPC backend and forces array-style params.
func NewPepepowRPC(cfg json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	base, err := btc.NewBitcoinRPC(cfg, pushHandler)
	if err != nil {
		return nil, err
	}
	rpc := &PepepowRPC{base.(*btc.BitcoinRPC)}
	rpc.RPCMarshaler = btc.JSONMarshalerV1{} // positional params
	return rpc, nil
}

// Initialize selects network params and attaches the parser.
func (r *PepepowRPC) Initialize() error {
	ci, err := r.GetChainInfo()
	if err != nil {
		return err
	}
	params := dash.GetChainParams(ci.Chain)

	r.Parser = NewPepepowParser(params, r.ChainConfig)
	if params.Net == dash.MainnetMagic {
		r.Testnet = false
		r.Network = "livenet"
	} else {
		r.Testnet = true
		r.Network = "testnet"
	}
	return nil
}

// GetTransactionForMempool proxies to the underlying RPC.
func (r *PepepowRPC) GetTransactionForMempool(txid string) (*bchain.Tx, error) {
	return r.GetTransaction(txid)
}
