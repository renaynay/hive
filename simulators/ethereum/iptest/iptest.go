package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/ethereum/hive/simulators/common"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/hive/simulators/common/providers/hive"
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

		type Payload struct {
			Jsonrpc string        `json:"jsonrpc"`
			Method  string        `json:"method"`
			Params  []interface{} `json:"params"`
			ID      int           `json:"id"`
		}

		data := Payload{
			Jsonrpc: "2.0",
			Method:  "web3_clientVersion",
			ID:      67,
		}
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			log.Error("could not marshal payload", "err", err.Error())
			os.Exit(1)
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s:8545"), body)
		if err != nil {
			log.Error("could not create new request", "err", err.Error())
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("could not send request to geth node", "err", err.Error())
			os.Exit(1)
		}

		spew.Dump(resp)

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
