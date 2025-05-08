package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
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
	blockNumber := big.NewInt(8268485)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())        // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
		fmt.Println(tx.Value().String())    // 100000000000000000
		fmt.Println(tx.Gas())               // 21000
		fmt.Println(tx.GasPrice().Uint64()) // 100000000000
		fmt.Println(tx.Nonce())             // 245132
		fmt.Println(tx.Data())              // []
		fmt.Println(tx.To().Hex())          // 0x8F9aFd209339088Ced7Bc0f57Fe08566ADda3587

		if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
			fmt.Println("sender", sender.Hex()) // 0x2CdA41645F2dBffB852a605E92B185501801FC28
		} else {
			log.Fatal(err)
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(receipt.Status) // 1
		fmt.Println(receipt.Logs)   // []
		break
	}

	fmt.Println("获取区块交易信息-------------------------")

	blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tx.Hash().Hex()) // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
		break
	}

	txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(isPending)
	fmt.Println(tx.Hash().Hex())

	fmt.Println("获取区块收据信息-------------------------")
	//查询收据
	//receiptByHash, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(blockHash, false))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//receiptsByNum, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(receiptByHash[0] == receiptsByNum[0]) // true
	//
	//for _, receipt := range receiptByHash {
	//	fmt.Println(receipt.Status)                // 1
	//	fmt.Println(receipt.Logs)                  // []
	//	fmt.Println(receipt.TxHash.Hex())          // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
	//	fmt.Println(receipt.TransactionIndex)      // 0
	//	fmt.Println(receipt.ContractAddress.Hex()) // 0x0000000000000000000000000000000000000000
	//	break
	//}
	//
	////txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")
	//receipt, err := client.TransactionReceipt(context.Background(), txHash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(receipt.Status)                // 1
	//fmt.Println(receipt.Logs)                  // []
	//fmt.Println(receipt.TxHash.Hex())          // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
	//fmt.Println(receipt.TransactionIndex)      // 0
	//fmt.Println(receipt.ContractAddress.Hex()) // 0x0000000000000000000000000000000000000000

	//创建钱包
	fmt.Println("创建钱包：-------------------------")
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	//这就是签署交易的私钥
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 去掉'0x'
	//返回公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("from pubKey:", hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04'
	//根据公钥返回公共地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)
	//公共地址其实就是公钥的 Keccak-256 哈希，然后我们取最后 40 个字符（20 个字节）并用“0x”作为前缀
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("full:", hexutil.Encode(hash.Sum(nil)[:]))
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 原长32位，截去12位，保留后20位

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

	value := big.NewInt(0)    // in wei (1 eth)
	gasLimit := uint64(21000) // in units
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

	signedTx, err := types.SignTx(tx1, types.NewEIP155Signer(chainIDs), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx1 sent: %s", signedTx.Hash().Hex())

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
