package main

import (
	"runtime"
)

func main() {
	print("Running main...")
	defer print("Finished main.")

	InitZobrist()

	runtime.GOMAXPROCS(runtime.NumCPU())

	print("Version", runtime.Version())
	print("NumCPU", runtime.NumCPU())
	print("GOMAXPROCS", runtime.GOMAXPROCS(0))
	print("Initialization complete.")
	print()

	test_play_self()
	// test_play_human()

	// test_benchmark()
}

func test_play_self() {
	game := game_from_opening("Sicilian Defense")
	engine_1 := new_light_blue_1_0()
	engine_2 := new_light_blue_1_0()
	play_self(&engine_1, &engine_2, game)
}

func test_play_human() {
	game := game_from_opening("Start Position")
	engine_1 := new_engine_version_4_0()
	engine_2 := new_engine_human()
	play_human(&engine_1, &engine_2, game)
}

func test_benchmark() {
	engine1 := new_light_blue_1_0()
	engines := []Engine{&engine1}
	benchmark_engines(engines, game_from_fen("rn1qkb1r/pp2pppp/5n2/3p1b2/3P4/2N1P3/PP3PPP/R1BQKBNR w KQkq - 0 1").Position())
}
