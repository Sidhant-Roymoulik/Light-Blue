package engine

import (
	"fmt"
	"time"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
	"github.com/bndr/gotabulate"
)

// Adapted from https://github.com/0hq/antikythera/blob/main/benchmark.go#L24

func benchmark(
	ply int, e *Engine, pos *chess.Position, rows [][]interface{},
) [][]interface{} {
	e.setBenchmarkMode(ply)
	e.resetZobrist()

	print("BEGIN BENCHMARKING -", e.getName())
	print("Benchmark Depth:", e.max_ply)

	start := time.Now()
	eval, move := e.run(pos)

	row := []interface{}{
		e.max_ply,
		move.String(),
		getMateOrCPScore(eval),
		(time.Since(start)).Round(time.Millisecond),
		e.counters.nodes_searched,
		e.counters.q_nodes_searched,
		e.counters.hashes_used,
		e.counters.check_extensions,
		e.counters.smp_pruned,
		e.counters.nmp_pruned,
		e.counters.razor_pruned,
		e.counters.iid_move_found,
		e.counters.lmp_pruned,
		e.counters.futility_pruned,
	}
	rows = append(rows, row)

	print("END BENCHMARKING -")
	print()

	return rows
}

func benchmark_range(plymin int, plymax int, e *Engine, pos *chess.Position) {
	rows := [][]interface{}{}
	for i := plymin; i <= plymax; i++ {
		rows = benchmark(i, e, pos, rows)

		t := gotabulate.Create(rows)
		t.SetHeaders([]string{
			"Depth",
			"Move",
			"Eval",
			"Time",
			"Nodes",
			"Q-Nodes",
			"Hashes",
			"Checks",
			"SMP",
			"NMP",
			"Razor",
			"IID",
			"LMP",
			"Futility",
		})
		fmt.Println(t.Render("grid"))
	}
}

func benchmark_engines(engines []*Engine, pos *chess.Position) {
	for _, e := range engines {
		e.reset()
		benchmark_range(1, 9, e, pos)
	}
}
