package pepepow

import (
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/bchain/coins/dash"
)

// NewPepepowParser reuses the Dash parser because PepePow is a Dash-derived chain, but
// overrides the address magic bytes to match the upstream PePe-core chainparams.
func NewPepepowParser(params *chaincfg.Params, cfg *btc.Configuration) *dash.DashParser {
	customParams := *params
	switch params.Net {
	case dash.MainnetMagic:
		customParams.PubKeyHashAddrID = []byte{55} // base58 prefix: P
		customParams.ScriptHashAddrID = []byte{16} // base58 prefix: 7
	case dash.TestnetMagic, dash.RegtestMagic:
		customParams.PubKeyHashAddrID = []byte{140} // base58 prefix: y
		customParams.ScriptHashAddrID = []byte{19}  // base58 prefix: 8 or 9
	}
	return dash.NewDashParser(&customParams, cfg)
}
