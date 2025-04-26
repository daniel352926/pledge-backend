package validate

import (
	"github.com/gin-gonic/gin"
	"io"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/request"
)

type BlockParams struct {
}

func NewBlockParams() *BlockParams {
	return &BlockParams{}
}

func (v *BlockParams) BlockParams(ctx *gin.Context, req *request.BlockParams) int {
	err := ctx.ShouldBind(req)
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		return statecode.CommonErrServerErr
	}
	return statecode.CommonSuccess
}
