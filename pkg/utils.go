package pkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
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

func StrHash(value string) int64 {
	intVal, err := strconv.ParseInt(value[0:3], 16, 64)
	if err != nil {
		return 0
	}

	return intVal / 4
}

func MakeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}
