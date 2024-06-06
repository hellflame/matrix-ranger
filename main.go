package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/hellflame/matrix-ranger/framework"
	"github.com/hellflame/matrix-ranger/pages"
	"github.com/hellflame/matrix-ranger/pages/stage"
	"github.com/hellflame/matrix-ranger/styles"
)

func main() {
	go func() {
		window := new(app.Window)

		window.Option(
			app.MaxSize(unit.Dp(500), unit.Dp(580)),
			app.MinSize(unit.Dp(500), unit.Dp(580)),
			app.Size(unit.Dp(500), unit.Dp(580)),
			app.Title("matrix ranger"))
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	ops := new(op.Ops)
	style := styles.CreateStyle(70, 10, 5, 5)
	style.SetTheme("default")
	// stage := new(stage.Stage)

	router := framework.NewRouter(new(pages.Background))
	ctx := framework.NewContext(ops, router)

	// default first page
	router.To("stage")

	router.Add("stage", stage.NewStage(style, 9))
	// router.Add("home", nil)
	// router.Add("menu", nil)

	for {
		switch e := window.Event().(type) {
		// case app.ViewEvent:
		// 	println("view event")

		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			ctx.Refresh(&e)

			router.Render(ctx)

			e.Frame(ops)
		}
	}
}
