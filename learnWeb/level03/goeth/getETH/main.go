package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
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
	fmt.Println("eth转账：-------------------------")
	//eth转账
	privateKeys, err := crypto.HexToECDSA("9923d7312ae8bae86f5323774271e76835ae35013b735f8ca709a39d798e3f1c")
	if err != nil {
		log.Fatal(err)
	}

	publicKeys := privateKeys.Public()
	publicKeyECDSAs, ok := publicKeys.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSAs)
	fmt.Println("fromAddress:", fromAddress.Hex())
	if "0x2bc264624a60BF7De7EA0068b71B63BaA51C27c6" == fromAddress.Hex() {
		fmt.Println("fromAddress:", "0x2bc264624a60BF7De7EA0068b71B63BaA51C27c6")
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(100000000000) // in wei (1 eth)
	gasLimit := uint64(21000)         // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0xeF421Da63310b49EEd742CA928f0e1156455c757")
	var data []byte
	tx1 := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainIDs, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx1, types.NewEIP155Signer(chainIDs), privateKeys)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx1 sent: %s", signedTx.Hash().Hex())

}
