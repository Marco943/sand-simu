package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const WIDTH, HEIGHT int = 640, 480

type Game struct {
	pixels []byte
	world  []bool
}

func (g *Game) Update() error {
	g.world = make([]bool, HEIGHT*WIDTH)
	randPix := rand.Intn(HEIGHT*WIDTH) + 1
	g.world[randPix] = true
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fmt.Println(x, y)
		g.world[y*WIDTH+x] = true
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, 4*HEIGHT*WIDTH)
	}
	for i, pix := range g.world {
		if pix {
			g.pixels[4*i] = 0xff
			g.pixels[4*i+1] = 0xff
			g.pixels[4*i+2] = 0xff
			g.pixels[4*i+3] = 0xff
		} else {
			g.pixels[4*i] = 0
			g.pixels[4*i+1] = 0
			g.pixels[4*i+2] = 0
			g.pixels[4*i+3] = 0
		}
	}
	screen.WritePixels(g.pixels)
	ebitenutil.DebugPrint(screen, strconv.FormatFloat(ebiten.ActualFPS(), 'f', 1, 64))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Sand Simulator")
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
