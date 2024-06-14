package framework

type Page interface {
	OnCreate(ctx *Context)
	Render(ctx *Context)
}
