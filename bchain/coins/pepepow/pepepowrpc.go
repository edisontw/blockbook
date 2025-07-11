package pepepow

import (
    "time"

    "github.com/trezor/blockbook/bchain"
    hclient "github.com/trezor/blockbook/bchain/client/http"
)

// init registers the PepePow RPC client constructor using the generic HTTP JSON-RPC client.
func init() {
    bchain.RegisterRPC("pepepow", parsePepePowRPC)
}

// parsePepePowRPC constructs an HTTP RPC client for PepePow based on RawChaincfg.
func parsePepePowRPC(raw *bchain.RawChaincfg) (*bchain.RPCClient, error) {
    // Use timeout from config (in seconds)
    timeout := time.Duration(raw.RPCTimeout) * time.Second
    // Create a new HTTP JSON-RPC client
    client, err := hclient.New(raw.RPCURL, raw.RPCUser, raw.RPCPass, timeout, raw.NamedParams)
    if err != nil {
        return nil, err
    }
    return client, nil
}
