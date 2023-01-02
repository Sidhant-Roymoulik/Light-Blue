package main

import (
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

type EngineClass struct {
	name              string
	max_ply           int
	start             time.Time
	time_limit        time.Duration
	upgrades          EngineUpgrades
	tt                TransTable[SearchEntry]
	age               uint8        // this is used to age off entries in the transposition table, in the form of a half move clock
	zobristHistory    [1024]uint64 // draw detection history
	zobristHistoryPly uint16       // draw detection ply
	prev_guess        int          // used to decide between mtd(f) and mtd(bi)
	use_mtd_f         bool
}

type Engine interface {
	getName() string
	time_up() bool
	run(*chess.Position) (int, *chess.Move)
	reset_TT(position *chess.Position)
	Add_Zobrist_History(hash uint64)
	Remove_Zobrist_History()
	Is_Draw_By_Repetition(hash uint64) bool
}

type EngineUpgrades struct {
	concurrent          bool
	move_ordering       bool
	alphabeta           bool
	q_search            bool
	delta_pruning       bool
	iterative_deepening bool
	transposition_table bool
	mtd                 bool
	lazy_smp            bool
}

func (engine *EngineClass) getName() string {
	return engine.name
}

func (engine *EngineClass) time_up() bool {
	return time.Since(engine.start) > engine.time_limit
}

func (engine *EngineClass) reset_TT(position *chess.Position) {
	engine.tt.Clear()
	engine.tt.Resize(64, 16)
	engine.zobristHistory[engine.zobristHistoryPly] = Zobrist.GenHash(position)
}

// adds to zobrist history, which is used for draw detection
func (engine *EngineClass) Add_Zobrist_History(hash uint64) {
	engine.zobristHistoryPly++
	engine.zobristHistory[engine.zobristHistoryPly] = hash
}

// decrements ply counter, which means history will be overwritten
func (engine *EngineClass) Remove_Zobrist_History() {
	engine.zobristHistoryPly--
}

func (engine *EngineClass) Is_Draw_By_Repetition(hash uint64) bool {
	for i := uint16(0); i < engine.zobristHistoryPly; i++ {
		if engine.zobristHistory[i] == hash {
			return true
		}
	}
	return false
}
