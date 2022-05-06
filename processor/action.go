package processor

import "gimg/logger"

type ActionCtx struct {
}

type Action interface {
	//Name Operate name
	Name() string
	SetParams(params Params)
	Do(p Processor) error
}

const (
	Nop = iota
	Resize
	Thumbnail
	Flip
	Rotate
)

//NopAction do nothing actually
type NopAction struct {
}

func (na *NopAction) SetParams(params Params) {

}

func (na *NopAction) Name() string {
	return "nop"
}

func (na *NopAction) Do(p Processor) error {
	return nil
}

//ResizeAction can resize image
type ResizeAction struct {
	width  int
	height int
}

func (ra *ResizeAction) SetParams(params Params) {
	ra.width = params.GetInt("w", 0)
	ra.height = params.GetInt("h", 0)
}

func (ra *ResizeAction) Name() string {
	return "resize"
}

func (ra *ResizeAction) Do(p Processor) error {
	p.GetLogger().Info("Resize image file ", logger.Int("Width", ra.width), logger.Int("Height", ra.height))
	return p.Resize(uint(ra.width), uint(ra.height))
}

//ThumbnailAction can generate thumbnail
type ThumbnailAction struct {
	width  int
	height int
}

func (ta *ThumbnailAction) SetParams(params Params) {
	ta.width = params.GetInt("w", 0)
	ta.height = params.GetInt("h", 0)
}

func (ta *ThumbnailAction) Name() string {
	return "thumbnail"
}

func (ta *ThumbnailAction) Do(p Processor) error {
	p.GetLogger().Info("Thumbnail image file ", logger.Int("Width", ta.width), logger.Int("Height", ta.height))
	return p.Thumbnail(uint(ta.width), uint(ta.height))
}

func NewAction(typ int) Action {
	if typ == Resize {
		return &ResizeAction{}
	} else if typ == Thumbnail {
		return &ThumbnailAction{}
	}

	return &NopAction{}
}
