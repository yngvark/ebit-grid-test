package main

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yngvark/ebit-grid-test/pkg/game"
	"os"
)

func runGame() error {
	// Logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Game
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
		log.Error().Err(err).Msg("running game failed")
	}
}
