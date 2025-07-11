package pepepow

import (
    "github.com/trezor/blockbook/bchain"
    "github.com/trezor/blockbook/bchain/coins/dash"
)

// init registers the PepePow parser using the Dash parser as base.
func init() {
    bchain.RegisterParser("pepepow", parsePepePowParams)
}

// parsePepePowParams delegates to Dash parser because PepePow is a fork of Dash core v0.12.2.
// It then overrides any chain-specific parameters from the raw configuration.
func parsePepePowParams(raw *bchain.RawChaincfg) (*bchain.Chaincfg, error) {
    // Use the Dash parser for most settings
    cfg, err := dash.Parse(raw)
    if err != nil {
        return nil, err
    }

    // Override XPub magic values specific to PepePow
    cfg.XpubMagic = raw.XpubMagic
    cfg.XpubMagicSegwitP2sh = raw.XpubMagicSegwitP2sh
    cfg.XpubMagicSegwitNative = raw.XpubMagicSegwitNative

    return cfg,
