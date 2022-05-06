package processor

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
	return p.Resize(uint(ra.width), uint(ra.height))
}

func NewAction(typ int) Action {
	if typ == Resize {
		return &ResizeAction{}
	}

	return &NopAction{}
}
