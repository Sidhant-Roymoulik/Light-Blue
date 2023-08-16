package engine

import (
	"fmt"
	"time"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
)

type Engine struct {
	EngineClass
	max_ply           int
	start             time.Time
	counters          EngineCounters
	timer             TimeManager
	tt                TransTable[SearchEntry]
	age               uint8
	zobristHistory    [1024]uint64
	zobristHistoryPly uint16
	prev_guess        int
	killer_moves      [100][2]*chess.Move
}

type EngineClass struct {
	name     string
	author   string
	upgrades EngineUpgrades
}

type EngineUpgrades struct {
	iterative_deepening bool
}

type EngineCounters struct {
	nodes_searched   uint64
	q_nodes_searched uint64
	hashes_used      uint64
	check_extensions uint64
	smp_pruned       uint64
	nmp_pruned       uint64
	razor_pruned     uint64
	iid_move_found   uint64
	lmp_pruned       uint64
	futility_pruned  uint64
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

func (e *Engine) getName() string {
	return e.name
}

func (e *Engine) printSearchStats() {
	print("Nodes explored:", e.counters.nodes_searched)
	print("Q-Nodes explored:", e.counters.q_nodes_searched)
	print("Hashes Used:", e.counters.hashes_used)
	print("")
	print("Check Extensions:", e.counters.check_extensions)
	print("SMP Prunes:", e.counters.smp_pruned)
	print("NMP Prunes:", e.counters.nmp_pruned)
	print("Razor Prunes:", e.counters.razor_pruned)
	print("Futility Prunes:", e.counters.futility_pruned)
	print("IID Moves Found:", e.counters.iid_move_found)
}

func (e *Engine) setBenchmarkMode(ply int) {
	e.upgrades.iterative_deepening = false
	e.max_ply = ply
	e.timer.MaxDepth = uint8(ply)
}

func (e *Engine) resetCounters() {
	e.counters.nodes_searched = 0
	e.counters.q_nodes_searched = 0
	e.counters.hashes_used = 0
	e.counters.check_extensions = 0
	e.counters.smp_pruned = 0
	e.counters.nmp_pruned = 0
	e.counters.razor_pruned = 0
	e.counters.futility_pruned = 0
	e.counters.iid_move_found = 0
}

func (e *Engine) resizeTT(sizeInMB uint64, entrySize uint64) {
	e.tt.Resize(sizeInMB, entrySize)
}

func (e *Engine) clearTT() {
	e.tt.Clear()
}

func (e *Engine) uninitializeTT() {
	e.tt.Unitialize()
}

func (e *Engine) reset() {
	e.max_ply = 0
	e.age = 0
	e.prev_guess = 0

	e.resetCounters()
	e.resetKillerMoves()
	e.resetZobrist()
	e.tt.Clear()

	e.resizeTT(DefaultTTSize, SearchEntrySize)
}
