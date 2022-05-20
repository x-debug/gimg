package processor

type ImageOp interface {
	Resize(width, height uint) error
	Thumbnail(width, height uint) error
	Rotate(deg float64) error
}
