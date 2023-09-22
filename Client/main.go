package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"
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
	mode               int
	playerX            int
	playerY            int
	players            []PlayerInfo
	myPlayer           PlayerInfo
	wall               *wall // 壁の配列を追加
	runes              []rune
	text               string
	counter            int
	npcs               []PlayerInfo
	speedMultiplier    float64
	maxSpeedMultiplier float64
	timePassed         float64 // 経過時間（秒）
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
	g := &Game{
		maxSpeedMultiplier: 10.0, // この値は任意で設定できます。例として3倍速とします。
	}
	g.init()
	return g
}

// Init method
func (g *Game) init() {
	first := InitPlayer(g.text, true)
	g.myPlayer = first
	g.timePassed = 0

	// 壁の初期配置
	g.wall = &wall{
		leftX:   0,
		rightX:  screenX - wallWidth,
		topY:    0,
		bottomY: screenY - wallHeight,
		size:    50,
	}

	// NPCの初期化: 前回のNPCをクリアしてから新しいNPCを追加
	g.npcs = []PlayerInfo{}
	g.npcs = append(g.npcs, InitNPC("NPC1"))
	g.npcs = append(g.npcs, InitNPC("NPC2"))
	g.npcs = append(g.npcs, InitNPC("NPC3"))

	g.speedMultiplier = 1.0 // 初期の乗数
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

		g.timePassed += 1 / 60.0 // 60FPSを仮定

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

		//g.wall.move(0.01) // 速度は任意で設定可能

		// Check for collision with wall
		if g.isPlayerCollidingWithWall() {
			g.mode = modeGameOver
		}

		if g.isPlayerCollidingWithOtherPlayers() {
			g.mode = modeGameOver
		}

		fmt.Println("プレイヤー情報", g.myPlayer)

		// NPCsを更新
		for i := range g.npcs {
			g.moveNPC(&g.npcs[i])
		}

		// Check for collision with NPCs
		if g.isPlayerCollidingWithNPCs() {
			g.mode = modeGameOver
		}

		g.speedMultiplier += 0.001 // この値は微調整する必要があります。
		if g.speedMultiplier > g.maxSpeedMultiplier {
			g.speedMultiplier = g.maxSpeedMultiplier
		}

	case modeGameOver:
		if g.isKeySpaceJustPressed() {
			g.init()
			g.mode = modeGame
		}
	}

	return nil
}

func InitNPC(name string) PlayerInfo {
	var npc PlayerInfo
	npc.username = name
	// Adjusting the initial position of the NPC considering the collision offset
	npc.x = rand.Intn(screenX-playerWidth*2) + playerWidth/2
	npc.y = rand.Intn(screenY-playerHeight*2) + playerHeight/2
	return npc
}

func (g *Game) isPlayerCollidingWithNPCs() bool {
	px1, py1, px2, py2 := g.myPlayer.bounds()

	collisionOffset := 50 // Adjust this value to increase or decrease the collision offset

	for _, npc := range g.npcs {
		nx1, ny1, nx2, ny2 := npc.x+collisionOffset, npc.y+collisionOffset, npc.x+playerWidth-collisionOffset, npc.y+playerHeight-collisionOffset
		if px1 < nx2 && px2 > nx1 && py1 < ny2 && py2 > ny1 {
			return true
		}
	}
	return false
}

func (g *Game) moveNPC(npc *PlayerInfo) {
	direction := rand.Intn(4)             // 0:上, 1:下, 2:左, 3:右
	moveAmount := 5.0 * g.speedMultiplier // 乗数を考慮して移動量を計算

	switch direction {
	case 0:
		npc.y -= int(moveAmount)
	case 1:
		npc.y += int(moveAmount)
	case 2:
		npc.x -= int(moveAmount)
	case 3:
		npc.x += int(moveAmount)
	}

	// NPCがスクリーンから出ないようにする
	if npc.x < 0 {
		npc.x = 0
	}
	if npc.y < 0 {
		npc.y = 0
	}
	if npc.x > screenX-playerWidth {
		npc.x = screenX - playerWidth
	}
	if npc.y > screenY-playerHeight {
		npc.y = screenY - playerHeight
	}
}

// Draw method
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	g.drawWall(screen) // 壁を描画
	g.drawPlayer(screen, g.myPlayer)

	for i := 0; i < len(g.players); i++ {
		g.drawPlayer(screen, g.players[i])
	}

	for _, npc := range g.npcs {
		g.drawNpcPlayer(screen, npc)
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
		timeText := fmt.Sprintf("%.1f", g.timePassed)
		text.Draw(screen, timeText, arcadeFont, screenX-50, 20, color.White)
	case modeGameOver:
		screen.Fill(color.White) // Clear the screen

		users, err := getUserData()
		if err == nil {
			// リストの開始位置を定義
			yPosition := 260

			// 配列内の各ユーザー情報を表示
			for _, user := range users {
				text.Draw(screen, fmt.Sprintf("Name: %s", user.Name), arcadeFont, 275, yPosition, color.Black)
				yPosition += 20 // 次の行の位置に移動
				text.Draw(screen, fmt.Sprintf("HighScore: %d", user.HighScore), arcadeFont, 275, yPosition, color.Black)
				yPosition += 20 // 次の行の位置に移動
			}
		}
		text.Draw(screen, "GAME OVER", arcadeFont, 275, 240, color.Black)
	}
}

func (g *Game) Close() {

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

func (g *Game) drawNpcPlayer(screen *ebiten.Image, player PlayerInfo) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.x), float64(player.y))
	op.ColorM.Scale(1, 0, 0, 1)
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

func main() {
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("Dinosaur Jump")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

type User struct {
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
}

func getUserData() ([]User, error) {
	resp, err := http.Get("http://localhost:8080/users/get")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}
