package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"

	"fmt"

	"golang.org/x/image/font/basicfont"
)

const (
	ballSpeed   = 3.0
	paddleSpeed = 6.0
)

var (
	screenWidth  int
	screenHeight int
)

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
}

type Ball struct {
	Object
	dxdt float64 // x velocity per tick
	dydt float64 // y velocity per tick
}

type Game struct {
	playerPaddle   Paddle
	computerPaddle Paddle
	ball           Ball
	playerScore    int
	computerScore  int
	highScore      int
	ballSpeed      float64
}

func main() {
	// Get the screen size for the PC
	monitor := ebiten.Monitor()
	fullScreenWidth, fullScreenHeight := monitor.Size()

	// Set the game screen size
	screenWidth = fullScreenWidth
	screenHeight = fullScreenHeight

	// Optional: Scale the screen size for a windowed mode
	screenWidth /= 2
	screenHeight /= 2

	// Initialize the game components
	// Initialize player paddle
	playerPaddle := Paddle{
		Object: Object{
			X: screenWidth - 40, // Player paddle on the right side
			Y: (screenHeight - 100) / 2,
			W: 15,
			H: 100,
		},
	}

	// Initialize computer paddle
	computerPaddle := Paddle{
		Object: Object{
			X: 25, // Computer paddle on the left side
			Y: (screenHeight - 100) / 2,
			W: 15,
			H: 100,
		},
	}
	ball := Ball{
		Object: Object{
			X: screenWidth / 2,
			Y: screenHeight / 2,
			W: 15,
			H: 15,
		},
		dxdt: ballSpeed,
		dydt: ballSpeed,
	}

	// Create the game instance
	g := &Game{
		playerPaddle:   playerPaddle,
		computerPaddle: computerPaddle,
		ball:           ball,
		ballSpeed:      ballSpeed,
	}

	// Set the window properties
	ebiten.SetWindowTitle("Pong")
	// Enable full-screen mode
	ebiten.SetFullscreen(true)
	// ebiten.SetWindowSize(screenWidth, screenHeight)

	// Run the game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Player paddle
	vector.DrawFilledRect(
		screen,
		float32(g.playerPaddle.X),
		float32(g.playerPaddle.Y),
		float32(g.playerPaddle.W),
		float32(g.playerPaddle.H),
		color.White,
		false,
	)

	// Computer paddle
	vector.DrawFilledRect(
		screen,
		float32(g.computerPaddle.X),
		float32(g.computerPaddle.Y),
		float32(g.computerPaddle.W),
		float32(g.computerPaddle.H),
		color.White,
		false,
	)

	// Ball
	vector.DrawFilledCircle(
		screen,
		float32(g.ball.X)+float32(g.ball.W)/2, // Circle's X center
		float32(g.ball.Y)+float32(g.ball.H)/2, // Circle's Y center
		float32(g.ball.W)/2,                   // Circle's radius
		// float32(g.ball.H),
		color.White,
		false,
	)

	// Player Score
	playerScoreStr := "Player Score: " + fmt.Sprint(g.playerScore)
	text.Draw(screen, playerScoreStr, basicfont.Face7x13, 10, 10, color.White)

	// Computer Score
	computerScoreStr := "Computer Score: " + fmt.Sprint(g.computerScore)
	text.Draw(screen, computerScoreStr, basicfont.Face7x13, 10, 30, color.White)

	// High Score
	highScoreStr := "High Score: " + fmt.Sprint(g.highScore)
	text.Draw(screen, highScoreStr, basicfont.Face7x13, 10, 50, color.White)
}

func (g *Game) Update() error {
	g.playerPaddle.MoveOnKeyPress()

	// Computer paddle AI (slightly slower than the ball)
	g.computerPaddle.MoveToFollowBall(g.ball.Y+g.ball.H/2, int(paddleSpeed-3))

	g.ball.Move()
	g.CollideWithWall()
	g.CollideWithPlayerPaddle()
	g.CollideWithComputerPaddle()
	return nil
}

func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if p.Y+p.H < screenHeight {
			p.Y += int(paddleSpeed)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if p.Y > 0 {
			p.Y -= int(paddleSpeed)
		}
	}
}

func (b *Ball) Move() {
	b.X += int(b.dxdt)
	b.Y += int(b.dydt)
}

func (g *Game) Reset() {
	// Center the ball on the screen
	g.ball.X = (screenWidth - g.ball.W) / 2
	g.ball.Y = (screenHeight - g.ball.H) / 2

	/// Reverse ball direction to the last scorer
	g.ball.dxdt = -g.ball.dxdt

	// Reset the ball speed
	g.ballSpeed = ballSpeed
	g.ball.dxdt = g.ballSpeed
	g.ball.dydt = g.ballSpeed
}

func (g *Game) CollideWithWall() {
	// Ball hits the player's side wall
	if g.ball.X >= screenWidth {
		g.computerScore++
		g.Reset()
	}

	// Ball hits the computer's side wall
	if g.ball.X <= 0 {
		g.playerScore++
		// Update the high score if the player's score is higher
		if g.playerScore > g.highScore {
			g.highScore = g.playerScore
		}
		g.Reset()
	}

	// Ball hits the top or bottom of the screen
	if g.ball.Y <= 0 {
		g.ball.dydt = ballSpeed
	} else if g.ball.Y >= screenHeight {
		g.ball.dydt = -ballSpeed
	}
}

func (g *Game) CollideWithPlayerPaddle() {
	if g.ball.X+g.ball.W >= g.playerPaddle.X && g.ball.Y+g.ball.H >= g.playerPaddle.Y && g.ball.Y <= g.playerPaddle.Y+g.playerPaddle.H {
		g.ball.dxdt = -g.ball.dxdt
	}
}

func (g *Game) CollideWithComputerPaddle() {
	if g.ball.X <= g.computerPaddle.X+g.computerPaddle.W && g.ball.Y+g.ball.H >= g.computerPaddle.Y && g.ball.Y <= g.computerPaddle.Y+g.computerPaddle.H {
		g.ball.dxdt = -g.ball.dxdt

		// Increase the ball's speed
		g.ballSpeed += 1          // Increase speed by a fixed amount (adjust as needed)
		g.ball.dxdt = g.ballSpeed // Apply the new horizontal speed
		g.ball.dydt = g.ballSpeed // Optionally, increase the vertical speed similarly
	}
}

func (p *Paddle) MoveToFollowBall(ballY int, speed int) {
	center := p.Y + p.H/2

	// Move up or down based on the ball's position
	if center < ballY {
		p.Y += int(speed)
	}

	if center > ballY {
		p.Y -= int(speed)
	}

	// Keep paddle within screen bounds
	if p.Y < 0 {
		p.Y = 0
	}

	if p.Y+p.H > screenHeight {
		p.Y = screenHeight - p.H
	}
}
