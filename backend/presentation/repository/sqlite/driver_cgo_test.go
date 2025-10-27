//go:build cgo

package sqlite

import _ "github.com/mattn/go-sqlite3"

const testDriverName = "sqlite3"
