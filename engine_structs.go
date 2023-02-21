package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

type EngineClass struct {
	name              string
	author            string
	max_ply           int
	max_q_ply         int
	start             time.Time
	time_limit        time.Duration
	counters          EngineCounters
	upgrades          EngineUpgrades
	tt                TransTable[SearchEntry]
	age               uint8        // this is used to age off entries in the transposition table, in the form of a half move clock
	zobristHistory    [1024]uint64 // draw detection history
	zobristHistoryPly uint16       // draw detection ply
	prev_guess        int          // Used in MTD(f) and aspiration window
	use_mtd_f         bool
	quit_mtd          bool
	killer_moves      [100][2]*chess.Move
	threads           int
	quit_search       bool
}

type Engine interface {
	getName() string
	getAuthor() string
	getDepth() string
	getNodesSearched() int
	getQNodesSearched() int
	getHashesUsed() int
	saveTTPosition(uint64, int, *chess.Move, int, int, uint8)
	probeTTPosition(uint64, int, int, int, int) (int, bool, *chess.Move)
	setBenchmarkMode(int)
	addKillerMove(*chess.Move, int)
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
	hashes_written   int
}

// -----------------------------------------------------------------------------
// 		Principal Variation Stuff
// 		(Adapted from https://github.com/algerbrex/blunder)
// -----------------------------------------------------------------------------

type PVLine struct {
	Moves []*chess.Move
}

// Clear the principal variation line.
func (pvLine *PVLine) clear() {
	pvLine.Moves = nil
}

// Update the principal variation line with a new best move,
// and a new line of best play after the best move.
func (pvLine *PVLine) update(move *chess.Move, newPVLine PVLine) {
	pvLine.clear()
	pvLine.Moves = append(pvLine.Moves, move)
	pvLine.Moves = append(pvLine.Moves, newPVLine.Moves...)
}

// Get the best move from the principal variation line.
func (pvLine *PVLine) getPVMove() *chess.Move {
	if len(pvLine.Moves) == 0 {
		return nil
	}
	return pvLine.Moves[0]
}

func (pvLine PVLine) String() string {
	pv := fmt.Sprintf("%s", pvLine.Moves)
	return pv[1 : len(pv)-1]
}

// -----------------------------------------------------------------------------

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
	return fmt.Sprint(engine.max_ply) + "/" + fmt.Sprint(engine.max_q_ply)
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

func (engine *EngineClass) saveTTPosition(hash uint64, score int, best *chess.Move, ply int, depth int, flag uint8) {
	if !engine.time_up() && best != nil {
		var entry *SearchEntry = engine.tt.Store(hash, depth, engine.age)
		entry.Set(hash, score, best, ply, depth, flag, engine.age)
		engine.counters.hashes_written++
	}
}

func (engine *EngineClass) probeTTPosition(hash uint64, ply int, depth int, alpha int, beta int) (int, bool, *chess.Move) {
	var entry *SearchEntry = engine.tt.Probe(hash)
	var tt_eval, should_use, tt_move = entry.Get(hash, 0, depth, alpha, beta)
	return tt_eval, should_use, tt_move
}

func (engine *EngineClass) setBenchmarkMode(ply int) {
	engine.upgrades.iterative_deepening = false
	engine.max_ply = ply
}

func (engine *EngineClass) addKillerMove(move *chess.Move, ply int) {
	if !move.HasTag(chess.Capture) && move.Promo() == chess.NoPieceType &&
		move != engine.killer_moves[ply][0] {
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
}

func (engine *EngineClass) reset() {
	engine.max_ply = 0
	engine.time_limit = TIME_LIMIT
	engine.resetCounters()
	engine.tt = TransTable[SearchEntry]{}
	engine.age = 0
	engine.prev_guess = 0
	engine.resetKillerMoves()
	engine.threads = runtime.GOMAXPROCS(0)
	engine.resetZobrist()

	engine.tt.Clear()
	engine.tt.Resize(64, 16)
}
