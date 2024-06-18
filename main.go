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
	"github.com/hellflame/matrix-ranger/utils"
)

func main() {
	cfg := utils.Parse()
	if cfg == nil {
		return
	}
	go func() {
		window := new(app.Window)

		window.Option(
			app.MaxSize(unit.Dp(cfg.Width), unit.Dp(cfg.Height)),
			app.MinSize(unit.Dp(cfg.Width), unit.Dp(cfg.Height)),
			app.Size(unit.Dp(cfg.Width), unit.Dp(cfg.Height)),
			app.Title("matrix ranger"))
		err := run(window, cfg.BlockSize, cfg.Seed)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window, blockSize, seed int) error {
	ops := new(op.Ops)
	style := styles.CreateStyle(blockSize, 10, 5, 5)
	style.SetTheme("default")
	// stage := new(stage.Stage)

	router := framework.NewRouter(new(pages.Background))
	ctx := framework.NewContext(ops, router)

	// default first page
	router.To("stage")

	router.Add("stage", stage.NewStage(style, int64(seed)))
	// router.Add("home", nil)
	// router.Add("menu", nil)

	for {
		switch e := window.Event().(type) {
		// case app.ViewEvent:
		// 	println("view event")

		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			ctx.NewFrame(&e)
			router.Render(ctx)

			e.Frame(ops)
		}
	}
}
