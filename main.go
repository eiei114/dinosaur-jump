package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	debug    = false
	screenX  = 640
	screenY  = 480
	groundY  = 400
	fontSize = 10

	// game modes
	modeTitle    = 0
	modeLogin    = 1
	modeGame     = 2
	modeGameover = 3

	// image sizes
	dinosaurHeight = 50
	dinosaurWidth  = 100
	groundHeight   = 50
	groundWidth    = 50
)

//go:embed resources/images/dinosaur_01.png
var byteDinosaur1Img []byte

//go:embed resources/images/dinosaur_02.png
var byteDinosaur2Img []byte

//go:embed resources/images/ground.png
var byteGroundImg []byte

var (
	dinosaur1Img *ebiten.Image
	dinosaur2Img *ebiten.Image
	groundImg    *ebiten.Image
	arcadeFont   font.Face
)

func init() {
	rand.Seed(time.Now().UnixNano())

	img, _, err := image.Decode(bytes.NewReader(byteDinosaur1Img))
	if err != nil {
		log.Fatal(err)
	}
	dinosaur1Img = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(byteDinosaur2Img))
	if err != nil {
		log.Fatal(err)
	}
	dinosaur2Img = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(byteGroundImg))
	if err != nil {
		log.Fatal(err)
	}
	groundImg = ebiten.NewImageFromImage(img)

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

type ground struct {
	x int
	y int
}

func (g *ground) move(speed int) {
	g.x -= speed
	if g.x < -groundWidth {
		g.x = g.x + groundWidth
	}
}

// Game struct
type Game struct {
	mode      int
	count     int
	score     int
	hiscore   int
	dinosaurX int
	dinosaurY int
	ground    *ground
	runes     []rune
	text      string
	counter   int
}

// NewGame method
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

// Init method
func (g *Game) init() {
	g.hiscore = g.score
	g.count = 0
	g.score = 0
	g.dinosaurX = 100
	g.dinosaurY = 100
	g.ground = &ground{y: groundY - 30}
}

// Update method
func (g *Game) Update() error {
	switch g.mode {
	case modeTitle:
		if g.isKeySpaceJustPressed() {
			g.mode = modeLogin
		}
	case modeLogin:
		g.runes = ebiten.AppendInputChars(g.runes[:0])
		g.text += string(g.runes)
		// Adjust the string to be at most 10 lines.
		ss := strings.Split(g.text, "\n")
		if len(ss) > 10 {
			g.text = strings.Join(ss[len(ss)-10:], "\n")
		}

		// If the backspace key is pressed, remove one character.
		if repeatingKeyPressed(ebiten.KeyBackspace) {
			if len(g.text) >= 1 {
				g.text = g.text[:len(g.text)-1]
			}
		}

		g.counter++

		if g.isKeySpaceJustPressed() {
			g.mode = modeGame
		}
	case modeGame:
		//todo 障害物に当たったらゲームオーバーになるようにする
		g.count++
		g.score = g.count / 5

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.dinosaurY -= 100
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.dinosaurY += 100
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			g.dinosaurX += 100
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			g.dinosaurX -= 100
		}

	case modeGameover:
		if g.isKeySpaceJustPressed() {
			g.init()
			g.mode = modeGame
		}
	}

	return nil
}

// Draw method
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	text.Draw(screen, fmt.Sprintf("Hisore: %d", g.hiscore), arcadeFont, 300, 20, color.Black)
	text.Draw(screen, fmt.Sprintf("Score: %d", g.score), arcadeFont, 500, 20, color.Black)
	var xs [3]int
	var ys [3]int

	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(
			"g.y: %d\nTree1 x:%d, y:%d\nTree2 x:%d, y:%d\nTree3 x:%d, y:%d",
			g.dinosaurY,
			xs[0],
			ys[0],
			xs[1],
			ys[1],
			xs[2],
			ys[2],
		))
	}

	g.drawGround(screen)
	g.drawDinosaur(screen)

	switch g.mode {
	case modeTitle:
		text.Draw(screen, "PRESS SPACE KEY", arcadeFont, 245, 240, color.Black)
	case modeLogin:
		//todo 名前を入力してとテキストを入れたい

		// Blink the cursor.
		t := g.text
		if g.counter%60 < 30 {
			t += "_"
		}
		text.Draw(screen, t, arcadeFont, 275, 240, color.Black)
	case modeGame:

	case modeGameover:
		text.Draw(screen, "GAME OVER", arcadeFont, 275, 240, color.Black)
	}
}

func (g *Game) drawDinosaur(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.dinosaurX), float64(g.dinosaurY))
	op.Filter = ebiten.FilterLinear
	if (g.count/5)%2 == 0 {
		screen.DrawImage(dinosaur1Img, op)
		return
	}
	screen.DrawImage(dinosaur2Img, op)
}

func (g *Game) drawGround(screen *ebiten.Image) {
	for i := 0; i < 14; i++ {
		x := float64(groundWidth * i)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, float64(g.ground.y))
		op.GeoM.Translate(float64(g.ground.x), 0.0)
		op.Filter = ebiten.FilterLinear
		screen.DrawImage(groundImg, op)
	}
}

// Layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (g *Game) isKeySpaceJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func (g *Game) isKeyEnterJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}
	return false
}

func main() {
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("Dinosaur Jump")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
