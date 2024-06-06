package framework

import (
	"fmt"
	"strings"
	"time"
)

const routeSep = "/"

type Router struct {
	paths []string

	route *route
}

type route struct {
	page     Page
	created  bool
	children map[string]*route
}

func (r *route) Add(name string, p Page) *route {
	if r.children == nil {
		r.children = make(map[string]*route)
	}
	r.children[name] = &route{
		page: p,
	}
	return r.children[name]
}

func (r *route) Render(ctx *Context) {
	if !r.created {
		r.page.OnCreate()
		r.created = true
	}
	r.page.Render(ctx)
}

func NewRouter(base Page) *Router {
	r := &Router{
		route: &route{
			page: base,
		},
	}
	return r
}

func (r *Router) Add(name string, p Page) *route {
	return r.route.Add(name, p)
}

func (r *Router) To(name string) bool {
	// windows/ok
	println("to page", name)
	if r.CurrentRoute() != name {
		r.paths = strings.Split(strings.TrimPrefix(name, routeSep), routeSep)
		return true
	}
	return false
}

func (r *Router) CurrentRoute() string {
	return strings.Join(r.paths, routeSep)
}

func (r *Router) Render(ctx *Context) {
	route := r.route
	for _, p := range r.paths {
		route.Render(ctx)
		route = route.children[p]
	}
	route.Render(ctx)
	fmt.Println("render", r.paths, time.Now().String())
}

func (r *Router) Refresh() {

}
