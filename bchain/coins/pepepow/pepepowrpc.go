package pepepow

import (
    "encoding/json"

    "github.com/trezor/blockbook/bchain"
    "github.com/trezor/blockbook/bchain/coins/btc"
)

// -----------------------------------------------------------------------------
//  PepePow RPC driver – 100 % Dash‑compatible but **no** DIP / special‑tx logic
// -----------------------------------------------------------------------------
// PepePOW is forked from Dash 0.12.2.*. That version predates DIP0001, so there
// are *no* SpecialTransactions and we do **not** need the extra GetBlock logic
// found in newer Dash (or Blockbook’s dashrpc.go). We therefore rely entirely on
// BitcoinRPC’s native implementation and only tweak the RPC marshaller so that
// params are sent as **positional arrays** (JSONMarshalerV1).
// -----------------------------------------------------------------------------

type PepepowRPC struct {
    *btc.BitcoinRPC
}

// NewPepepowRPC constructs the RPC backend, forces array‑style params, and is
// registered via init() below so that chaincfg can reference "pepepow".
func NewPepepowRPC(cfg json.RawMessage, push func(bchain.NotificationType)) (bchain.BlockChain, error) {
    base, err := btc.NewBitcoinRPC(cfg, push)
    if err != nil {
        return nil, err
    }

    r := &PepepowRPC{base.(*btc.BitcoinRPC)}
    r.RPCMarshaler = btc.JSONMarshalerV1{} // positional params for legacy core
    return r, nil
}

func init() {
    bchain.RegisterRPC("pepepow", NewPepepowRPC)
}

// Initialize chooses the correct chain params and boots the parser.
func (r *PepepowRPC) Initialize() error {
    ci, err := r.GetChainInfo()
    if err != nil {
        return err
    }
    params := GetChainParams(ci.Chain)

    r.Parser = NewPepepowParser(params, r.ChainConfig)
    if params.Net == MainnetMagic {
        r.Testnet = false
        r.Network = "livenet"
    } else {
        r.Testnet = true
        r.Network = "testnet"
    }
    return nil
}

// GetTransactionForMempool simply proxies to the underlying RPC.
func (r *PepepowRPC) GetTransactionForMempool(txid string) (*bchain.Tx, error) {
    return r.GetTransaction(txid)
}
