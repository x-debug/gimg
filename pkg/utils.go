package pkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
)

func CalcMd5File(file multipart.File) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, file)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func CalcMd5Str(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}
