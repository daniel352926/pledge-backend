package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/api/models/ethereum"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
	"pledge-backend/db"
	"pledge-backend/log"
	"strconv"
)

type EthController struct {
}

// Test curl "http://localhost:8080/api/v21/eth/test"
func (c *EthController) Test(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	res.OK(ctx, "ok")
	return
}

// GetBlock curl "http://localhost:8080/api/v21/eth/block/head"
// curl "http://localhost:8080/api/v21/eth/block/8196808"
// curl "http://localhost:8080/api/v21/eth/block/8196808?full=true"
func (c *EthController) GetBlock(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.BlockParams{}

	errCode := validate.NewBlockParams().BlockParams(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	// 获取到 blockNum 请求参数
	blockNum := ctx.Params.ByName("block_num")

	var blockMap map[string]interface{}
	var blockMapPtr *map[string]interface{}
	switch blockNum {
	case "head", "finalized", "safe":
		var err error
		blockMapPtr, err = GetBlockByTag(blockNum)
		if err != nil {
			log.Logger.Error("GetBlockByRpcCall failed:" + err.Error())
			res.Response(ctx, statecode.CommonErrServerErr, nil)
			return
		}
	default:
		// 字符串转 uint（用 strconv.ParseUint）
		num, err := strconv.ParseUint(blockNum, 10, 64)
		if err != nil {
			log.Logger.Error("ParseUint failed:" + err.Error())
			res.Response(ctx, statecode.CommonErrServerErr, nil)
			return
		}
		blockMapPtr, err = c.GetBlockByNum(num)
	}

	if blockMapPtr != nil {
		blockMap = *blockMapPtr
	}

	// 这里将map转换为自定义的block结构体
	block, err := mapToBlock(blockMap)
	if err != nil {
		log.Logger.Error("mapToBlock failed:" + err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}
	// full参数=false时，返回block结构体
	if !req.Full {
		res.OK(ctx, block)
		return
	}

	// 获取交易信息
	// 获取 transactions 字段
	hashes, err := services.NewTx().GetByHashes(block.Transactions)
	if err != nil {
		log.Logger.Error("GetByHashes failed:" + err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	if hashes == nil {
		// 这里循环从链上获取交易信息，大概率会因接口调用过于频繁而报错
		for _, txHash := range block.Transactions {
			tx, err := ethereum.GetTxByHash(txHash)
			if err != nil {
				log.Logger.Error("GetTxByHash failed:" + txHash + ",error:" + err.Error())
				continue
			}
			block.TxList = append(block.TxList, tx)
		}
	} else {
		block.TxList = append(block.TxList, hashes)
	}

	res.OK(ctx, block)
	return
}

func mapToBlock(data map[string]interface{}) (*ethereum.Block, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var block ethereum.Block
	err = json.Unmarshal(jsonData, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (c *EthController) GetBlockByNum(blockNum uint64) (*map[string]interface{}, error) {
	//转成十六进制字符串
	arg := fmt.Sprintf("0x%x", blockNum)

	// 1. 先查数据库
	block, err := services.NewBlockService().GetByNum(arg)
	if err != nil {
		return nil, err
	}
	if block != nil {
		blockJson, err := json.Marshal(block)
		if err != nil {
			return nil, err
		}
		var blockMap map[string]interface{}
		err = json.Unmarshal(blockJson, &blockMap)
		if err != nil {
			return nil, err
		}
		return &blockMap, err
	}

	// 2. 数据库不存在，再查RPC
	blockMap, err := ethereum.GetBlockByRpcCall(arg)
	if err != nil {
		return nil, err
	}

	// 3. 存入数据库
	// 提取其他字段并转换为字符串
	model, err := MapToBlockModel(*blockMap)

	_ = services.NewBlockService().Save(model)

	return blockMap, nil
}

func MapToBlockModel(data map[string]interface{}) (*models.Block, error) {
	// 处理 transactions 和 uncles 两个数组
	transactions, err := json.Marshal(data["transactions"])
	if err != nil {
		return nil, err
	}
	uncles, err := json.Marshal(data["uncles"])
	if err != nil {
		return nil, err
	}

	// 安全提取字段（如果字段缺失不会panic）
	getString := func(key string) string {
		if v, ok := data[key]; ok && v != nil {
			return v.(string)
		}
		return ""
	}

	blockModel := &models.Block{
		Number:           getString("number"),
		Hash:             getString("hash"),
		ParentHash:       getString("parentHash"),
		Nonce:            getString("nonce"),
		Sha3Uncles:       getString("sha3Uncles"),
		LogsBloom:        getString("logsBloom"),
		TransactionsRoot: getString("transactionsRoot"),
		StateRoot:        getString("stateRoot"),
		ReceiptsRoot:     getString("receiptsRoot"),
		Miner:            getString("miner"),
		Difficulty:       getString("difficulty"),
		TotalDifficulty:  getString("totalDifficulty"),
		ExtraData:        getString("extraData"),
		Size:             getString("size"),
		GasLimit:         getString("gasLimit"),
		GasUsed:          getString("gasUsed"),
		Timestamp:        getString("timestamp"),
		Transactions:     string(transactions),
		Uncles:           string(uncles),
	}

	return blockModel, nil
}

func GetBlockByTag(tag string) (*map[string]interface{}, error) {
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

// GetTxByHash curl "http://localhost:8080/api/v21/eth/tx/0xb48abc9e971287a7adea4d3a3902551f7f589f7f634935438babd1abe5c93e14"
func (c *EthController) GetTxByHash(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	txHash := ctx.Params.ByName("tx_hash")
	if txHash == "" {
		res.Response(ctx, statecode.ParameterEmptyErr, nil)
		return
	}

	// 从数据库里获取交易信息
	tx, err := services.NewTx().GetByHash(txHash)
	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	if tx != nil {
		res.OK(ctx, tx)
		return
	}

	// 从链上获取交易信息
	transaction, err := ethereum.GetTxByHash(txHash)
	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
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
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	//fmt.Println(fmt.Sprintf("0x%x", transaction.Type()))
	//fmt.Println("!!!!0x" + transaction.GasPrice().Text(16))
	err = services.NewTx().Save(&models.Tx{
		Hash:                 transaction.Hash().Hex(),
		Type:                 transaction.Type(),
		Nonce:                transaction.Nonce(),
		GasPrice:             "0x" + transaction.GasPrice().Text(16),
		MaxPriorityFeePerGas: "0x" + transaction.GasTipCap().Text(16),
		MaxFeePerGas:         "0x" + transaction.GasFeeCap().Text(16),
		Gas:                  fmt.Sprintf("0x%x", transaction.Gas()),
		Value:                "0x" + transaction.Value().Text(16),
		Input:                string(transaction.Data()),
		V:                    "0x" + v.Text(16),
		R:                    "0x" + r.Text(16),
		S:                    "0x" + s.Text(16),
		To:                   transaction.To().Hex(),
		ChainId:              "0x" + transaction.ChainId().Text(16),
		AccessList:           string(accessList),
	})
	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	res.OK(ctx, transaction)
}

// GetReceiptByHash curl "http://localhost:8080/api/v21/etc/tx_receipt/0xb48abc9e971287a7adea4d3a3902551f7f589f7f634935438babd1abe5c93e14"
func (c *EthController) GetReceiptByHash(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	txHash := ctx.Params.ByName("tx_hash")

	if txHash == "" {
		res.Response(ctx, statecode.ParameterEmptyErr, nil)
		return
	}

	// 1. 从数据库里获取
	receiptByMysql, err := services.NewReceipt().GetByHash(txHash)
	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}
	if receiptByMysql != nil {
		res.OK(ctx, receiptByMysql)
		return
	}

	// 2. 从链上获取
	receipt, err := ethereum.GetReceiptByHash(txHash)
	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	logs, err := json.Marshal(receipt.Logs)
	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	// 写入数据库
	err = services.NewReceipt().Save(&models.Receipt{
		Hash:              receipt.TxHash.Hex(),
		Type:              fmt.Sprintf("0x%x", receipt.Type),
		Status:            fmt.Sprintf("0x%x", receipt.Status),
		Root:              "0x" + hex.EncodeToString(receipt.PostState),
		CumulativeGasUsed: fmt.Sprintf("0x%x", receipt.CumulativeGasUsed),
		LogsBloom:         fmt.Sprintf("0x%x", receipt.Bloom[:]),
		ContractAddress:   receipt.ContractAddress.Hex(),
		GasUsed:           fmt.Sprintf("0x%x", fmt.Sprintf("0x%x", receipt.Status)),
		BlockHash:         receipt.BlockHash.Hex(),
		BlockNumber:       receipt.BlockNumber.String(),
		TransactionIndex:  strconv.Itoa(int(receipt.TransactionIndex)),
		Logs:              string(logs),
	})
	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, receipt)
		return
	}

	res.OK(ctx, receipt)
}
