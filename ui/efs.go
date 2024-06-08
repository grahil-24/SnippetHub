package ui

import "embed"

// The comment directive must be placed immediately above the variable in which you want
// to store the embedded files.
//
//go:embed "html" "static"
var Files embed.FS
