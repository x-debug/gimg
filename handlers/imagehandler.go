package handlers

import (
	"fmt"
	"gimg/cache"
	"gimg/logger"
	"gimg/pkg"
	"os"

	"github.com/gin-gonic/gin"
)

// GetHandler Get image
func GetHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		op := c.DefaultQuery("op", "")
		remote := c.DefaultQuery("remote", "")
		ctx.Logger.Info("Raw query ", logger.String("Query", c.Request.URL.RawQuery), logger.String("Op", op))

		//Client set file hash or set remote file url, IF set the remote file url,
		//Engine will download the image file into savepath before process the image file
		if remote != "" {
			ctx.Logger.Info("Fetch remote file", logger.String("RemoteUrl", remote))
			proxy := pkg.NewProxy(ctx.Conf.Engine.SavePath, ctx.Conf.Proxy, ctx.Logger)
			req := proxy.CloneRequest(c.Request)
			hash = req.HashVal()
			if err := proxy.Do(req, hash); err != nil {
				pkg.Fail(c, fmt.Sprintf("download file error, reason of error: %s", err.Error()))
				return
			}
		}

		processor := ctx.Engine.NewProcessor(ctx, ctx.Logger, ctx.Conf, hash).SetParam(c.Request.URL.RawQuery)
		processor.SetupActions(op)
		defer processor.Destroy()

		//Hit high cache? the cache save image with some parameters
		imageBlob, err := ctx.Cache.Get(processor.ActionFinger())
		if err != nil {
			if cache.CacheMiss != err {
				ctx.Logger.Info("Cache brocker get error ", logger.Error(err))
				pkg.Fail(c, fmt.Sprintf("Cache brocker configuare error or shutdown?, reason of error: %s", err.Error()))
				return
			}
			ctx.Logger.Info("Cache pass")
		} else {
			_, err = c.Writer.Write(imageBlob)
			if err != nil {
				pkg.Fail(c, fmt.Sprintf("Write image blob to http stream, reason of error: %s", err.Error()))
			} else {
				ctx.Logger.Info("Image cache hit", logger.String("CacheKey", processor.ActionFinger()))
			}
			return
		}

		//Hit file cache? the cache cost some io performance
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

		//Read origin image, processor will cache these result for performance after processed
		rFile, closef, err = processor.Read()
		defer func(f func()) {
			if f != nil {
				closef()
			}
		}(closef)
		if err != nil {
			pkg.Fail(c, fmt.Sprintf("Read original image file error, reason of error: %s", err.Error()))
			return
		}

		if processor.ActionOnlyNop() {
			ctx.RenderFile(c, processor, rFile)
			return
		}

		ctx.Logger.Info("Processor fit image ", logger.String("FileName", rFile.Name()))
		err = processor.Fit(rFile)
		if err != nil {
			pkg.Fail(c, fmt.Sprintf("Fit image file error, reason of error: %s", err.Error()))
			return
		}

		//Write file object into disk, and write back the stream to body of response
		//Now, the image of processed MUST write into the files, BUT they can communicate through memory stream
		//TODO communicate through memory stream
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

		md5, _ := pkg.CalcMd5File(fObj)
		err = ctx.SaveFile(md5, fObj)
		if err != nil {
			pkg.Fail(c, "Save file error")
			return
		}

		_, _ = c.Writer.WriteString(md5)
	}
}
