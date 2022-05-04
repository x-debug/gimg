package processor

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
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
func (p *ImagickEngine) NewProcessor(reader io.Reader) (Processor, error) {
	return newImagickProcessor(reader) //fit imagick
}
