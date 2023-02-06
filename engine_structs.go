package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

type EngineClass struct {
	name                 string
	author               string
	max_ply              int
	max_q_ply            int
	start                time.Time
	time_limit           time.Duration
	counters             EngineCounters
	upgrades             EngineUpgrades
	tt                   TransTable[SearchEntry]
	age                  uint8        // this is used to age off entries in the transposition table, in the form of a half move clock
	zobristHistory       [1024]uint64 // draw detection history
	zobristHistoryPly    uint16       // draw detection ply
	prev_guess           int          // Used in MTD(f)
	use_mtd_f            bool
	quit_mtd             bool
	killer_moves         [100][2]*chess.Move
	threads              int
	quit_search_at_depth [100]bool
}

type Engine interface {
	getName() string
	getAuthor() string
	getDepth() string
	getNodesSearched() int
	getQNodesSearched() int
	getHashesUsed() int
	setBenchmarkMode(int)
	time_up() bool
	Add_Zobrist_History(uint64)
	Remove_Zobrist_History()
	Is_Draw_By_Repetition(uint64) bool
	resetCounters()
	resetKillerMoves()
	resetZobrist()
	reset()
	run(*chess.Position) (int, *chess.Move)
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
	killer_moves        bool
	lazy_smp            bool
	pvs                 bool
	aspiration_window   bool
}

type EngineCounters struct {
	nodes_searched   int
	q_nodes_searched int
	hashes_used      int
}

type Result struct {
	eval  int
	move  *chess.Move
	depth int
}

func (engine *EngineClass) getName() string {
	return engine.name
}

func (engine *EngineClass) getAuthor() string {
	if engine.author == "" {
		return "Sidhant Roymoulik"
	}
	return engine.author
}

func (engine *EngineClass) getDepth() string {
	return fmt.Sprint(engine.max_ply)
}

func (engine *EngineClass) getNodesSearched() int {
	return engine.counters.nodes_searched
}

func (engine *EngineClass) getQNodesSearched() int {
	return engine.counters.q_nodes_searched
}

func (engine *EngineClass) getHashesUsed() int {
	return engine.counters.hashes_used
}

func (engine *EngineClass) setBenchmarkMode(ply int) {
	// engine.upgrades.lazy_smp = false
	engine.upgrades.iterative_deepening = false
	engine.max_ply = ply
}

func (engine *EngineClass) addKillerMove(move *chess.Move, ply int) {
	if !move.HasTag(chess.Capture) && move != engine.killer_moves[ply][0] {
		engine.killer_moves[ply][1] = engine.killer_moves[ply][0]
		engine.killer_moves[ply][0] = move
	}
}

func (engine *EngineClass) time_up() bool {
	return engine.upgrades.iterative_deepening && time.Since(engine.start) >= engine.time_limit
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

func (engine *EngineClass) resetCounters() {
	engine.counters.nodes_searched = 0
	engine.counters.q_nodes_searched = 0
	engine.counters.hashes_used = 0
}

func (engine *EngineClass) resetKillerMoves() {
	for i := 0; i < len(engine.killer_moves); i++ {
		engine.killer_moves[i][0] = nil
		engine.killer_moves[i][1] = nil
	}
}

func (engine *EngineClass) resetZobrist() {
	engine.zobristHistory = [1024]uint64{}
	engine.zobristHistoryPly = 0
	engine.tt.Clear()
	engine.tt.Resize(64, 16)
}

func (engine *EngineClass) reset() {
	engine.max_ply = 0
	engine.time_limit = TIME_LIMIT
	engine.resetCounters()
	engine.tt = TransTable[SearchEntry]{}
	engine.age = 0
	engine.prev_guess = 0
	engine.use_mtd_f = false
	engine.quit_mtd = false
	engine.resetKillerMoves()
	engine.threads = runtime.GOMAXPROCS(0)
	engine.resetZobrist()
}
