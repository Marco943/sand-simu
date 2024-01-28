package main

import (
	"math/rand"
	"strconv"

	"github.com/crazy3lf/colorconv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  int     = 640
	screenHeight         = 480
	scale                = 5
	brushSize            = 3
	width                = screenWidth / scale
	height               = screenHeight / scale
	S            float64 = 1
	L            float64 = 0.5
)

var H float64 = 1

type Game struct {
	pixels []byte
	world  []float64
}

func (g *Game) Get(x int, y int) float64 {
	if x < 0 || x >= width || y < 0 || y >= height {
		return 1
	}
	return g.world[y*width+x]
}

func (g *Game) Set(x int, y int, v float64) {
	g.world[y*width+x] = v
}

func (g *Game) UpdatePixel(x int, y int) {
	vAtual := g.Get(x, y)
	// Checa se está vivo
	if vAtual > 0 {
		// Checa se não está no chão
		if y != height-1 {
			// Checa se não tem nada embaixo
			if g.Get(x, y+1) == 0 {
				// Cai para baixo
				g.Set(x, y+1, vAtual)
				g.Set(x, y, 0)
			} else {
				// Se tiver algo embaixo, checa os vizinhos de baixo para deslizar
				neighboorLeft := g.Get(x-1, y+1) > 0
				neighboorRight := g.Get(x+1, y+1) > 0
				switch {
				// Se está vazio dos dois lados, cai em direção aleatória
				case !neighboorLeft && !neighboorRight:
					if rand.Intn(2) == 0 {
						g.Set(x-1, y+1, vAtual)
					} else {
						g.Set(x+1, y+1, vAtual)
					}
					g.Set(x, y, vAtual)
				// Vazio na esquerda
				case !neighboorLeft:
					g.Set(x-1, y+1, vAtual)
					g.Set(x, y, 0)
				// Vazio na direita
				case !neighboorRight:
					g.Set(x, y, 0)
					g.Set(x+1, y+1, vAtual)
				}
			}
		}
	}
}

func (g *Game) Update() error {
	if g.world == nil {
		g.world = make([]float64, width*height)
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
				if nx < width && ny < height && nx >= 0 && ny >= 0 && g.Get(nx, ny) == 0 {
					if H <= 358 {
						H += 0.2
					} else {
						H = 1
					}
					g.Set(nx, ny, H)
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
		red, green, blue, _ := colorconv.HSLToRGB(float64(pix), S, L)
		switch {
		case pix > 0:
			g.pixels[4*i] = byte(red)
			g.pixels[4*i+1] = byte(green)
			g.pixels[4*i+2] = byte(blue)
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
	ebitenutil.DebugPrintAt(screen, strconv.FormatInt(int64(H), 10), width-100, 0)
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
