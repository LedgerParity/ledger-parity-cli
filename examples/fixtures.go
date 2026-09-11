package examples

import "embed"

// Files are synthetic offline fixtures, not observed ledger data.
//
//go:embed config.json internal.json onchain.json sdp/*
var Files embed.FS
