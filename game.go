package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
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
	createGame(&game, width, height)
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
				s.SetContent(game.Player.X, game.Player.Y, ' ', nil, defStyle)
				move(ev.Key(), &game)
				s.SetContent(game.Player.X, game.Player.Y, 'X', nil, defStyle)
            } else if ev.Key() == tcell.KeyCtrlR {
                // Handle Ctrl+R
                createGame(&game, width, height)
                drawGame(&game, s, defStyle)
            }
		}
	}
}

func drawGame(game *Game, s tcell.Screen, defStyle tcell.Style) {
	for y, row := range game.World {
		for x, ch := range row {
			s.SetContent(x, y, ch, nil, defStyle)
		}
	}

    for _, food := range game.Food {
        s.SetContent(food.X, food.Y, 'o', nil, defStyle)
        if game.World[food.Y][food.X] == '.' {
            game.World[food.Y][food.X] = ' '
        }
    }
	s.SetContent(game.Player.X, game.Player.Y, 'X', nil, defStyle)
}



type Food struct {
    X int
    Y int
    Duration int
}

type Player struct {
    X int
    Y int
    Score int
}

type Game struct {
    World [][]rune
    Player Player
    Food []Food
    Width int
    Height int
}
func createGame(game *Game, width, height int) {

    player := Player{
        X: rand.Intn(width),
        Y: rand.Intn(height),
        Score: 0,
    }

	world := make([][]rune, height)
	for i := range world {
		world[i] = make([]rune, width)
	}

    foods := []Food{}
    for i := 0; i < 10; i++ { 
        food := Food{
            X: rand.Intn(width),
            Y: rand.Intn(height),
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
