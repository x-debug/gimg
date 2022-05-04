package processor

import (
	"bytes"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"os"
)

type ImagickProcessor struct {
	mw *imagick.MagickWand
}

func newImagickProcessor(reader io.Reader) (Processor, error) {
	buffer := &bytes.Buffer{}
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return nil, err
	}

	mw := imagick.NewMagickWand()
	err = mw.ReadImageBlob(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return &ImagickProcessor{mw: mw}, nil
}

func (p *ImagickProcessor) WriteTo(filename string) error {
	return p.mw.WriteImage(filename)
}

func (p *ImagickProcessor) WriteToFile(fObj *os.File) error {
	return p.mw.WriteImageFile(fObj)
}

func (p *ImagickProcessor) Destroy() {
	if p.mw != nil {
		p.mw.Destroy()
	}
}

func (p *ImagickProcessor) Resize(width, height uint) error {
	return p.mw.ResizeImage(width, height, imagick.FILTER_LANCZOS)
}
