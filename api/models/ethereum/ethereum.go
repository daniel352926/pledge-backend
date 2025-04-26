package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"pledge-backend/config"
	"pledge-backend/log"
)

type Block struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []string      `json:"transactions"`
	TxList           []interface{} `json:"txList"`
	Uncles           []string      `json:"uncles"`
}

func GetClient() *ethclient.Client {
	ethConfig := config.Config.Eth
	client, err := ethclient.Dial(ethConfig.RawUrl)
	if err != nil {
		log.Logger.Fatal("获取Eth客户端失败")
	}

	return client
}

func GetRpcClient() *rpc.Client {
	ethConfig := config.Config.Eth
	fmt.Println(ethConfig)
	client, err := rpc.Dial(ethConfig.RawUrl)
	if err != nil {
		log.Logger.Fatal("获取Eth rpc客户端失败")
	}

	return client
}

// GetBlockByRpcCall 通过原始JSON RPC调用方式获取区块信息
func GetBlockByRpcCall(arg string) (*map[string]interface{}, error) {
	client, err := rpc.Dial(config.Config.Eth.RawUrl)
	if err != nil {
		return nil, err
	}
	//defer client.Close()
	var blockMap map[string]interface{}
	err = client.CallContext(context.Background(), &blockMap, "eth_getBlockByNumber", arg, false)

	return &blockMap, err
}

// GetTxByHash 根据hash获取交易信息
func GetTxByHash(txHash string) (*types.Transaction, error) {
	client, err := ethclient.Dial(config.Config.Eth.RawUrl)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	transaction, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// GetReceiptByHash 根据hash获取交易回执信息
func GetReceiptByHash(txHash string) (*types.Receipt, error) {
	client, err := ethclient.Dial(config.Config.Eth.RawUrl)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}

	return receipt, nil
}
