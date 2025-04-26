package controllers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/response"
	"pledge-backend/config"
	"pledge-backend/contract/bindings"
	"pledge-backend/log"
)

type PlgrController struct {
}

type Fee struct {
	BorrowFee  uint64
	LendFee    uint64
	PoolLength uint64
}

// GetFee curl "http://localhost:8080/api/v21/plgr/usdt/get_fee"
func (*PlgrController) GetFee(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	ethereumConn, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error(err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}
	pledgePoolToken, err := bindings.NewPledgePoolToken(common.HexToAddress(config.Config.TestNet.PledgePoolToken), ethereumConn)

	if nil != err {
		log.Logger.Error(err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	borrowFee, err := pledgePoolToken.PledgePoolTokenCaller.BorrowFee(nil)
	if nil != err {
		log.Logger.Error(err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	lendFee, err := pledgePoolToken.PledgePoolTokenCaller.LendFee(nil)
	if nil != err {
		log.Logger.Error(err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	pooLength, err := pledgePoolToken.PledgePoolTokenCaller.PoolLength(nil)
	if nil != err {
		log.Logger.Error(err.Error())
		res.Response(ctx, statecode.CommonErrServerErr, nil)
		return
	}

	res.OK(ctx, Fee{
		BorrowFee:  borrowFee.Uint64(),
		LendFee:    lendFee.Uint64(),
		PoolLength: pooLength.Uint64(),
	})
}
