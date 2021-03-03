package main

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/hive/hivesim"
)

var hiveParams = map[string]string{
	"HIVE_NETWORK_ID":     "19763",
	"HIVE_CHAIN_ID":       "19763",
	"HIVE_FORK_HOMESTEAD": "0",
	"HIVE_FORK_TANGERINE": "0",
	"HIVE_FORK_SPURIOUS":  "0",
	"HIVE_FORK_BYZANTIUM": "0",
	"HIVE_LOGLEVEL":       "5",
}

func main() {
	suite := hivesim.Suite{
		Name: "block production",
		Description: "This suite tests a client's ability to accurately produce and " +
			"propagate blocks throughout an ethereum network",
	}
	suite.Add(hivesim.TestSpec{
		Name: "client launch",
		Description: `This test launches the client and runs the test tool.
Results from the test tool are reported as individual sub-tests.`,
		Run: initNodes,
	})
	hivesim.MustRunSuite(hivesim.New(), suite)
}

func initNodes(t *hivesim.T) {
	// get clients and make sure there are at least 2 or 3?
	clients, err := t.Sim.ClientTypes()
	if err != nil {
		t.Fatal("could not get client types: %v", err)
	}

	// TODO fix
	_, err = t.Sim.StartTest(t.SuiteID, "block production test", "This suite tests a client's ability to accurately produce and "+
		"propagate blocks throughout an ethereum network\"")
	if err != nil {
		t.Fatalf("could not start test: %v", err)
	}

	for i, client := range clients {
		// make the last client the mining node
		if i == len(clients)-1 {
			mineParams := hiveParams
			mineParams["HIVE_MINER"] = "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
			t.RunClient(client, hivesim.ClientTestSpec{
				Name:        client,
				Description: "mining node",
				Parameters:  mineParams,
				Files: map[string]string{
					"genesis.json": "./init/genesis.json",
				},
				Run: runMiner,
			})
		}
		//t.RunClient(client, hivesim.ClientTestSpec{
		//	Name:        client,
		//	Description: "non-mining node",
		//	Parameters:  hiveParams,
		//	// todo files?
		//	Run: runNodes,
		//})
	}

}

type result struct {
	Result string `json:"result"`
}

// TODO what to do here ? send a bunch of txs to the node?
func runMiner(t *hivesim.T, c *hivesim.Client) {
	//// read genesis file
	//gen := readGenesisFile(t)
	//gblock := gen.ToBlock(nil)
	////// Load chain.rlp.
	////blocks := toBlocks(t, gblock)
	//
	//sendTxsToMiner(t, c, gblock)
	//// check if it's successful
	//
	//
	//time.Sleep(2 * time.Second)
	//params := map[string]string{
	//	"from":     "0x71562b71999873db5b286df957af199ec94617f7",
	//	"to":       "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
	//	"gas":      "0x76c0", // 30400
	//	"gasPrice": "0x1388", // 5000
	//	"value":    "0x64",   // 100
	//}
	//
	//res := new(result)
	//if err := c.RPC().Call(res, "eth_sendTransaction", params); err != nil {
	//	t.Fatalf("could not complete rpc call: %v", err)
	//}

	time.Sleep(5000 * time.Second)
	t.Sim.EndTest(t.SuiteID, t.TestID, hivesim.TestResult{})
}

//func sendTxsToMiner(t *hivesim.T, c *hivesim.Client, chain *Chain) {
//	// get mining node rpc endpoint
//	miner := c.RPC()
//
//	tx := generateTx()
//
//	//// send txs starting from block to mining node
//	//for i := 1; i <= 5; i++ {
//	//	txs := chain.blocks[i].Transactions()
//	//	for _, tx := range txs {
//	//		// get tx as message to get sender
//	//		txMsg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
//	//		if err != nil {
//	//			t.Fatalf("could not get tx as message: %v", txMsg)
//	//		}
//	//
//	//		var hash [32]byte
//	//		callParams := map[string]string{
//	//			"from":     txMsg.From().String(),
//	//			"to":       tx.To().String(),
//	//			"gas":      hexutil.EncodeUint64(tx.Gas()),
//	//			"gasPrice": hexutil.EncodeBig(tx.GasPrice()), // 10000000000000
//	//			"value":    hexutil.EncodeBig(tx.Value()),    // 2441406250
//	//			"data":     string(tx.Data()),
//	//		}
//	//		miner.Call(hash, "eth_sendTransaction", callParams)
//	//	}
//	//
//	//	_ = tESTING(t, c) // TODO REMOVE
//	//
//	//	//block := waitForBlockPropagation(t)
//	//	//if block.Number() != chain.blocks[i].Number() {
//	//	//	t.Fatalf("block number mismatch: expected %d, got %d", i, block.Number())
//	//	//}
//	//}
//}

func generateTx() *types.Transaction {
	return nil
}

func tESTING(t *hivesim.T, c *hivesim.Client) bool {
	time.Sleep(time.Millisecond * 1000)
	headers := make(chan *types.Header)
	sub, err := c.RPC().EthSubscribe(context.Background(), headers)
	if err != nil {
		t.Fatalf("could not subscribe to new headers: %v", err)
	}

	select {
	case err := <-sub.Err():
		t.Fatalf("got subscription error: %v", err)
	case header := <-headers:
		t.Log("GOT BLOCK HEADER \nHEADER NUMBER: %d\n CONTENTS: %v", header.Number, header.Hash().Hex())
	}

	return true
}

//
//func waitForBlockPropagation(t *hivesim.T) *types.Block {
//
//	// todo query other nodes over and over til you get a new block announcement or something like that
//}

//func readGenesisFile(t *hivesim.T) core.Genesis {
//	rawGen, err := ioutil.ReadFile("genesis.json")
//	if err != nil {
//		t.Fatal("could not read genesis file: %v", err)
//	}
//	var gen core.Genesis
//	if err := json.Unmarshal(rawGen, &gen); err != nil {
//		t.Fatalf("could not unmarshal genesis file: %v", err)
//	}
//	return gen
//}

//
//func toBlocks(t *hivesim.T, gblock *types.Block) []*types.Block {
//	fh, err := os.Open("/chain.rlp")
//	if err != nil {
//		t.Fatal("could not read chain file: %v", err)
//	}
//	defer fh.Close()
//	var reader io.Reader = fh
//	stream := rlp.NewStream(reader, 0)
//	var blocks = make([]*types.Block, 1)
//	blocks[0] = gblock
//	for i := 0; ; i++ {
//		var b types.Block
//		if err := stream.Decode(&b); err == io.EOF {
//			break
//		} else if err != nil {
//			t.Fatal(err)
//		}
//		if b.NumberU64() != uint64(i+1) {
//			t.Fatalf("block at index %d has wrong number %d", i, b.NumberU64())
//		}
//		blocks = append(blocks, &b)
//	}
//	return blocks
//}

func runNodes(t *hivesim.T, c *hivesim.Client) {
	// TODO
	// time.Wait
	//
	// query nodes to see if they got new blocks from mining node
	// maybe send some txs to regular nodes, see if txs get propagated to miner, then query miner
}
