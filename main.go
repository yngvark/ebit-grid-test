package main

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yngvark/ebit-grid-test/pkg/game"
	"log"
)

func runGame() error {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Hello, World!")

	g, err := game.NewGame()

	err = ebiten.RunGame(g)
	if err != nil {
		return fmt.Errorf("running game: %w", err)
	}

	return nil
}

func main() {
	err := runGame()
	if err != nil {
		log.Fatal(err)
	}
}
