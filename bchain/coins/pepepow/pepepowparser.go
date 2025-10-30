package pepepow

import (
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/bchain/coins/dash"
)

const (
	// MainnetMagic mirrors Dash's magic constant for network detection.
	MainnetMagic = dash.MainnetMagic
	// TestnetMagic mirrors Dash's testnet magic constant.
	TestnetMagic = dash.TestnetMagic
	// RegtestMagic mirrors Dash's regtest magic constant.
	RegtestMagic = dash.RegtestMagic
)

var (
	mainNetParams chaincfg.Params
	testNetParams chaincfg.Params
	regtestParams chaincfg.Params
)

func init() {
	mainNetParams = *dash.GetChainParams("main")
	mainNetParams.Name = "Pepepow"
	mainNetParams.Net = MainnetMagic
	mainNetParams.PubKeyHashAddrID = []byte{55} // base58 prefix: P
	mainNetParams.ScriptHashAddrID = []byte{56} // base58 prefix: 3

	testNetParams = *dash.GetChainParams("test")
	testNetParams.Name = "Pepepow testnet"
	testNetParams.Net = TestnetMagic

	regtestParams = *dash.GetChainParams("regtest")
	regtestParams.Name = "Pepepow regtest"
	regtestParams.Net = RegtestMagic
}

// GetChainParams returns the Pepepow-specific chain parameters for the selected network.
func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&mainNetParams) {
		if err := chaincfg.Register(&mainNetParams); err != nil && err != chaincfg.ErrDuplicateNet {
			panic(err)
		}
		if err := chaincfg.Register(&testNetParams); err != nil && err != chaincfg.ErrDuplicateNet {
			panic(err)
		}
		if err := chaincfg.Register(&regtestParams); err != nil && err != chaincfg.ErrDuplicateNet {
			panic(err)
		}
	}

	switch chain {
	case "test":
		return &testNetParams
	case "regtest":
		return &regtestParams
	default:
		return &mainNetParams
	}
}

// NewPepepowParser creates a Dash-compatible parser with Pepepow specific overrides applied.
func NewPepepowParser(params *chaincfg.Params, cfg *btc.Configuration) (*dash.DashParser, error) {
	parser := dash.NewDashParser(params, cfg)

	parser.BitcoinLikeParser.XPubMagic = cfg.XPubMagic
	parser.BitcoinLikeParser.XPubMagicSegwitP2sh = cfg.XPubMagicSegwitP2sh
	parser.BitcoinLikeParser.XPubMagicSegwitNative = cfg.XPubMagicSegwitNative
	parser.BitcoinLikeParser.Slip44 = cfg.Slip44

	return parser, nil
}
