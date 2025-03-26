package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
)

func fillArray(arr [][]rune) {
	for i := range arr {
		for j := range arr[i] {
			if rand.Intn(100) < 3 {
				arr[i][j] = '.'
			} else {
				arr[i][j] = ' '
			}
		}
	}

}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	// boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

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

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	width, height := s.Size()

	arr := make([][]rune, height)
	for i := range arr {
		arr[i] = make([]rune, width)
	}

	fillArray(arr)
	for y, row := range arr {
		for x, ch := range row {
			s.SetContent(x, y, ch, nil, defStyle)
		}
	}

	playerX := width / 2
	playerY := height / 2
	s.SetContent(playerX, playerY, 'X', nil, defStyle)

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
				s.SetContent(playerX, playerY, ' ', nil, defStyle)
				move(ev.Key(), &playerX, &playerY, width, height)
				s.SetContent(playerX, playerY, 'X', nil, defStyle)
			}

		}
	}
}

func move(key tcell.Key, playerX *int, playerY *int, width, height int) {
    switch key { 
    case tcell.KeyUp:
        if *playerY > 0 {
            *playerY--
        }
    case tcell.KeyDown:
        if *playerY < height-1 {
            *playerY++
        }
    case tcell.KeyLeft:
        if *playerX > 0 {
            *playerX--
        }
    case tcell.KeyRight:
        if *playerX < width-1 {
            *playerX++
        }
    default:
        // Do nothing for other keys
        
    }
}
