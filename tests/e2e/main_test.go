package e2e

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.StringVar(&wasmContractPath, "contracts-path", "../testdata", "Set path to dir with wasm contracts")
	flag.BoolVar(&wasmContractGZipped, "gzipped", false, "Don't use `.gz` file ending when set")
	flag.Parse()

	os.Exit(m.Run())
}
