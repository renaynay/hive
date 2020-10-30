package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/hive/simulators/common/providers/hive"
)

func main() {
	host := hive.New()

	availableClients, err := host.GetClientTypes()
	if err != nil {
		fatalf("failed to get client types: %s", err.Error())
	}
	log.Info("Got clients", "clients", availableClients)

	logFile, _ := os.LookupEnv("HIVE_SIMLOG")

	for _, client := range availableClients {
		suiteID, err := host.StartTestSuite("iptest", "ip test", logFile)
		if err != nil {
			fatalf("failed to start test suite: %s", err.Error())
		}
		testID, err := host.StartTest(suiteID, "iptest", "iptest")
		if err != nil {
			fatalf("failed to start test: %s", err.Error())
		}
		env := map[string]string{
			"CLIENT": client,
		}
		files := map[string]string{}
		// get client
		containerID, ip, _, err := host.GetNode(suiteID, testID, env, files)
		if err != nil {
			fatalf("could not get node: %s", err.Error())
		}
		// create network1
		networkID, err := host.CreateNetwork(suiteID, "network1")
		if err != nil {
			fatalf("could not create network: %s", err.Error())
		}
		// connect client to network1
		if err := host.ConnectContainer(suiteID, networkID, containerID); err != nil {
			fatalf("could not connect container to network: %s", err.Error())
		}
		// connect sim to network
		if err := host.ConnectContainer(suiteID, networkID, "simulation"); err != nil {
			fatalf("could not connect container to network: %s", err.Error())
		}
		// get client ip
		clientIP, err := host.GetContainerNetworkIP(suiteID, networkID, containerID)
		if err != nil {
			fatalf("could not get client network ip address: %s", err.Error())
		}
		_, err = host.GetContainerNetworkIP(suiteID, networkID, "simulation")
		if err != nil {
			fatalf("could not get client network ip address for simulation container: %s", err.Error())
		}

		// TODO dial IPs
		_, err = net.Dial("tcp", fmt.Sprintf("%s:%d", clientIP, 8545))
		if err != nil {
			fatalf("failed to dial client: %s", err.Error())
		}

		// disconnect client from network1
		if err := host.DisconnectContainer(suiteID, networkID, containerID); err != nil {
			fatalf("could not disconnect container from network: %s", err.Error())
		}
		// disconnect simulation from network1
		if err := host.DisconnectContainer(suiteID, networkID, "simulation"); err != nil {
			fatalf("could not disconnect container from network: %s", err.Error())
		}
		// remove network1
		if err := host.RemoveNetwork(suiteID, networkID); err != nil {
			fatalf("could not remove network: %s", err.Error())
		}

		host.KillNode(suiteID, testID, containerID)
		host.EndTest(suiteID, testID, nil, nil) // &common.TestResult{Pass: true, Details: fmt.Sprint("clientIP: %s", clientIP)} // TODO
		host.EndTestSuite(suiteID)
	}
}

func fatalf(format string, args ...interface{}) {
	fatal(fmt.Errorf(format, args...))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
