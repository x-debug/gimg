package handlers

import (
	"fmt"
	"gimg/cache"
	"gimg/logger"
	"gimg/pkg"
	pl "gimg/processor"
	"os"

	"github.com/gin-gonic/gin"
)

// GetHandler Get image
func GetHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		op := c.DefaultQuery("op", "")

		ctx.Logger.Info("Raw query ", logger.String("Query", c.Request.URL.RawQuery), logger.String("Op", op))
		processor := ctx.Engine.NewProcessor(ctx, ctx.Logger, ctx.Conf, hash).SetParam(c.Request.URL.RawQuery)
		//Cache processed image with parameters
		imageBlob, err := ctx.Cache.Get(processor.ActionFinger())
		if err != nil {
			if cache.CacheMiss != err {
				ctx.Logger.Info("Cache brocker get error ", logger.Error(err))
				pkg.Fail(c, fmt.Sprintf("Cache brocker configuare error or shutdown?, error:%s", err.Error()))
				return
			}
			ctx.Logger.Info("Cache pass")
		} else {
			_, err = c.Writer.Write(imageBlob)
			if err != nil {
				pkg.Fail(c, fmt.Sprintf("Write image blob to http stream:%s", err.Error()))
			} else {
				ctx.Logger.Info("Image cache hit", logger.String("CacheKey", processor.ActionFinger()))
			}
			return
		}
		defer processor.Destroy()

		if op == "resize" {
			processor.AddAction(pl.NewAction(pl.Resize))
		} else if op == "thumbnail" {
			processor.AddAction(pl.NewAction(pl.Thumbnail))
		} else if op == "flip" {
			processor.AddAction(pl.NewAction(pl.Flip))
		} else if op == "rotate" {
			processor.AddAction(pl.NewAction(pl.Rotate))
		} else if op == "lua" {
			processor.AddAction(pl.NewAction(pl.LUA))
		} else if op == "gray" {
			processor.AddAction(pl.NewAction(pl.GRAY))
		} else if op == "crop" {
			processor.AddAction(pl.NewAction(pl.CROP))
		} else if op == "quality" {
			processor.AddAction(pl.NewAction(pl.QUALITY))
		} else if op == "format" {
			processor.AddAction(pl.NewAction(pl.FORMAT))
		} else if op == "round" {
			processor.AddAction(pl.NewAction(pl.ROUND))
		} else {
			processor.AddAction(pl.NewAction(pl.Nop))
		}

		rFile, closef, err := processor.ReadCached()
		defer func(f func()) {
			if f != nil {
				closef()
			}
		}(closef)
		if err == nil {
			ctx.Logger.Info("Cached file hit, read from cache ", logger.String("Filename", rFile.Name()))
			ctx.RenderFile(c, processor, rFile)
			return
		}

		rFile, closef, err = processor.Read()
		defer func(f func()) {
			if f != nil {
				closef()
			}
		}(closef)
		if err != nil {
			pkg.Fail(c, fmt.Sprintf("Read original image file error, %s", err.Error()))
			return
		}

		if processor.ActionOnlyNop() {
			ctx.RenderFile(c, processor, rFile)
			return
		}

		ctx.Logger.Info("Processor fit image ", logger.String("FileName", rFile.Name()))
		err = processor.Fit(rFile)
		if err != nil {
			pkg.Fail(c, fmt.Sprintf("Fit image file error, %s", err.Error()))
			return
		}

		//write file object into disk
		wfile, _ := os.CreateTemp("/tmp", "")
		defer wfile.Close()

		err = processor.WriteToFile(wfile)
		if err != nil {
			ctx.Logger.Error("Write Image To File", logger.Error(err))
			pkg.Fail(c, "Copy file error")
			return
		}

		filename := ctx.File(processor.ActionFinger())
		err = os.Rename(wfile.Name(), filename)
		if err != nil {
			pkg.Fail(c, "Rename file error")
			return
		}

		ctx.RenderFile(c, processor, wfile)
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
