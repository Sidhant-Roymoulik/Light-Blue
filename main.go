package main

import (
	"runtime"
)

var timeLeft int64 = 10 * 1000
var increment int64 = 100
var moveTime int64 = NoValue
var movesToGo int16 = 40
var maxDepth uint8 = 100
var maxNodeCount uint64 = 1000000000

func main() {
	// print("Running main...")
	// defer print("Finished main.")

	InitBitboards()
	InitEvalBitboards()
	InitTables()
	InitZobrist()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// print("Version", runtime.Version())
	// print("NumCPU", runtime.NumCPU())
	// print("GOMAXPROCS", runtime.GOMAXPROCS(0))
	// print("Initialization complete.")
	// print()

	// test_play_self()
	// test_play_human()

	// test_benchmark()

	run_uci()
}

func test_play_self() {
	game := game_from_opening("Start Position")
	engine_1 := new_light_blue()
	engine_1.timer.Setup(
		timeLeft,
		increment,
		moveTime,
		movesToGo,
		maxDepth,
		maxNodeCount,
	)
	engine_2 := new_light_blue()
	engine_2.timer.Setup(
		timeLeft,
		increment,
		moveTime,
		movesToGo,
		maxDepth,
		maxNodeCount,
	)
	play_self(&engine_1, &engine_2, game)
}

func test_play_human() {
	game := game_from_opening("Start Position")
	engine_1 := new_light_blue()
	engine_1.timer.Setup(
		timeLeft,
		increment,
		moveTime,
		movesToGo,
		maxDepth,
		maxNodeCount,
	)
	engine_2 := new_engine_human()
	play_human(&engine_1, &engine_2, game)
}

func test_benchmark() {
	engine_1 := new_light_blue()
	engine_1.timer.Setup(
		InfiniteTime,
		NoValue,
		NoValue,
		movesToGo,
		maxDepth,
		maxNodeCount,
	)
	engines := []Engine{&engine_1}
	benchmark_engines(engines, game_from_fen("rn1qkb1r/pp2pppp/5n2/3p1b2/3P4/2N1P3/PP3PPP/R1BQKBNR w KQkq - 0 1").Position())
}

func run_uci() {
	uci_engine := &UCIEngine{}
	uci_engine.loop()
}
