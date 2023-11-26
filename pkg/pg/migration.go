package pg

import "embed"

//go:embed migrations/*.sql
var ReadModelFS embed.FS
