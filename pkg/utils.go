package pkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
)

func CalcMd5(file multipart.File) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, file)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
