package processor

import (
	"errors"
	"gimg/logger"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

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
	LUA
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

//RotateAction can rotate image file
type RotateAction struct {
	deg float64
}

func (ra *RotateAction) SetParams(params Params) {
	ra.deg = params.GetFloat64("deg", 0.0)
}

func (ra *RotateAction) Name() string {
	return "rotate"
}

func (ra *RotateAction) Do(p Processor) error {
	return p.Rotate(ra.deg)
}

//LuaAction can custom with lua
type LuaAction struct {
	scriptName string
}

func (la *LuaAction) SetParams(params Params) {
	la.scriptName = params.GetString("f", "")
}

func (la *LuaAction) Name() string {
	return "lua"
}

func (la *LuaAction) Do(p Processor) error {
	imgObj := NewFromProcessor(p)
	L := lua.NewState()
	defer L.Close()

	if la.scriptName == "" {
		return errors.New("script name is error")
	}

	L.SetGlobal("G", luar.New(L, imgObj))
	p.GetLogger().Info("Lua script ", logger.String("ScriptName", la.scriptName))
	conf := p.GetActionConf()
	if conf == nil {
		p.GetLogger().Error("Load lua action conf error")
		return errors.New("load lua conf error")
	}

	filename := conf.LoadScriptPath + "/" + la.scriptName + ".lua"
	if err := L.DoFile(filename); err != nil {
		return err
	}
	return nil
}

func NewAction(typ int) Action {
	if typ == Resize {
		return &ResizeAction{}
	} else if typ == Thumbnail {
		return &ThumbnailAction{}
	} else if typ == Rotate {
		return &RotateAction{}
	} else if typ == LUA {
		return &LuaAction{}
	}

	return &NopAction{}
}
