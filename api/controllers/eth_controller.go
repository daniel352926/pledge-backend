package controllers

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-gonic/gin"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/ethereum"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
	"pledge-backend/log"
	"strconv"
)

type EthController struct {
}

var ethService = services.NewEthService()

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

	var blockNumber int64
	if blockNum == "head" {
		blockNumber = int64(rpc.LatestBlockNumber)
	} else if blockNum == "finalized" {
		blockNumber = int64(rpc.FinalizedBlockNumber)
	} else if blockNum == "safe" {
		blockNumber = int64(rpc.SafeBlockNumber)
	} else {
		// 字符串blockNum转 bigInt
		num, err := strconv.ParseInt(blockNum, 10, 64)
		if err != nil {
			res.Response(ctx, statecode.CommonErrServerErr, nil)
			return
		}
		blockNumber = num
	}
	block, err := ethService.GetBlockByNumber(blockNumber)
	if err != nil {
		log.Logger.Error("GetBlockByNumber failed:" + err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	resBlock := ethereum.Block{
		Block:  *block,
		TxList: nil,
	}

	// full参数=false时，返回block结构体
	if !req.Full {
		res.OK(ctx, resBlock)
		return
	}

	// 获取交易信息
	// 获取 transactions 字段
	txList, err := services.NewTx().GetByBlockNum(block.Number)
	if err != nil {
		log.Logger.Error("GetByHashes failed:" + err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	if txList == nil {
		// 这里循环从链上获取交易信息，大概率会因接口调用过于频繁而报错
		var transactions types.Transactions
		// 把block.Transactions字符串json.Unmarsha解析成types.Transactions
		err := json.Unmarshal([]byte(block.Transactions), &transactions)
		if err != nil {
			log.Logger.Error("Transactions json.Unmarshal failed:" + err.Error())
			res.Response(ctx, statecode.CommonErrServerErr, nil)
			return
		}

		for _, transaction := range transactions {
			tx, err := ethService.GetTxByHash(transaction.Hash().Hex())
			if err != nil {
				log.Logger.Error("GetTxByHash failed:" + transaction.Hash().Hex() + ",error:" + err.Error())
				continue
			}
			resBlock.TxList = append(resBlock.TxList, *tx)
		}
	} else {
		resBlock.TxList = append(resBlock.TxList, *txList...)
	}

	res.OK(ctx, resBlock)
}

//func (c *EthController) GetBlockByNum(blockNum uint64) (*map[string]interface{}, error) {
//	//转成十六进制字符串
//	arg := fmt.Sprintf("0x%x", blockNum)
//
//	blockMap, err := ethService.GetBlockByNumber(arg)
//	if err != nil {
//		return nil, err
//	}
//
//	return blockMap, nil
//}

// GetTxByHash curl "http://localhost:8080/api/v21/eth/tx/0xb48abc9e971287a7adea4d3a3902551f7f589f7f634935438babd1abe5c93e14"
func (c *EthController) GetTxByHash(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	txHash := ctx.Params.ByName("tx_hash")
	if txHash == "" {
		res.Response(ctx, statecode.ParameterEmptyErr, nil)
		return
	}

	transaction, err := ethService.GetTxByHash(txHash)
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

	receipt, err := ethService.GetReceiptByHash(txHash)

	if err != nil {
		res.Response(ctx, statecode.CommonErrServerErr, receipt)
		return
	}

	res.OK(ctx, receipt)
}
