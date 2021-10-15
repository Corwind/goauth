package routes

import "github.com/apple/foundationdb/bindings/go/src/fdb"

type Env struct {
	DB fdb.Database
}
