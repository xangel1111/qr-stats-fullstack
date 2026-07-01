// Package db embeds the SQL migration files into the binary so the service is
// fully self-contained: no external migration files need to ship alongside it.
package db

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS
