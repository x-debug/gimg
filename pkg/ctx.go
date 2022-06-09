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
	"strings"

	"github.com/gin-gonic/gin"
)

const ImageParamSpliter = "_"

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
		fc.Logger.Error("Copy buffer", lg.Error(err))
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

//GetStoragePath return IF fingerprint is original image,set L1(3Byte) and L2(3Byte) for hierarchical directory
func (fc *Ctx) GetStoragePath(fingerprint string) (string, error) {
	originalDir := GetOrignalFileName(fingerprint)
	if strings.Contains(fingerprint, ImageParamSpliter) {
		if err := MakeDirectoryIfNotExists(fc.Conf.Engine.CachePath); err != nil {
			return "", err
		}

		l1Pair := StrHash(fingerprint)
		l2Pair := StrHash(fingerprint[3:])
		pathL1 := fmt.Sprintf("%s/%d", fc.Conf.Engine.CachePath, l1Pair)
		if err := MakeDirectoryIfNotExists(pathL1); err != nil {
			return "", err
		}
		pathL2 := fmt.Sprintf("%s/%d/%d", fc.Conf.Engine.CachePath, l1Pair, l2Pair)
		if err := MakeDirectoryIfNotExists(pathL2); err != nil {
			return "", err
		}
		fileDir := fmt.Sprintf("%s/%d/%d/%s", fc.Conf.Engine.CachePath, l1Pair, l2Pair, originalDir)
		if err := MakeDirectoryIfNotExists(fileDir); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/%d/%d/%s/%s", fc.Conf.Engine.CachePath, l1Pair, l2Pair, originalDir, fingerprint), nil
	} else {
		if err := MakeDirectoryIfNotExists(fc.Conf.Engine.SavePath); err != nil {
			return "", err
		}

		l1Pair := StrHash(fingerprint)
		l2Pair := StrHash(fingerprint[3:])
		pathL1 := fmt.Sprintf("%s/%d", fc.Conf.Engine.SavePath, l1Pair)
		if err := MakeDirectoryIfNotExists(pathL1); err != nil {
			return "", err
		}
		pathL2 := fmt.Sprintf("%s/%d/%d", fc.Conf.Engine.SavePath, l1Pair, l2Pair)
		if err := MakeDirectoryIfNotExists(pathL2); err != nil {
			return "", err
		}
		fileDir := fmt.Sprintf("%s/%d/%d/%s", fc.Conf.Engine.SavePath, l1Pair, l2Pair, originalDir)
		if err := MakeDirectoryIfNotExists(fileDir); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/%d/%d/%s/%s", fc.Conf.Engine.SavePath, l1Pair, l2Pair, originalDir, fingerprint), nil
	}
}

func (fc *Ctx) ReadFile(fingerprint string) (*os.File, func(), error) {
	filename, err := fc.GetStoragePath(fingerprint)
	if err != nil {
		return nil, nil, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		_ = file.Close()
	}, nil
}

func (fc *Ctx) File(filename string) string {
	path, err := fc.GetStoragePath(filename)
	if err != nil {
		return ""
	}
	return path
}

//SaveFile save file to path of fingerprint
func (fc *Ctx) SaveFile(fingerprint string, file multipart.File) error {
	filename, err := fc.GetStoragePath(fingerprint)
	if err != nil {
		return err
	}
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
