package fs

import (
	"mime/multipart"
	"os"
)

type FileSystem interface {
	ReadFile(fingerprint string) (*os.File, func(), error)
	File(filename string) string
	SaveFile(fingerprint string, file multipart.File) error
}
