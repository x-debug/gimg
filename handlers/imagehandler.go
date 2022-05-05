package handlers

import (
	"fmt"
	lg "gimg/logger"
	"gimg/pkg"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

// GetHandler Get image
func GetHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		width, _ := strconv.Atoi(c.DefaultQuery("w", "0"))
		height, _ := strconv.Atoi(c.DefaultQuery("h", "0"))
		ctx.Logger.Info("Params hash", lg.Int("Width", width), lg.Int("Height", height))

		fileHash := fmt.Sprintf("%s_w%d_h%d", hash, width, height)
		//Read file from local fs or remote file store
		rFile, closef, err := ctx.ReadFile(fileHash)
		defer func(f func()) {
			if f != nil {
				closef()
			}
		}(closef)

		if err == nil {
			ctx.RenderFile(c, rFile)
			return
		}

		rFile, closef, err = ctx.ReadFile(hash)
		defer func(f func()) {
			if f != nil {
				closef()
			}
		}(closef)

		processor, err := ctx.Engine.NewProcessor(rFile)
		if err != nil {
			pkg.Fail(c, "Read file error")
			return
		}
		defer processor.Destroy()

		if width > 0 && height > 0 {
			err = processor.Resize(uint(width), uint(height))
			if err != nil {
				pkg.Fail(c, fmt.Sprintf("Resize image file error, %s", err.Error()))
				return
			}
		}

		//write file object into disk
		wfile, _ := os.CreateTemp("/tmp", "")
		defer wfile.Close()

		err = processor.WriteToFile(wfile)
		if err != nil {
			pkg.Fail(c, "Copy file error")
			return
		}

		filename := ctx.File(fileHash)
		err = os.Rename(wfile.Name(), filename)
		if err != nil {
			pkg.Fail(c, "Rename file error")
			return
		}

		ctx.RenderFile(c, wfile)
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

		_, _ = c.Writer.WriteString(md5)
	}
}
