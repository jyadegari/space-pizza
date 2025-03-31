package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

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

	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyUp || ev.Key() == tcell.KeyDown || ev.Key() == tcell.KeyLeft || ev.Key() == tcell.KeyRight {
				s.SetContent(game.Player.X, game.Player.Y+1, ' ', nil, defStyle)
				move(ev.Key(), &game)
				s.SetContent(game.Player.X, game.Player.Y+1, 'X', nil, defStyle)
			} else if ev.Key() == tcell.KeyCtrlR {
				// Handle Ctrl+R
				createGame(&game, width, height-1) // Reserve one line for score
				drawGame(&game, s, defStyle)
			}
		}

		checkFood(&game)
		drawGame(&game, s, defStyle) // Redraw after updating score
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
	World  [][]rune
	Player Player
	Food   []Food
	Width  int
	Height int
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
