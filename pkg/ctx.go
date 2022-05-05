package pkg

import (
	"bytes"
	"fmt"
	"gimg/config"
	"gimg/logger"
	lg "gimg/logger"
	"gimg/processor"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
)

type Ctx struct {
	conf   *config.Config
	Engine processor.Engine
	Logger logger.Logger
}

func CreateCtx(conf *config.Config, logger logger.Logger, engine processor.Engine) *Ctx {
	return &Ctx{conf: conf, Engine: engine, Logger: logger}
}

func (fc *Ctx) RenderFile(c *gin.Context, file *os.File) {
	//set file begin of position
	_, _ = file.Seek(0, io.SeekStart)

	wBuffer := &bytes.Buffer{}
	_, err := io.Copy(wBuffer, file)
	if err != nil {
		Fail(c, "Copy buffer error")
		return
	}

	nBytes, err := c.Writer.Write(wBuffer.Bytes())
	fc.Logger.Info("Write buffer", lg.Int("Bytes", nBytes))
	if err != nil {
		Fail(c, "Write buffer error")
		return
	}
}

func (fc *Ctx) ReadFile(fingerprint string) (*os.File, func(), error) {
	filename := fmt.Sprintf("%s/%s", fc.conf.Engine.SavePath, fingerprint)
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		_ = file.Close()
	}, nil
}

func (fc *Ctx) File(filename string) string {
	return fmt.Sprintf("%s/%s", fc.conf.Engine.SavePath, filename)
}

//SaveFile save file to path of fingerprint
func (fc *Ctx) SaveFile(fingerprint string, file multipart.File) error {
	filename := fmt.Sprintf("%s/%s", fc.conf.Engine.SavePath, fingerprint)
	dstFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	//_, err = io.Copy(dstFile, file)
	buffer := make([]byte, 256)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	for {
		_, err := file.Read(buffer)
		//log.Printf("read %d bytes, err: %s\n", nBytes, err)
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
