package main

import (
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const screenWidth, screenHeight int = 640, 480
const scale int = 1
const brushSize int = 9
const width, height int = screenWidth / scale, screenHeight / scale

type Game struct {
	pixels []byte
	world  []bool
}

func (g *Game) Get(x int, y int) bool {
	if x < 0 || x >= width || y < 0 || y >= height {
		return true
	}
	return g.world[y*width+x]
}

func (g *Game) Set(x int, y int, v bool) {
	g.world[y*width+x] = v
}

func (g *Game) UpdatePixel(x int, y int) {
	// Checa se está vivo
	if g.Get(x, y) {
		// Checa se não está no chão
		if y != height-1 {
			// Checa se não tem nada embaixo
			if !g.Get(x, y+1) {
				// Cai para baixo
				g.Set(x, y+1, true)
				g.Set(x, y, false)
			} else {
				// Se tiver algo embaixo, checa os vizinhos de baixo para deslizar
				neighboorLeft := g.Get(x-1, y+1)
				neighboorRight := g.Get(x+1, y+1)
				switch {
				// Se está vazio dos dois lados, cai em direção aleatória
				case !neighboorLeft && !neighboorRight:
					if rand.Intn(2) == 0 {
						g.Set(x-1, y+1, true)
					} else {
						g.Set(x+1, y+1, true)
					}
					g.Set(x, y, false)
				// Vazio na esquerda
				case !neighboorLeft:
					g.Set(x-1, y+1, true)
					g.Set(x, y, false)
				// Vazio na direita
				case !neighboorRight:
					g.Set(x, y, false)
					g.Set(x+1, y+1, true)
				}
			}
		}
	}
}

func (g *Game) Update() error {
	if g.world == nil {
		g.world = make([]bool, width*height)
	}

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			g.UpdatePixel(x, y)
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		brushLength := brushSize / 2
		for nx := x - brushLength; nx <= x+brushLength; nx++ {
			for ny := y - brushLength; ny <= y+brushLength; ny++ {
				if nx < width && ny < height && nx >= 0 && ny >= 0 {
					g.Set(nx, ny, true)
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, 4*width*height)
	}
	for i, pix := range g.world {
		switch pix {
		case true:
			g.pixels[4*i] = 0xff
			g.pixels[4*i+1] = 0xff
			g.pixels[4*i+2] = 0xff
			g.pixels[4*i+3] = 0xff
		default:
			g.pixels[4*i] = 0
			g.pixels[4*i+1] = 0
			g.pixels[4*i+2] = 0
			g.pixels[4*i+3] = 0
		}
	}
	screen.WritePixels(g.pixels)
	ebitenutil.DebugPrint(screen, strconv.FormatFloat(ebiten.ActualFPS(), 'f', 1, 64))
}

func (g *Game) Layout(outsidescreenWidth, outsidescreenHeight int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sand Simulator")
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
