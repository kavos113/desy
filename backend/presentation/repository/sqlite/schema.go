package sqlite

import (
	_ "embed"
	"strings"
)

//go:embed schema.sql
var embeddedSchema string

func schemaStatements() []string {
	raw := strings.Split(embeddedSchema, ";")
	statements := make([]string, 0, len(raw))
	for _, stmt := range raw {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		statements = append(statements, trimmed)
	}
	return statements
}
