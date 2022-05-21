package processor

import (
	"bytes"
	"fmt"
	"gimg/config"
	"gimg/fs"
	"gimg/logger"
	"io"
	"os"
	"strings"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type ImagickProcessor struct {
	mw         *imagick.MagickWand
	params     Params
	fileHash   string
	actions    []Action
	fs         fs.FileSystem
	logger     logger.Logger
	actionConf *config.ActionConf
}

func newImagickProcessor(fs fs.FileSystem, logger logger.Logger, conf *config.Config, hash string) Processor {
	mw := imagick.NewMagickWand()

	self := &ImagickProcessor{mw: mw, fs: fs, fileHash: hash, logger: logger, actions: make([]Action, 0), actionConf: conf.Action}
	return self
}

func (p *ImagickProcessor) GetLogger() logger.Logger {
	return p.logger
}

func (p *ImagickProcessor) GetActionConf() *config.ActionConf {
	return p.actionConf
}

func (p *ImagickProcessor) Load(file *os.File) error {
	buffer := &bytes.Buffer{}
	_, err := io.Copy(buffer, file)
	if err != nil {
		return err
	}

	err = p.mw.ReadImageBlob(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (p *ImagickProcessor) SetParam(param string) Processor {
	p.params = NewParams(param)
	return p
}

func (p *ImagickProcessor) AddAction(action Action) {
	p.actions = append(p.actions, action)
}

func (p *ImagickProcessor) LenOfAction() int {
	return len(p.actions)
}

func (p *ImagickProcessor) ActionOnlyNop() bool {
	return len(p.actions) == 0 || (len(p.actions) == 1 && p.actions[0].Name() == "nop")
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

func (p *ImagickProcessor) ActionFinger() string {
	joinParams := make([]string, 0)
	for _, param := range p.params {
		joinParams = append(joinParams, param.Key+param.Value)
	}
	newHash := fmt.Sprintf("%s_%s", p.fileHash, strings.Join(joinParams, "_"))
	return newHash
}

//ReadCached return cached file object
func (p *ImagickProcessor) ReadCached() (*os.File, func(), error) {
	newHash := p.ActionFinger()
	return p.fs.ReadFile(newHash)
}

//Read return file object
func (p *ImagickProcessor) Read() (*os.File, func(), error) {
	return p.fs.ReadFile(p.fileHash)
}

//Fit process image object with actions
func (p *ImagickProcessor) Fit(file *os.File) error {
	err := p.Load(file)
	if err != nil {
		return err
	}

	for _, action := range p.actions {
		action.SetParams(p.params)
		err = action.Do(p)
		if err != nil {
			return err
		}
	}

	return nil
}

//Resize return new sized image object with width and height
func (p *ImagickProcessor) Resize(width, height uint) error {
	return p.mw.ResizeImage(width, height, imagick.FILTER_LANCZOS)
}

//Thumbnail return new thumbnail image object with width and height
func (p *ImagickProcessor) Thumbnail(width, height uint) error {
	return p.mw.ThumbnailImage(width, height)
}

func (p *ImagickProcessor) Rotate(deg float64) error {
	return p.mw.RotateImage(imagick.NewPixelWand(), deg)
}

func (p *ImagickProcessor) GrayScale() error {
	return p.mw.SetImageType(imagick.IMAGE_TYPE_GRAYSCALE)
}

func (p *ImagickProcessor) Crop(x, y int, width, height uint) error {
	return p.mw.CropImage(width, height, x, y)
}

func (p *ImagickProcessor) SetQuality(quality uint) error {
	return p.mw.SetCompressionQuality(quality)
}
