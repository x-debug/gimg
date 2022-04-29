package handlers

import (
	"gimg/pkg"
	"github.com/gin-gonic/gin"
	"io"
)

// GetHandler Get image
func GetHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		rFile, closef, err := ctx.ReadFile(hash)
		defer closef()
		
		if err != nil {
			pkg.Fail(c, "Read file error")
			return
		}

		_, err = io.Copy(c.Writer, rFile)
		if err != nil {
			pkg.Fail(c, "Copy file error")
			return
		}
	}
}

// UploadHandler Upload image
func UploadHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("filename")

		if err != nil {
			pkg.Fail(c, "Parse argument error")
			return
		}

		fObj, err := file.Open()
		if err != nil {
			pkg.Fail(c, "Open temp file error")
			return
		}
		defer fObj.Close()

		md5, _ := pkg.CalcMd5(fObj)
		err = ctx.SaveFile(md5, fObj)
		if err != nil {
			pkg.Fail(c, "Save file error")
			return
		}
	}
}
