package framework

type Page interface {
	OnCreate()
	Render(ctx *Context)
}
