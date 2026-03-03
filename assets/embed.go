package assets

import "embed"

// FS embeds all campaign and story content JSON files.
//
//go:embed levels/*.json story/*.json
var FS embed.FS
