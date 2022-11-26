package main

import (
	"github.com/Sidhant-Roymoulik/chess"
)

type Engine interface {
	getName() string
	run(*chess.Position) (int, *chess.Move)
}

type EngineClass struct {
	name     string
	upgrades EngineUpgrades
}

type EngineUpgrades struct {
	move_ordering       bool
	alphabeta           bool
	iterative_deepening bool
	q_search            bool
	concurrent          bool
}

func (engine *EngineClass) getName() string {
	return engine.name
}
