package processor

import (
	"gimg/config"
	"gimg/fs"
	"gimg/logger"
	"os"
)

const Imagick = iota

//Engine is image process engine
type Engine interface {
	//Initialize You must call it, processor will initialize some unmanaged resources with no gc
	Initialize()

	//Terminate You must call it when program exit, otherwise will resource leak
	Terminate()

	//NewProcessor build a new processor
	NewProcessor(fs fs.FileSystem, logger logger.Logger, conf *config.Config, originalHash string) Processor

	GetConfig() *config.EngineConf
}

//NewEngine build processor with type
func NewEngine(typ int, config *config.EngineConf) Engine {
	if typ == Imagick {
		return &ImagickEngine{cfg: config}
	}

	return nil
}

type FileW interface {
	WriteTo(filename string) error
	WriteToFile(fObj *os.File) error
}

type Loggable interface {
	GetLogger() logger.Logger
}

type Processor interface {
	FileW
	ImageOp
	Loggable
	HttpFinger

	Load(file *os.File) error
	SetParam(param string) Processor
	GetActionConf() *config.ActionConf
	Destroy()
	AddAction(action Action)
	SetupActions(typ string)
	LenOfAction() int
	ActionOnlyNop() bool
	ReadCached() (*os.File, func(), error)
	Read() (*os.File, func(), error)
	Fit(file *os.File) error
}
