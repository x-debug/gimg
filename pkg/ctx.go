package pkg

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type Ctx struct {
	savePath string
}

func CreateCtx(path string) *Ctx {
	return &Ctx{savePath: path}
}

func (fc *Ctx) ReadFile(fingerprint string) (io.Reader, func(), error) {
	filename := fmt.Sprintf("%s/%s", fc.savePath, fingerprint)
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		_ = file.Close()
	}, nil
}

//SaveFile save file to path of fingerprint
func (fc *Ctx) SaveFile(fingerprint string, file multipart.File) error {
	filename := fmt.Sprintf("%s/%s", fc.savePath, fingerprint)
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
