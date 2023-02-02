package main

import (
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

// Adapted from https://github.com/0hq/antikythera/blob/main/benchmark.go#L24

func benchmark(ply int, engine Engine, pos *chess.Position) float64 {
	engine.setBenchmarkMode(ply)

	print("BEGIN BENCHMARKING -", engine.getName())
	// print("Starting at time", time.Now())

	start := time.Now()
	eval, move := engine.run(pos)
	elapsed := time.Since(start)

	print("Depth:", engine.getDepth())

	// print("Complete at time", time.Now())
	print("Best Move:", move.String())
	if eval > 100000 {
		print("Eval: Mate in", (CHECKMATE_VALUE-eval+1)/2)
	} else if eval < -100000 {
		print("Eval: Mate in", (CHECKMATE_VALUE+eval+1)/2)
	} else {
		print("Eval:", float32(eval*getMultiplier(pos.Turn() == chess.White))/100.0)
	}
	print("Time Taken:", (time.Since(start)).Round(time.Millisecond))
	print("Nodes explored:", engine.getNodesSearched())
	print("Q-Nodes explored:", engine.getQNodesSearched())
	print("Hashes Used:", engine.getHashesUsed())
	print("END BENCHMARKING -")
	print()

	return elapsed.Seconds()
}

func benchmark_range(plymin int, plymax int, engine Engine, pos *chess.Position) {
	engine.reset()
	for i := plymin; i <= plymax; i++ {
		benchmark(i, engine, pos)
	}
}

func benchmark_engines(engines []Engine, pos *chess.Position) {
	for _, engine := range engines {
		benchmark_range(1, 1, engine, pos)
	}
}
