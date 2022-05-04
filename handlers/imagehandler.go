package handlers

import (
	"bytes"
	"fmt"
	"gimg/pkg"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"strconv"
)

// GetHandler Get image
func GetHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")

		//Read file from local fs or remote file store
		rFile, closef, err := ctx.ReadFile(hash)
		defer closef()

		if err != nil {
			pkg.Fail(c, "Read file error")
			return
		}

		width, _ := strconv.Atoi(c.DefaultQuery("w", "0"))
		height, _ := strconv.Atoi(c.DefaultQuery("h", "0"))
		log.Printf("Params hash: %s, width: %d, height: %d\n", hash, width, height)

		if width > 0 && height > 0 {
			processor, err := ctx.Engine.NewProcessor(rFile)
			if err != nil {
				pkg.Fail(c, "Read file error")
				return
			}
			defer processor.Destroy()

			err = processor.Resize(uint(width), uint(height))
			if err != nil {
				pkg.Fail(c, fmt.Sprintf("Resize image file error, %s", err.Error()))
				return
			}

			wfile, _ := os.CreateTemp("/tmp", "")
			defer wfile.Close()
			err = processor.WriteToFile(wfile)
			if err != nil {
				pkg.Fail(c, "Copy file error")
				return
			}

			_, _ = wfile.Seek(0, io.SeekStart)
			md5, _ := pkg.CalcMd5(wfile)
			fileHash := fmt.Sprintf("%s_w%d_h%d", md5, width, height)
			filename := ctx.File(fileHash)
			err = os.Rename(wfile.Name(), filename)
			if err != nil {
				pkg.Fail(c, "Rename file error")
				return
			}

			_, _ = wfile.Seek(0, io.SeekStart)
			rFile = wfile
		}

		wBuffer := &bytes.Buffer{}
		_, err = io.Copy(wBuffer, rFile)
		if err != nil {
			pkg.Fail(c, "Copy buffer error")
			return
		}
		nBytes, err := c.Writer.Write(wBuffer.Bytes())
		log.Printf("Write buffer %d bytes\n", nBytes)
		if err != nil {
			pkg.Fail(c, "Write buffer error")
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
