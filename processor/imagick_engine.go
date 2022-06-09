package processor

import (
	"github.com/x-debug/gimg/config"
	"github.com/x-debug/gimg/fs"
	"github.com/x-debug/gimg/logger"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type ImagickEngine struct {
	cfg *config.EngineConf
}

func (p *ImagickEngine) Initialize() {
	imagick.Initialize()
}

func (p *ImagickEngine) Terminate() {
	imagick.Terminate()
}

func (p *ImagickEngine) GetConfig() *config.EngineConf {
	return p.cfg
}

//NewProcessor build processor
func (p *ImagickEngine) NewProcessor(fs fs.FileSystem, logger logger.Logger, conf *config.Config, originalHash string) Processor {
	return newImagickProcessor(fs, logger, conf, originalHash) //fit imagick
}
