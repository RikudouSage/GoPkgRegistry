package migrations

import "embed"

// Assets contains the SQL migration files bundled with the application.
//
//go:embed *.sql
var Assets embed.FS
