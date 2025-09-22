package main

import (
	"github.com/agntcy/dir/sdk/wasm/internal/dir/sdk/add"
)

func init() {
	add.Exports.Add = func(x uint32, y uint32) uint32 {
		return x + y
	}
}

// main is required for the `wasi` target, even if it isn't used.
func main() {}
