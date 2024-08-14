package main

import (
	"fmt"
	"math"

	"github.com/gen2brain/raylib-go/raylib"
)

var SCREEN_WIDTH float32 = 480
var SCREEN_HEIGHT float32 = 720

type Player struct {
	Pos rl.Vector2
	Sz  rl.Vector2
	Spd rl.Vector2
	HP  int8
}

func (p *Player) Move(x float32, y float32) {
	nextPosx := p.Pos.X + x

	if nextPosx < 0 {
		p.Pos.X = 0
	} else if nextPosx+p.Sz.X > SCREEN_WIDTH {
		p.Pos.X = SCREEN_WIDTH - p.Sz.X
	} else {
		p.Pos.X = nextPosx
	}

	nextPosY := p.Pos.Y + y

	if nextPosY < 0 {
		p.Pos.Y = 0
	} else if nextPosY+p.Sz.Y > SCREEN_HEIGHT {
		p.Pos.Y = SCREEN_HEIGHT - p.Sz.Y
	} else {
		p.Pos.Y = nextPosY
	}
}

type Enemy struct {
	Pos rl.Vector2
	Sz  rl.Vector2
	Spd rl.Vector2

	Projectiles []Projectile
	Behaviours  []Behaviour
}

func (e *Enemy) RunBehaviour(frame int32) {
	var bte Behaviour

	for i, value := range e.Behaviours {
		frame -= value.Frames

		if frame < 0 {
			bte = e.Behaviours[i]
			break
		}
	}

	for i, projType := range bte.ProjectileTypes {
		if frame%projType.Reload == 0 {
			e.Projectiles = append(e.Projectiles, bte.ProjectileTypes[i].CreateProjectiles(rl.NewVector2(e.Pos.X+e.Sz.X/2, e.Pos.Y+e.Sz.Y/2))...)
		}
	}

}

func (e *Enemy) UpdateProjectilesPos() {
	for i, proj := range e.Projectiles {
		e.Projectiles[i].Pos.X = proj.Spd.X + proj.Pos.X
		e.Projectiles[i].Pos.Y = proj.Spd.Y + proj.Pos.Y
	}

	for i, proj := range e.Projectiles {
		outOfBounds := CheckOutOfBounds(proj.Pos)
		if outOfBounds {
			e.Projectiles = append(e.Projectiles[:i], e.Projectiles[i+1:]...)
		}
	}
}

type Projectile struct {
	Pos      rl.Vector2
	Sz       rl.Vector2
	Spd      rl.Vector2
	Damage   int8
	Harmless bool
	Color    rl.Color
}

type ProjectileType struct {
	Name                 string
	Reload               int32
	AtOnce               int8
	AngleOffsetIncrement float32
	CurrentAngleOffset   float32
}

func (pt *ProjectileType) CreateProjectiles(pos rl.Vector2) []Projectile {
	projectiles := make([]Projectile, pt.AtOnce)
	var speedMultiplier float32
	var size rl.Vector2
	var damage int8
	var color rl.Color

	switch pt.Name {
	case "small":
		speedMultiplier = 7.0
		size = rl.NewVector2(10, 10)
		damage = 1
		color = rl.Purple
	case "big":
		speedMultiplier = 3.0
		size = rl.NewVector2(30, 30)
		damage = 3
		color = rl.SkyBlue
	default:
		speedMultiplier = 5.0
		size = rl.NewVector2(20, 20)
		damage = 2
		color = rl.Maroon
	}

	var coneAngle float64

	if pt.AtOnce < 5 {
		coneAngle = 60.0
	} else {
		coneAngle = 160.0
	}

	coneAngleRad := coneAngle * (math.Pi / 180)

	for i := 0; i < int(pt.AtOnce); i++ {
		var spd rl.Vector2

		angle := -coneAngleRad/2 + float64(i)*(coneAngleRad/float64(pt.AtOnce-1))

		spd.X = float32(math.Sin(angle+float64(pt.CurrentAngleOffset))) * speedMultiplier
		spd.Y = float32(math.Cos(angle+float64(pt.CurrentAngleOffset))) * speedMultiplier

		projectiles[i] = Projectile{
			Pos:      pos,
			Spd:      spd,
			Sz:       size,
			Damage:   damage,
			Harmless: false,
			Color:    color,
		}
	}

	pt.CurrentAngleOffset = pt.CurrentAngleOffset + pt.AngleOffsetIncrement

	// Rotating cone slightly
	if pt.CurrentAngleOffset > 0.3 {
		pt.AngleOffsetIncrement = -pt.AngleOffsetIncrement
	}

	return projectiles
}

type Behaviour struct {
	Frames          int32
	ProjectileTypes []ProjectileType
}

type Game struct {
	ScreenWidth  int32
	ScreenHeight int32
	Frames       int32

	GameOver bool
	Pause    bool

	Player
	Enemy
}

