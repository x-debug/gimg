package processor

import (
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
	NewProcessor(fs fs.FileSystem, logger logger.Logger, originalHash string) Processor
}

//NewEngine build processor with type
func NewEngine(typ int) Engine {
	if typ == Imagick {
		return &ImagickEngine{}
	}

	return nil
}

type FileW interface {
	WriteTo(filename string) error
	WriteToFile(fObj *os.File) error
}

type ImageOp interface {
	Load(file *os.File) error
	Resize(width, height uint) error
}

type Processor interface {
	FileW
	ImageOp

	SetParam(param string) Processor
	Destroy()
	AddAction(action Action)
	ActionFinger() string
	LenOfAction() int
	ActionOnlyNop() bool
	ReadCached() (*os.File, func(), error)
	Read() (*os.File, func(), error)
	Fit(file *os.File) error
}
