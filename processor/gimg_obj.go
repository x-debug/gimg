package processor

//GimgObj return new lua image object with image op, and request object
type GimgObj struct {
	processor Processor
}

func NewFromProcessor(processor Processor) *GimgObj {
	return &GimgObj{processor: processor}
}

func (g *GimgObj) Resize(width, height uint) error {
	return g.processor.Resize(width, height)
}

func (g *GimgObj) Thumbnail(width, height uint) error {
	return g.processor.Thumbnail(width, height)
}

func (g *GimgObj) Rotate(deg float64) error {
	return g.processor.Rotate(deg)
}
