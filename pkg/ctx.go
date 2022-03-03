package pkg

import "mime/multipart"

type Ctx struct {
	savePath string
}

func CreateCtx(path string) *Ctx {
	return &Ctx{savePath: path}
}

func (fc *Ctx) SaveFile(hash string, file multipart.File) error {
	return nil
}
