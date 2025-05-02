package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gomodule/redigo/redis"
	"pledge-backend/api/ethereum"
	"pledge-backend/api/models"
	"pledge-backend/db"
	"pledge-backend/log"
	"strconv"
	"time"
)

type EthService struct {
}

var blockService = NewBlockService()
var txService = NewTx()
var receiptService = NewReceipt()

func NewEthService() *EthService {
	return &EthService{}
}

// GetBlockByNum 根据块号获取块信息
//func (e *EthService) GetBlockByNum(blockNum string) (*map[string]interface{}, error) {
//	// 1. 先查数据库
//	block, err := blockService.GetByNum(blockNum)
//	if err != nil {
//		return nil, err
//	}
//	if block != nil {
//		blockJson, err := json.Marshal(block)
//		if err != nil {
//			return nil, err
//		}
//		var blockMap map[string]interface{}
//		err = json.Unmarshal(blockJson, &blockMap)
//		if err != nil {
//			return nil, err
//		}
//		return &blockMap, err
//	}
//
//	// 2. 数据库不存在，再查RPC
//	blockMap, err := ethereum.GetBlockByRpcCall(blockNum)
//	if err != nil {
//		return nil, err
//	}
//
//	// 3. 存入数据库
//	// 提取其他字段并转换为字符串
//	model, err := MapToBlockModel(*blockMap)
//
//	_ = block.Save(model)
//
//	return blockMap, nil
//}

