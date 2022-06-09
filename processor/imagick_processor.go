package processor

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/x-debug/gimg/config"
	"github.com/x-debug/gimg/fs"
	"github.com/x-debug/gimg/logger"

	"gopkg.in/gographics/imagick.v3/imagick"
)

const HTTP_SCHEMA = "http://"
const HTTPS_SCHEMA = "https://"

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
	h := md5.New()
	joinParams := make([]string, 0)
	for _, param := range p.params {
		//File system path can't contain the special char, so Hash(value) if query param is url
		if strings.Contains(param.Value, HTTP_SCHEMA) || strings.Contains(param.Value, HTTPS_SCHEMA) {
			h.Reset()
			h.Write([]byte(param.Value))
			hash := fmt.Sprintf("%x", h.Sum(nil))
			joinParams = append(joinParams, param.Key+hash)
		} else {
			joinParams = append(joinParams, param.Key+param.Value)
		}
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

func (p *ImagickProcessor) SetFormat(format string) error {
	return p.mw.SetFormat(format)
}

func (p *ImagickProcessor) RoundCorner(rx, ry float64) error {
	p.logger.Info("RoundCorner", logger.Float64("RX", rx), logger.Float64("RY", ry))

	mw := imagick.NewMagickWand()
	pw := imagick.NewPixelWand()
	dw := imagick.NewDrawingWand()

	imgWidth := p.mw.GetImageWidth()
	imgHeight := p.mw.GetImageHeight()

	p.logger.Info("RoundCorner", logger.Int("Image Width", int(imgWidth)), logger.Int("Image Height", int(imgHeight)))
	// Create the initial 640x480 transparent canvas
	pw.SetColor("none")
	mw.NewImage(imgWidth, imgHeight, pw)

	pw.SetColor("white")
	dw.SetFillColor(pw)
	dw.RoundRectangle(0, 0, float64(imgWidth-1), float64(imgHeight-1), rx, ry)
	mw.DrawImage(dw)

	// Note that MagickSetImageCompose is usually only used for the MagickMontageImage
	// function and isn't used or needed by MagickCompositeImage
	err := mw.CompositeImage(p.mw, imagick.COMPOSITE_OP_SRC_IN, true, 0, 0)
	p.logger.Info("Image format", logger.String("format", p.mw.GetFormat()))
	mw.SetFormat("png") //Set png format
	p.mw = mw
	return err
}

func (p *ImagickProcessor) SetupActions(op string) {
	if op == "resize" {
		p.AddAction(NewAction(Resize))
	} else if op == "thumbnail" {
		p.AddAction(NewAction(Thumbnail))
	} else if op == "flip" {
		p.AddAction(NewAction(Flip))
	} else if op == "rotate" {
		p.AddAction(NewAction(Rotate))
	} else if op == "lua" {
		p.AddAction(NewAction(LUA))
	} else if op == "gray" {
		p.AddAction(NewAction(GRAY))
	} else if op == "crop" {
		p.AddAction(NewAction(CROP))
	} else if op == "quality" {
		p.AddAction(NewAction(QUALITY))
	} else if op == "format" {
		p.AddAction(NewAction(FORMAT))
	} else if op == "round" {
		p.AddAction(NewAction(ROUND))
	} else {
		p.AddAction(NewAction(Nop))
	}
}
