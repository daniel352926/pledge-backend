package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"pledge-backend/api/models"
	"pledge-backend/config"
	"pledge-backend/log"
)

type Block struct {
	Block  models.Block
	TxList []models.Tx `json:"txList"`
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

func GetBlockByNum(num *int64) (*types.Block, error) {
	client, err := ethclient.Dial(config.Config.Eth.RawUrl)
	if err != nil {
		log.Logger.Fatal("获取Eth客户端失败")
		return nil, err
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(*num))
	if err != nil {
		return nil, err
	}

	return block, nil
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
