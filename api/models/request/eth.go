package request

type BlockParams struct {
	Full bool `form:"full,default=false" binding:"omitempty"`
}
