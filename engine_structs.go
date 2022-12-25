package main

import (
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

type Engine interface {
	getName() string
	run(*chess.Position) (int, *chess.Move)
}

type EngineClass struct {
	name       string
	max_ply    int
	start      time.Time
	time_limit time.Duration
	upgrades   EngineUpgrades
}

type EngineUpgrades struct {
	concurrent          bool
	move_ordering       bool
	alphabeta           bool
	q_search            bool
	iterative_deepening bool
	transposition_table bool
	mtd_f               bool
	lazy_smp            bool
}

func (engine *EngineClass) getName() string {
	return engine.name
}
