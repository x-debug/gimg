package pkg

import (
	"fmt"
	"gimg/config"
	"gimg/processor"
	"io"
	"mime/multipart"
	"os"
)

type Ctx struct {
	conf   *config.Config
	Engine processor.Engine
}

func CreateCtx(conf *config.Config, engine processor.Engine) *Ctx {
	return &Ctx{conf: conf, Engine: engine}
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
