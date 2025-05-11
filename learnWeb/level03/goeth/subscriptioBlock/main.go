package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {

	//client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	//client, err := ethclient.Dial("wss://ropsten.infura.io/ws/v3/YOUR_PROJECT_ID")
	//client, err := ethclient.Dial("wss://sepolia.infura.io/v3/e7099d8da3594659a6ffc36d3e68d24b")
	//client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/e7099d8da3594659a6ffc36d3e68d24b")
	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/e7099d8da3594659a6ffc36d3e68d24b")
	if err != nil {
		log.Fatal(err)
	}
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case headers := <-headers:
			fmt.Println("headers:", headers.Hash().Hex())
			block, err := client.BlockByHash(context.Background(), headers.Hash())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64())   // 3477413
			fmt.Println(block.Time())              // 1529525947
			fmt.Println(block.Nonce())             // 130524141876765836
			fmt.Println(len(block.Transactions())) // 7
		}
	}
}
