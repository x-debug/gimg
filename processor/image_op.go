package processor

type ImageOp interface {
	Resize(width, height uint) error
	Thumbnail(width, height uint) error
	Rotate(deg float64) error
	GrayScale() error
	Crop(x, y int, width, height uint) error
	SetQuality(quality uint) error
}
