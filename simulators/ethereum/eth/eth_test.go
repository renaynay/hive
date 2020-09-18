package eth

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/hive/simulators/common"
	"github.com/ethereum/hive/simulators/common/providers/hive"
	"github.com/ethereum/hive/simulators/common/providers/local"
	"os"
	"testing"
)

var (
	genesisFile *string
	chainFile   *string

	host common.TestSuiteHost
	err error
)

func init() {
	hive.Support()
	local.Support()
}

func TestMain(m *testing.M) {
	genesisFile = flag.String("genesis", "", "path to genesis file")
	chainFile = flag.String("chainFile", "", "path to chain.rlp file")

	simProviderType := flag.String("simProvider", "", "the simulation provider type (local|hive)")
	providerConfigFile := flag.String("providerConfig", "", "the config json file for the provider")

	flag.Parse()

	host, err = common.InitProvider(*simProviderType, *providerConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialise provider: %s", err.Error())
	}

	os.Exit(RunTestSuite(m))
}

func RunTestSuite(m *testing.M) int {
	return m.Run()
}

func TestEth(t *testing.T) {
	log.Root().SetHandler(log.StdoutHandler)
	logFile, _ := os.LookupEnv("HIVE_SIMLOG")
	//start the test suite
	testSuite, err := host.StartTestSuite("Eth protocol test suite",
		`TODO`, logFile) // TODO needs description
	if err != nil {
		t.Fatalf("Simulator error. Failed to start test suite. %v ", err)
	}
	defer host.EndTestSuite(testSuite)
	//get all client types required to test
	availableClients, err := host.GetClientTypes()
	if err != nil {
		t.Fatalf("Simulator error. Cannot get client types. %v", err)
	}
	
}


