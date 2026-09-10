package starter

import "embed"

// Assets contains the templates used by the CLI adapter.
//
//go:embed all:create
var Assets embed.FS
