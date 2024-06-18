package utils

import (
	"fmt"

	"github.com/hellflame/argparse"
)

type Config struct {
	Width  int
	Height int

	BlockSize int

	Seed int
}

func Parse() *Config {
	p := argparse.NewParser("", "", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})
	width := p.Int("", "width", &argparse.Option{Default: "500"})
	height := p.Int("", "height", &argparse.Option{Default: "600"})

	blockSize := p.Int("s", "size", &argparse.Option{Default: "70"})
	seed := p.Int("", "seed", &argparse.Option{Default: "10"})

	if e := p.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return nil
	}
	return &Config{
		Height: *height, Width: *width,
		BlockSize: *blockSize, Seed: *seed,
	}
}
