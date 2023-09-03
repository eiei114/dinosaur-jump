package main

import (
	"bytes"
	_ "embed"
	"golang.org/x/net/websocket"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	debug    = false
	screenX  = 640
	screenY  = 640
	fontSize = 10

	// game modes
	modeTitle    = 0
	modeLogin    = 1
	modeGame     = 2
	modeGameOver = 3

	// image sizes
	playerHeight = 100
	playerWidth  = 100
	wallHeight   = 50
	wallWidth    = 50
)

//go:embed resources/images/player.png
var bytePlayerImg []byte

//go:embed resources/images/wall.png
var byteWallImg []byte

var (
	playerImg  *ebiten.Image
	wallImg    *ebiten.Image
	arcadeFont font.Face
)

func init() {
	rand.Seed(time.Now().UnixNano())

	img, _, err := image.Decode(bytes.NewReader(bytePlayerImg))
	if err != nil {
		log.Fatal(err)
	}
	playerImg = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(byteWallImg))
	if err != nil {
		log.Fatal(err)
	}
	wallImg = ebiten.NewImageFromImage(img)

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

type wall struct {
	leftX   float64
	rightX  float64
	topY    float64
	bottomY float64
	size    int // 壁のサイズ、初期値は50で
}

// 関数を更新して、速度に応じて壁を移動させる
func (w *wall) move(speed float64) {
	w.leftX += speed
	w.rightX -= speed
	w.topY += speed
	w.bottomY -= speed
}

// Game struct
type Game struct {
	mode     int
	playerX  int
	playerY  int
	players  []PlayerInfo
	myPlayer PlayerInfo
	wall     *wall // 壁の配列を追加
	runes    []rune
	text     string
	counter  int
	wsConn   *websocket.Conn
}

type PlayerInfo struct {
	x        int
	y        int
	username string
	id       int
	isMine   bool
}

// NewGame method
func NewGame() *Game {
	g := &Game{}
	g.init()
	//g.connectToServer()
	return g
}

// Init method
func (g *Game) init() {
	first := InitPlayer(g.text, true)
	g.myPlayer = first

	// 壁の初期配置
	g.wall = &wall{
		leftX:   0,
		rightX:  screenX - wallWidth,
		topY:    0,
		bottomY: screenY - wallHeight,
		size:    50,
	}
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

		if g.isKeyEnterJustPressed() {

			g.mode = modeGame
		}
	case modeGame:
		//todo 障害物に当たったらゲームオーバーになるようにする

		//todo マルチプラットフォームになるようにメソッド化する
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.myPlayer.y -= 25
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.myPlayer.y += 25

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.myPlayer.x += 25

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.myPlayer.x -= 25
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			//サーバーにリクエストを送るコードをここに書いて
			//go g.sendMessageToServer("A message from the client") // using goroutine so it doesn't block the main loop
		}

		g.wall.move(0.01) // 速度は任意で設定可能

		// Check for collision with wall
		if g.isPlayerCollidingWithWall() {
			g.mode = modeGameOver
		}

		if g.isPlayerCollidingWithOtherPlayers() {
			g.mode = modeGameOver
		}

	case modeGameOver:
		if g.isKeySpaceJustPressed() {
			g.myPlayer.x = 100
			g.myPlayer.y = 100
			g.init()
			g.mode = modeGame
		}
	}

	return nil
}

// Draw method
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	g.drawWall(screen) // 壁を描画
	g.drawPlayer(screen, g.myPlayer)

	for i := 0; i < len(g.players); i++ {
		g.drawPlayer(screen, g.players[i])
	}

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

	case modeGameOver:
		text.Draw(screen, "GAME OVER", arcadeFont, 275, 240, color.Black)
	}
}

func (g *Game) Close() {
	if g.wsConn != nil {
		g.wsConn.Close()
	}
}

func InitPlayer(username string, isMine bool) PlayerInfo {
	var player PlayerInfo
	player.username = username
	player.isMine = isMine
	player.x = 100
	player.y = 100
	return player
}

