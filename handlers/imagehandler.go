package handlers

import (
	"fmt"
	"gimg/pkg"
	"github.com/gin-gonic/gin"
)

// GetHandler Get image
func GetHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

// UploadHandler Upload image
func UploadHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("filename")

		if err != nil {
			pkg.Fail(c, "Parse argument error")
			fmt.Println(err)
			return
		}

		fObj, err := file.Open()
		if err != nil {
			pkg.Fail(c, "Open temp file error")
			return
		}
		defer fObj.Close()

		md5, _ := pkg.CalcMd5(fObj)
		fmt.Println(md5)
		err = ctx.SaveFile(md5, fObj)
		if err != nil {
			pkg.Fail(c, "Save file error")
			return
		}
	}
}
