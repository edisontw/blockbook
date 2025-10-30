package pepepow

import (
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
	"github.com/trezor/blockbook/bchain/coins/dash"
)

const (
	// xpub defaults copied from Bitcoin as PepePow reuses the same serialization.
	xpubMagic             uint32 = 0x0488b21e
	xpubMagicSegwitP2sh   uint32 = 0x049d7cb2
	xpubMagicSegwitNative uint32 = 0x04b24746
)

// PepepowParser implements the bitcoin-like parser logic for PepePow.
type PepepowParser struct {
	*btc.BitcoinLikeParser
	baseparser *bchain.BaseParser
}

func applyPepepowConfig(c *btc.Configuration) *btc.Configuration {
	if c == nil {
		c = &btc.Configuration{}
	}
	if c.XPubMagic == 0 {
		c.XPubMagic = xpubMagic
	}
	if c.XPubMagicSegwitP2sh == 0 {
		c.XPubMagicSegwitP2sh = xpubMagicSegwitP2sh
	}
	if c.XPubMagicSegwitNative == 0 {
		c.XPubMagicSegwitNative = xpubMagicSegwitNative
	}
	return c
}

// NewPepepowParser returns a new parser instance configured for PepePow.
func NewPepepowParser(params *chaincfg.Params, c *btc.Configuration) *PepepowParser {
	cfg := applyPepepowConfig(c)
	return &PepepowParser{
		BitcoinLikeParser: btc.NewBitcoinLikeParser(params, cfg),
		baseparser:        &bchain.BaseParser{},
	}
}

// GetChainParams proxies to the Dash chain parameters as PepePow inherits them.
func GetChainParams(chain string) *chaincfg.Params {
	return dash.GetChainParams(chain)
}

// PackTx packs transaction to byte array using protobuf.
func (p *PepepowParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array.
func (p *PepepowParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
