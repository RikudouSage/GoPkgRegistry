package migrations

import "embed"

//go:embed *.sql
var Assets embed.FS