func (g *Game) Init() {
	g.ScreenHeight = 720
	g.ScreenWidth = 480

	g.Frames = 0
	g.GameOver = false
	g.Pause = false

	g.Player = Player{
		Sz:  rl.NewVector2(20, 20),
		Pos: rl.NewVector2(230, 660),
		Spd: rl.NewVector2(5, 5),
		HP:  10,
	}

	g.Enemy = Enemy{
		Sz:  rl.NewVector2(20, 20),
		Pos: rl.NewVector2(230, 20),
		Spd: rl.NewVector2(10, 10),
		Behaviours: []Behaviour{
			{
				Frames: 60 * 2,
				ProjectileTypes: []ProjectileType{
					{
						Name:                 "big",
						Reload:               60,
						AtOnce:               7,
						AngleOffsetIncrement: 0.08,
						CurrentAngleOffset:   0,
					},
				},
			},
			{
				Frames: 60 * 3,
				ProjectileTypes: []ProjectileType{
					{
						Name:                 "small",
						Reload:               20,
						AtOnce:               10,
						AngleOffsetIncrement: 0.05,
						CurrentAngleOffset:   0,
					},
				},
			},
			{
				Frames: 60 * 3,
				ProjectileTypes: []ProjectileType{
					{
						Reload:               25,
						AtOnce:               7,
						AngleOffsetIncrement: 0.1,
						CurrentAngleOffset:   0,
					},
				},
			},
			{
				Frames: 60 * 3,
				ProjectileTypes: []ProjectileType{
					{
						Name:                 "small",
						Reload:               2,
						AtOnce:               12,
						AngleOffsetIncrement: 0.01,
						CurrentAngleOffset:   0,
					},
				},
			},
			{
				Frames: 60 * 20,
				ProjectileTypes: []ProjectileType{
					{
						Name:                 "big",
						Reload:               40,
						AtOnce:               18,
						AngleOffsetIncrement: 0.025,
						CurrentAngleOffset:   0,
					},
					{
						Name:                 "small",
						Reload:               25,
						AtOnce:               6,
						AngleOffsetIncrement: 0.01,
						CurrentAngleOffset:   0,
					},
				},
			},
		},
	}
}

func (g *Game) CheckCollisions() {
	playerRect := rl.NewRectangle(g.Player.Pos.X, g.Player.Pos.Y, g.Player.Sz.X, g.Player.Sz.Y)

	for i, proj := range g.Enemy.Projectiles {
		if !proj.Harmless && rl.CheckCollisionCircleRec(proj.Pos, proj.Sz.X, playerRect) {
			fmt.Println("HP", g.Player.HP)
			g.Player.HP -= proj.Damage
			g.Enemy.Projectiles[i].Harmless = true
		}
	}
}

func (g *Game) DrawHP() {
	rl.DrawText(
		fmt.Sprintf("HP%d", g.Player.HP),
		10,
		10,
		24,
		rl.White,
	)
}

func (g *Game) CheckGameOver() {
	if g.Player.HP <= 0 {
		g.GameOver = true
	}
}

func CheckOutOfBounds(pos rl.Vector2) bool {
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	if pos.X < 0 || pos.X > float32(screenWidth) || pos.Y < 0 || pos.Y > float32(screenHeight) {
		return true
	}

	return false
}

func (g *Game) Update() {

	if rl.IsKeyPressed(rl.KeyP) {
		g.Pause = !g.Pause
	}

	if g.GameOver || g.Pause {
		return
	}

	if rl.IsKeyDown(rl.KeyD) {
		g.Player.Move(g.Player.Spd.X, 0)
	}
	if rl.IsKeyDown(rl.KeyA) {
		g.Player.Move(-g.Player.Spd.X, 0)
	}
	if rl.IsKeyDown(rl.KeyS) {
		g.Player.Move(0, g.Player.Spd.Y)
	}
	if rl.IsKeyDown(rl.KeyW) {
		g.Player.Move(0, -g.Player.Spd.Y)
	}

	g.Enemy.RunBehaviour(g.Frames)

	g.Enemy.UpdateProjectilesPos()

	g.CheckCollisions()

	g.CheckGameOver()

	g.Frames += 1
}

func (g *Game) Render() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	rl.DrawRectangleV(g.Player.Pos, g.Player.Sz, rl.Magenta)

	rl.DrawRectangleV(g.Enemy.Pos, g.Enemy.Sz, rl.Red)

	for _, proj := range g.Enemy.Projectiles {
		rl.DrawCircleV(proj.Pos, proj.Sz.X, proj.Color)
	}

	g.DrawHP()

	if g.Pause {
		rl.DrawText("GAME PAUSED", g.ScreenWidth/2-rl.MeasureText("GAME PAUSED", 40)/2, g.ScreenHeight/2-40, 40, rl.Gray)
	}

	if g.GameOver {
		rl.DrawText("YOU DIED", g.ScreenWidth/2-rl.MeasureText("YOU DIED", 40)/2, g.ScreenHeight/2-40, 40, rl.Gray)
	}

	rl.EndDrawing()
}

func main() {
	game := Game{}
	game.Init()

	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "TooHoo")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		game.Update()

		game.Render()
	}

	rl.CloseWindow()
}
