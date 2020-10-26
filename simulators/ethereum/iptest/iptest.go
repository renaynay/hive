package main

import (
	"fmt"
	"os"

	"github.com/ethereum/hive/simulators/common"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/hive/simulators/common/providers/hive"
	telnet "github.com/reiver/go-telnet"
)

func main() {
	host := hive.New()

	availableClients, err := host.GetClientTypes()
	if err != nil {
		log.Error("could not get client types: ", err.Error())
	}
	log.Info("Got clients", "clients", availableClients)

	logFile, _ := os.LookupEnv("HIVE_SIMLOG")

	for _, client := range availableClients {
		suiteID, err := host.StartTestSuite("iptest", "ip test", logFile)
		if err != nil {
			log.Error("unable to start test suite: ", err.Error())
			os.Exit(1)
		}

		testID, err := host.StartTest(suiteID, "iptest", "iptest")
		if err != nil {
			log.Error("unable to start test: ", err.Error())
			os.Exit(1)
		}

		env := map[string]string{
			"CLIENT": client,
		}
		files := map[string]string{}

		containerID, ip, _, err := host.GetNode(suiteID, testID, env, files)
		if err != nil {
			log.Error("could not get node", "err", err.Error())
			os.Exit(1)
		}

		ourOwnContainerID, err := host.GetSimContainerID(suiteID)
		if err != nil {
			log.Error("could not get sim container IP", "err", err.Error())
			os.Exit(1)
		}
		log.Info("OUR OWN CONTAINER ID", "ID", ourOwnContainerID)

		networkID, err := host.CreateNetwork(suiteID, "network1")
		if err != nil {
			log.Error("could not create network", "err", err.Error())
			os.Exit(1)
		}
		// TODO how to connect own sim container to this network

		// connect client to network
		if err := host.ConnectContainerToNetwork(suiteID, networkID, containerID); err != nil {
			log.Error("could not connect container to network", "err", err.Error())
			os.Exit(1)
		}
		// connect sim to network
		if err := host.ConnectContainerToNetwork(suiteID, networkID, ourOwnContainerID); err != nil {
			log.Error("could not connect container to network", "err", err.Error())
			os.Exit(1)
		}

		// get client ip
		clientIP, err := host.GetContainerNetworkIP(suiteID, networkID, containerID)
		if err != nil {
			log.Error("could not get client network ip addresses", "err", err.Error())
			os.Exit(1)
		}

		conn, err := telnet.DialTo(fmt.Sprintf("%s:30303", clientIP))
		if err != nil {
			log.Error("could not dial", "err", err.Error())
			os.Exit(1)
		}

		if _, err := conn.Write([]byte("ping")); err != nil {
			log.Error("could not write to conn", "err", err.Error())
			os.Exit(1)
		}

		//get our own ip
		simIP, err := host.GetContainerNetworkIP(suiteID, networkID, ourOwnContainerID)
		if err != nil {
			log.Error("could not get client network ip addresses", "err", err.Error())
			os.Exit(1)
		}

		log.Info("got bridge IP: ", "ip", ip)
		log.Info("got network1 ip for client", clientIP)
		log.Info("got network1 ip for sim", simIP)

		host.KillNode(suiteID, testID, containerID)
		host.EndTest(suiteID, testID, &common.TestResult{Pass: true, Details: fmt.Sprint("clientIP: %s, simIP: %s", clientIP, simIP)}, nil)
		host.EndTestSuite(suiteID)
	}
}
