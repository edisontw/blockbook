package pepepow

import (
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/bchain/coins/dash"
)

// NewPepepowParser reuses the Dash parser because PepePow is a Dash-derived chain.
func NewPepepowParser(params *chaincfg.Params, cfg *btc.Configuration) *dash.DashParser {
	return dash.NewDashParser(params, cfg)
}