func (g *Game) drawPlayer(screen *ebiten.Image, player PlayerInfo) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.x), float64(player.y))
	op.ColorM.Scale(0, 0.99, 0.89, 1) // この例では、赤はそのまま、緑は0.99倍、青は0.89倍にスケーリングされます。
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(playerImg, op)
}

func (g *Game) drawWall(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// 左側の壁
	op.GeoM.Reset()
	op.GeoM.Scale(float64(g.wall.size)/float64(wallImg.Bounds().Dx()), float64(screenY)/float64(wallImg.Bounds().Dy()))
	op.GeoM.Translate(g.wall.leftX, 0)
	screen.DrawImage(wallImg, op)

	// 右側の壁
	op.GeoM.Reset()
	op.GeoM.Scale(float64(g.wall.size)/float64(wallImg.Bounds().Dx()), float64(screenY)/float64(wallImg.Bounds().Dy()))
	op.GeoM.Translate(g.wall.rightX, 0)
	screen.DrawImage(wallImg, op)

	// 上側の壁
	op.GeoM.Reset()
	op.GeoM.Scale(float64(screenX)/float64(wallImg.Bounds().Dx()), float64(g.wall.size)/float64(wallImg.Bounds().Dy()))
	op.GeoM.Translate(0, g.wall.topY)
	screen.DrawImage(wallImg, op)

	// 下側の壁
	op.GeoM.Reset()
	op.GeoM.Scale(float64(screenX)/float64(wallImg.Bounds().Dx()), float64(g.wall.size)/float64(wallImg.Bounds().Dy()))
	op.GeoM.Translate(0, g.wall.bottomY)
	screen.DrawImage(wallImg, op)
}

func (g *Game) isPlayerCollidingWithWall() bool {
	playerRight := float64(g.myPlayer.x + playerImg.Bounds().Dx())
	playerBottom := float64(g.myPlayer.y + playerImg.Bounds().Dy())

	// 左側の壁との衝突
	if float64(g.myPlayer.x) < g.wall.leftX+float64(wallWidth) {
		return true
	}
	// 右側の壁との衝突
	if playerRight > g.wall.rightX {
		return true
	}
	// 上側の壁との衝突
	if float64(g.myPlayer.y) < g.wall.topY+float64(wallHeight) {
		return true
	}
	// 下側の壁との衝突
	if playerBottom > g.wall.bottomY {
		return true
	}
	return false
}

func (g *Game) isPlayerCollidingWithOtherPlayers() bool {
	px1, py1, px2, py2 := g.myPlayer.bounds()

	for _, otherPlayer := range g.players {
		if otherPlayer.isMine {
			continue
		}

		ox1, oy1, ox2, oy2 := otherPlayer.bounds()
		if px1 < ox2 && px2 > ox1 && py1 < oy2 && py2 > oy1 {
			return true
		}
	}
	return false
}

func (p *PlayerInfo) bounds() (x1, y1, x2, y2 int) {
	x1 = p.x
	y1 = p.y
	x2 = p.x + playerWidth
	y2 = p.y + playerHeight
	return
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

func (g *Game) connectToServer() {
	serverAddr := "localhost:8080" // またはお使いのサーバーアドレス
	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/ws"}
	conn, err := websocket.Dial(u.String(), "", "http://"+u.Host)
	if err != nil {
		log.Println("サーバーへの接続エラー:", err)
		return
	}

	g.wsConn = conn
}

func (g *Game) sendMessageToServer(message string) {
	if g.wsConn == nil {
		log.Println("WebSocket接続が存在しません。")
		return
	}

	err := websocket.Message.Send(g.wsConn, message)
	if err != nil {
		log.Println("メッセージの送信エラー:", err)
		return
	}

	// サーバからの応答を受け取る (応答が不要な場合はこれをスキップできます)
	var received string
	err = websocket.Message.Receive(g.wsConn, &received)
	if err != nil {
		log.Println("サーバーメッセージの受信エラー:", err)
		return
	}

	log.Println("サーバーからの応答:", received)
}

func main() {
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("Dinosaur Jump")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
