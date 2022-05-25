package pkg

import (
	"bytes"
	"fmt"
	"gimg/cache"
	"gimg/config"
	"gimg/logger"
	lg "gimg/logger"
	"gimg/processor"
	"io"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
)

type Ctx struct {
	Conf   *config.Config
	Engine processor.Engine
	Logger logger.Logger
	Cache  cache.Cache
}

func CreateCtx(conf *config.Config, cache cache.Cache, logger logger.Logger, engine processor.Engine) *Ctx {
	return &Ctx{Conf: conf, Cache: cache, Engine: engine, Logger: logger}
}

func (fc *Ctx) RenderFile(c *gin.Context, finger processor.HttpFinger, file *os.File) {
	//set file begin of position
	_, _ = file.Seek(0, io.SeekStart)

	wBuffer := &bytes.Buffer{}
	_, err := io.Copy(wBuffer, file)
	if err != nil {
		Fail(c, "Copy buffer error")
		return
	}

	fc.Logger.Info("Write cache brocker")
	//Ignore the error
	_ = fc.Cache.Set(finger.ActionFinger(), wBuffer.Bytes())
	/*
		if err != nil {
			Fail(c, "Write buffer error")
			return
		}
	*/

	nBytes, err := c.Writer.Write(wBuffer.Bytes())
	fc.Logger.Info("Write buffer", lg.Int("Bytes", nBytes))
	if err != nil {
		Fail(c, "Write buffer error")
		return
	}
}

func (fc *Ctx) ReadFile(fingerprint string) (*os.File, func(), error) {
	filename := fmt.Sprintf("%s/%s", fc.Conf.Engine.SavePath, fingerprint)
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		_ = file.Close()
	}, nil
}

func (fc *Ctx) File(filename string) string {
	return fmt.Sprintf("%s/%s", fc.Conf.Engine.SavePath, filename)
}

//SaveFile save file to path of fingerprint
func (fc *Ctx) SaveFile(fingerprint string, file multipart.File) error {
	filename := fmt.Sprintf("%s/%s", fc.Conf.Engine.SavePath, fingerprint)
	dstFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	buffer := make([]byte, 256)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	for {
		_, err := file.Read(buffer)
		if err == io.EOF { //file reached EOF, stop reading
			break
		}
		_, err = dstFile.Write(buffer)
		if err != nil {
			return err
		}
	}

	_ = dstFile.Close()
	return err
}
