//go:build !cgo

package sqlite

import _ "modernc.org/sqlite"

const testDriverName = "sqlite"
