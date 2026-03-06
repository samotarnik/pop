package assets

import "embed"

// FS embeds all campaign and story content JSON files.
//
//go:embed levels/*.json story/*.json sprites/prince/*
var FS embed.FS
