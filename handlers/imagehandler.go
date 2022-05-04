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

func renderFile(c *gin.Context, file *os.File) {
	//set file begin of position
	_, _ = file.Seek(0, io.SeekStart)

	wBuffer := &bytes.Buffer{}
	_, err := io.Copy(wBuffer, file)
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

// GetHandler Get image
func GetHandler(ctx *pkg.Ctx) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		width, _ := strconv.Atoi(c.DefaultQuery("w", "0"))
		height, _ := strconv.Atoi(c.DefaultQuery("h", "0"))
		log.Printf("Params hash: %s, width: %d, height: %d\n", hash, width, height)

		fileHash := fmt.Sprintf("%s_w%d_h%d", hash, width, height)
		//Read file from local fs or remote file store
		rFile, closef, err := ctx.ReadFile(fileHash)
		defer func(f func()) {
			if f != nil {
				log.Printf("Close fd: %s\n", fileHash)
				closef()
			}
		}(closef)

		if err == nil {
			log.Printf("Cache hit, hash: %s\n", fileHash)
			renderFile(c, rFile)
			return
		}

		rFile, closef, err = ctx.ReadFile(hash)
		defer func(f func()) {
			if f != nil {
				log.Printf("Close fd: %s\n", hash)
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

		renderFile(c, wfile)
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
