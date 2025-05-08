package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
)

// 获取区块信息
func main() {
	// 获取区块信息
	fmt.Println("获取区块信息-------------------------")
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/e7099d8da3594659a6ffc36d3e68d24b")
	if err != nil {
		log.Fatal(err)
	}
	//查询余额
	account := common.HexToAddress("0x2bc264624a60BF7De7EA0068b71B63BaA51C27c6")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)
	blockNumber1 := big.NewInt(5532993)
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balanceAt) // 25729324269165216042
	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue) // 25.729324269165216041
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance) // 25729324269165216042

}
