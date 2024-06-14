package framework

import (
	"time"

	"gioui.org/app"
	"gioui.org/op"
)

type Context struct {
	lastT int64

	R *Router

	DeltaT float32
	Ops    *op.Ops
	Event  *app.FrameEvent
}

func NewContext(ops *op.Ops, r *Router) *Context {
	return &Context{
		Ops: ops,
		R:   r,

		lastT: time.Now().UnixMilli(),
	}
}

func (c *Context) NewFrame(e *app.FrameEvent) {
	c.Ops.Reset()
	c.Event = e

	present := time.Now().UnixMilli()
	c.DeltaT = float32(present - c.lastT)
	c.lastT = present
}

func (c *Context) Refresh() {
	c.Event.Source.Execute(op.InvalidateCmd{})
}

func (c *Context) RouteTo(target string) {
	if c.R.To(target) {
		c.Refresh()
	}
}
