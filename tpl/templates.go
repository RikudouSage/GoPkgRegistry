package tpl

import "embed"

// Templates contains the HTML templates bundled with the application.
//
//go:embed *.tpl
var Templates embed.FS
