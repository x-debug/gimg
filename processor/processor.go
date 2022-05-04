package processor

import (
	"io"
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
	NewProcessor(reader io.Reader) (Processor, error)
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

type Processor interface {
	FileW

	Resize(width, height uint) error
	Destroy()
}
