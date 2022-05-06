package processor

import (
	"gimg/fs"
	"gimg/logger"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type ImagickEngine struct {
}

func (p *ImagickEngine) Initialize() {
	imagick.Initialize()
}

func (p *ImagickEngine) Terminate() {
	imagick.Terminate()
}

//NewProcessor build processor
func (p *ImagickEngine) NewProcessor(fs fs.FileSystem, logger logger.Logger, originalHash string) Processor {
	return newImagickProcessor(fs, logger, originalHash) //fit imagick
}
