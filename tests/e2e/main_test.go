package e2e

import (
	"flag"
	"os"
	"testing"

	appparams "github.com/babylonchain/babylon-sdk/demo/app/params"
)

func TestMain(m *testing.M) {
	flag.StringVar(&wasmContractPath, "contracts-path", "../testdata", "Set path to dir with wasm contracts")
	flag.BoolVar(&wasmContractGZipped, "gzipped", false, "Use `.gz` file ending when set")
	flag.Parse()

	// overwrite init messages in Babylon
	appparams.SetAddressPrefixes()

	// Run the tests
	exitCode := m.Run()

	// Teardown code here (if needed)

	// Exit with the appropriate code
	os.Exit(exitCode)
}
