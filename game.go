package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// InputEvent represents a game input event
type InputEvent struct {
	Key tcell.Key
}

func main() {
	defStyle, s := createScreen()
	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	width, height := s.Size()
	var game Game
	createGame(&game, width, height-1) // Reserve one line for score
	drawGame(&game, s, defStyle)

	// Create channel for input events
	inputChan := make(chan InputEvent, 10)

	// Start input handling goroutine
	go handleInput(s, inputChan)

	// Game loop
	for {
		// Process input events (non-blocking)
		select {
		case ev := <-inputChan:
			switch ev.Key {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
				s.SetContent(game.Player.X, game.Player.Y+1, ' ', nil, defStyle)
				move(ev.Key, &game)
				s.SetContent(game.Player.X, game.Player.Y+1, 'X', nil, defStyle)
			case tcell.KeyCtrlR:
				createGame(&game, width, height-1)
				drawGame(&game, s, defStyle)
			}
		default:
			// No input event, continue with game loop
		}

		// Update game state
		checkFood(&game)
		moveEnemies(&game)

		// Render
		drawGame(&game, s, defStyle)
		s.Show()

		// Small sleep to prevent CPU hogging
		time.Sleep(50 * time.Millisecond)
	}
}

// handleInput processes input events in a separate goroutine
func handleInput(s tcell.Screen, inputChan chan<- InputEvent) {
	for {
		ev := s.PollEvent()
		if ev == nil {
			continue
		}

		switch ev := ev.(type) {
		case *tcell.EventKey:
			inputChan <- InputEvent{Key: ev.Key()}
		}
	}
}

func drawGame(game *Game, s tcell.Screen, defStyle tcell.Style) {
	// Clear screen first
	s.Clear()

	// Draw score at the top
	scoreText := []rune(fmt.Sprintf("Score: %d", game.Player.Score))
	for i, ch := range scoreText {
		s.SetContent(i, 0, ch, nil, defStyle)
	}

	// Draw world with offset to leave room for score
	for y, row := range game.World {
		for x, ch := range row {
			s.SetContent(x, y+1, ch, nil, defStyle)
		}
	}

	for _, food := range game.Food {
		s.SetContent(food.X, food.Y+1, 'o', nil, defStyle)
		if game.World[food.Y][food.X] == '.' {
			game.World[food.Y][food.X] = ' '
		}
	}

	// Draw enemies
	for _, enemy := range game.Enemies {
		s.SetContent(enemy.X, enemy.Y+1, 'E', nil, defStyle)
	}

	s.SetContent(game.Player.X, game.Player.Y+1, 'X', nil, defStyle)
}

type Food struct {
	X        int
	Y        int
	Duration int
}

type Player struct {
	X     int
	Y     int
	Score int
}

type Game struct {
	World   [][]rune
	Player  Player
	Food    []Food
	Enemies []Enemy
	Width   int
	Height  int
}

type Enemy struct {
	X           int
	Y           int
	Direction   int // 0: up, 1: right, 2: down, 3: left
	MoveCounter int // For controlling enemy movement speed
}

func createGame(game *Game, width, height int) {
	player := Player{
		X:     rand.Intn(width),
		Y:     rand.Intn(height),
		Score: 0,
	}

	world := make([][]rune, height)
	for i := range world {
		world[i] = make([]rune, width)
	}

	foods := []Food{}
	for i := 0; i < 10; i++ {
		food := Food{
			X:        rand.Intn(width),
			Y:        rand.Intn(height),
			Duration: 10,
		}
		foods = append(foods, food)
	}

	// Create enemies
	enemies := []Enemy{}
	for i := 0; i < 3; i++ {
		enemy := Enemy{
			X:           rand.Intn(width),
			Y:           rand.Intn(height),
			Direction:   rand.Intn(4), // Random initial direction
			MoveCounter: 0,
		}
		enemies = append(enemies, enemy)
	}

	for i := range world {
		for j := range world[i] {
			if i == player.Y && j == player.X {
				continue
			}
			if rand.Intn(100) < 3 {
				world[i][j] = '.'
			} else {
				world[i][j] = ' '
			}
		}
	}

	game.World = world
	game.Player = player
	game.Food = foods
	game.Enemies = enemies
	game.Width = width
	game.Height = height
}

func createScreen() (tcell.Style, tcell.Screen) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	return defStyle, s
}

func move(key tcell.Key, game *Game) {

	switch key {
	case tcell.KeyUp:
		if game.Player.Y > 0 && game.World[game.Player.Y-1][game.Player.X] != '.' {
			game.Player.Y--
		}
	case tcell.KeyDown:
		if game.Player.Y < game.Height-1 && game.World[game.Player.Y+1][game.Player.X] != '.' {
			game.Player.Y++
		}
	case tcell.KeyLeft:
		if game.Player.X > 0 && game.World[game.Player.Y][game.Player.X-1] != '.' {
			game.Player.X--
		}
	case tcell.KeyRight:
		if game.Player.X < game.Width-1 && game.World[game.Player.Y][game.Player.X+1] != '.' {
			game.Player.X++
		}
	default:
		// Do nothing for other keys
	}
}

func checkFood(game *Game) {
	var remainingFood []Food

	for _, food := range game.Food {
		if game.Player.X == food.X && game.Player.Y == food.Y {
			// Player found food, increase score
			game.Player.Score += 10
		} else {
			// Keep the food in the list
			remainingFood = append(remainingFood, food)
		}
	}

	// Update the food list to only include remaining food
	game.Food = remainingFood
}

func moveEnemies(game *Game) {
	for i := range game.Enemies {
		enemy := &game.Enemies[i]

		// Only move every few frames to slow down enemies
		enemy.MoveCounter++
		if enemy.MoveCounter < 3 {
			continue
		}

		enemy.MoveCounter = 0

		// Randomly change direction occasionally
		if rand.Intn(10) < 2 {
			enemy.Direction = rand.Intn(4)
		}

		// Try to move in the current direction
		newX, newY := enemy.X, enemy.Y

		switch enemy.Direction {
		case 0: // up
			if enemy.Y > 0 && game.World[enemy.Y-1][enemy.X] != '.' {
				newY--
			}
		case 1: // right
			if enemy.X < game.Width-1 && game.World[enemy.Y][enemy.X+1] != '.' {
				newX++
			}
		case 2: // down
			if enemy.Y < game.Height-1 && game.World[enemy.Y+1][enemy.X] != '.' {
				newY++
			}
		case 3: // left
			if enemy.X > 0 && game.World[enemy.Y][enemy.X-1] != '.' {
				newX--
			}
		}

		// If we couldn't move in that direction, try a different one
		if newX == enemy.X && newY == enemy.Y {
			enemy.Direction = rand.Intn(4)
		} else {
			enemy.X, enemy.Y = newX, newY
		}
	}
}