func (e *EthService) GetTxByHash(txHash string) (*models.Tx, error) {
	// 从数据库里获取交易信息
	tx, err := txService.GetByHash(txHash)
	if err != nil {
		return nil, err
	}

	if tx != nil {
		return tx, nil
	}

	// 从链上获取交易信息
	transaction, err := ethereum.GetTxByHash(txHash)
	if err != nil {
		return tx, nil
	}

	// 从链上获取回执信息 (主要获取block上的高度和交易的t_index记录到交易表中)
	receipt, err := ethereum.GetReceiptByHash(txHash)
	if err != nil {
		return nil, err
	}

	// 存入数据库中
	v, r, s := transaction.RawSignatureValues()

	// 获取 MaxPriorityFeePerGas（EIP-1559 类型交易才有）
	//if transaction.Type() == types.DynamicFeeTxType {
	//	fmt.Println("MaxPriorityFeePerGas:", transaction.GasTipCap()) // *big.Int
	//	fmt.Println("MaxFeePerGas:", transaction.GasFeeCap())         // *big.Int
	//} else {
	//	fmt.Println("GasPrice:", transaction.GasPrice()) // *big.Int（用于非EIP-1559）
	//}

	accessList, err := json.Marshal(transaction.AccessList())
	if err != nil {
		return nil, err
	}

	model := models.Tx{
		Hash:                 transaction.Hash().Hex(),
		BlockNum:             receipt.BlockNumber.Uint64(),
		Index:                receipt.TransactionIndex,
		Type:                 transaction.Type(),
		Nonce:                transaction.Nonce(),
		GasPrice:             "0x" + transaction.GasPrice().Text(16),
		MaxPriorityFeePerGas: "0x" + transaction.GasTipCap().Text(16),
		MaxFeePerGas:         "0x" + transaction.GasFeeCap().Text(16),
		Gas:                  fmt.Sprintf("0x%x", transaction.Gas()),
		Value:                "0x" + transaction.Value().Text(16),
		Input:                common.BytesToHash(transaction.Data()).Hex(),
		V:                    "0x" + v.Text(16),
		R:                    "0x" + r.Text(16),
		S:                    "0x" + s.Text(16),
		To:                   transaction.To().Hex(),
		ChainId:              "0x" + transaction.ChainId().Text(16),
		AccessList:           string(accessList),
	}

	//fmt.Println(fmt.Sprintf("0x%x", transaction.Type()))
	//fmt.Println("!!!!0x" + transaction.GasPrice().Text(16))
	err = txService.Save(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (e *EthService) GetReceiptByHash(txHash string) (*models.Receipt, error) {
	// 1. 从数据库里获取
	receiptByMysql, err := receiptService.GetByHash(txHash)
	if err != nil {
		return nil, err
	}
	if receiptByMysql != nil {
		return receiptByMysql, nil
	}

	// 2. 从链上获取
	receipt, err := ethereum.GetReceiptByHash(txHash)
	if err != nil {
		return nil, err
	}

	logs, err := json.Marshal(receipt.Logs)
	if err != nil {
		return nil, err
	}

	model := models.Receipt{
		Hash:              receipt.TxHash.Hex(),
		Type:              receipt.Type,
		Status:            receipt.Status,
		Root:              common.Bytes2Hex(receipt.PostState),
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		LogsBloom:         fmt.Sprintf("0x%x", receipt.Bloom[:]),
		ContractAddress:   receipt.ContractAddress.Hex(),
		GasUsed:           fmt.Sprintf("0x%x", fmt.Sprintf("0x%x", receipt.Status)),
		BlockHash:         receipt.BlockHash.Hex(),
		BlockNumber:       receipt.BlockNumber.String(),
		TransactionIndex:  strconv.Itoa(int(receipt.TransactionIndex)),
		Logs:              string(logs),
	}

	// 写入数据库
	err = receiptService.Save(&model)

	return &model, err
}

func (e *EthService) GetBlockByTag(tag string) (*map[string]interface{}, error) {
	var arg string
	//把tag赋值给 arg，tag == "head"，arg = "latest"
	if tag == "head" {
		arg = "latest"
	} else {
		arg = tag
	}

	//1. 先从Redis里获取
	redisData, err := db.RedisGet(tag)
	if err != nil {
		//报错不退出，直接打log
		log.Logger.Error("redis get error:" + err.Error())
		return nil, err
	}

	var blockMap map[string]interface{}

	if redisData != nil {
		// 返回有数据则解析成map并返回
		err := json.Unmarshal(redisData, &blockMap)
		if err != nil {
			log.Logger.Error("Unmarshal json failed:" + err.Error())
		}
		return &blockMap, err
	}

	// redis里没有数据则从链上获取
	_blockMap, err := ethereum.GetBlockByRpcCall(arg)
	if err != nil {
		return nil, err
	}

	// 取到设置到Redis
	err = db.RedisSet(tag, _blockMap, 120)
	if err != nil {
		log.Logger.Error("Unmarshal json failed:" + err.Error())
	}
	return _blockMap, nil
}

func (e *EthService) GetBlockByNumber(blockNumber int64) (*models.Block, error) {
	redisKey := "RedisBlock" + strconv.FormatInt(blockNumber, 10)
	//1. 如果是特殊区块，先从Redis里获取
	if blockNumber < 0 {
		redisData, err := db.RedisGet(redisKey)
		if err != nil && !errors.Is(err, redis.ErrNil) {
			//报错不退出，直接打log
			log.Logger.Error("redis get error:" + err.Error())
			return nil, err
		}

		var block *models.Block

		if redisData != nil {
			// 返回有数据则解析成map并返回
			err := json.Unmarshal(redisData, &block)
			if err != nil {
				log.Logger.Error("Unmarshal json failed:" + err.Error())
			}
			return block, err
		}
	}

	//2. redis里没取到，从链上获取
	block, err := ethereum.GetBlockByNum(&blockNumber)
	if err != nil {
		log.Logger.Error("ParseUint failed:" + err.Error())
		return nil, err
	}

	// 构造返回体
	txList, err := json.Marshal(block.Transactions())
	if err != nil {
		log.Logger.Error("Transactions json.Marshal failed:" + err.Error())
		return nil, err
	}

	uncles, err := json.Marshal(block.Uncles())
	if err != nil {
		log.Logger.Error("uncles json.Marshal failed:" + err.Error())
		return nil, err
	}

	model := models.Block{
		Number:           block.Number().Uint64(),
		Hash:             block.Hash().Hex(),
		ParentHash:       block.ParentHash().Hex(),
		Nonce:            block.Nonce(),
		Sha3Uncles:       block.UncleHash().Hex(),
		LogsBloom:        common.Bytes2Hex(block.Bloom().Bytes()),
		TransactionsRoot: block.Header().TxHash.Hex(),
		StateRoot:        block.Header().Root.Hex(),
		ReceiptsRoot:     block.Header().ReceiptHash.Hex(),
		Miner:            block.Coinbase().Hex(),
		Difficulty:       block.Difficulty().String(),
		ExtraData:        common.Bytes2Hex(block.Header().Extra),
		Size:             strconv.FormatUint(block.Size(), 10),
		GasLimit:         strconv.FormatUint(block.GasLimit(), 10),
		GasUsed:          strconv.FormatUint(block.GasUsed(), 10),
		Timestamp:        block.Time(),
		Transactions:     string(txList),
		Uncles:           string(uncles),
		CreatedAt:        time.Time{},
	}

	if blockNumber > 0 {
		// 如果正常按照区块高度获取的数据，则保存到数据库
		err = blockService.Save(&model)
		if err != nil {
			log.Logger.Error("mysql save failed:" + err.Error())
			return &model, err
		}
	} else {
		// 特殊区块缓存到redis
		err = db.RedisSet(redisKey, block, 120)
		if err != nil {
			log.Logger.Error("redis cache failed:" + err.Error())
			return &model, err
		}
	}

	return &model, nil
}

//func MapToBlockModel(data map[string]interface{}) (*models.Block, error) {
//	// 处理 transactions 和 uncles 两个数组
//	transactions, err := json.Marshal(data["transactions"])
//	if err != nil {
//		return nil, err
//	}
//	uncles, err := json.Marshal(data["uncles"])
//	if err != nil {
//		return nil, err
//	}
//
//	// 安全提取字段（如果字段缺失不会panic）
//	getString := func(key string) string {
//		if v, ok := data[key]; ok && v != nil {
//			return v.(string)
//		}
//		return ""
//	}
//
//	blockModel := &models.Block{
//		Number:           getString("number"),
//		Hash:             getString("hash"),
//		ParentHash:       getString("parentHash"),
//		Nonce:            getString("nonce"),
//		Sha3Uncles:       getString("sha3Uncles"),
//		LogsBloom:        getString("logsBloom"),
//		TransactionsRoot: getString("transactionsRoot"),
//		StateRoot:        getString("stateRoot"),
//		ReceiptsRoot:     getString("receiptsRoot"),
//		Miner:            getString("miner"),
//		Difficulty:       getString("difficulty"),
//		TotalDifficulty:  getString("totalDifficulty"),
//		ExtraData:        getString("extraData"),
//		Size:             getString("size"),
//		GasLimit:         getString("gasLimit"),
//		GasUsed:          getString("gasUsed"),
//		Timestamp:        getString("timestamp"),
//		Transactions:     string(transactions),
//		Uncles:           string(uncles),
//	}
//
//	return blockModel, nil
//}
