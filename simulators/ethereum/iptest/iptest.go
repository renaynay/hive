package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/hive/simulators/common/providers/hive"
)

func main() {
	host := hive.New()

	availableClients, _ := host.GetClientTypes()
	log.Info("Got clients", "clients", availableClients)

	logFile, _ := os.LookupEnv("HIVE_SIMLOG")

	for _, client := range availableClients {
		suiteID, err := host.StartTestSuite("iptest", "ip test", logFile)
		if err != nil {
			log.Error("unable to start test suite: ", err.Error())
			os.Exit(1)
		}

		defer func() {
			if err := host.EndTestSuite(suiteID); err != nil {
				log.Error(fmt.Sprintf("Unable to end test suite: %s", err.Error()), err.Error())
				os.Exit(1)
			}
		}()

		testID, err := host.StartTest(suiteID, "iptest", "iptest")
		if err != nil {
			log.Error("unable to start test: ", err.Error())
			os.Exit(1)
		}

		env := map[string]string{
			"CLIENT": client,
		}
		files := map[string]string{}

		_, ip, _, err := host.GetNode(suiteID, testID, env, files)
		log.Info("IP ADDR", ip.String())
	}
}
