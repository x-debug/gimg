package pkg

import "github.com/gin-gonic/gin"

func Fail(c *gin.Context, msg string)  {
	c.JSON(500, ErrResp{Msg: msg})
}

func Success()  {

}
