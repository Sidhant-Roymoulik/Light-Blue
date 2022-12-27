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
}

func test_play_self() {
	game := game_from_fen(CHESS_START_POSITION)
	// engine_1 := new_engine_minimax_mo_ab_q_id()
	// engine_2 := new_engine_minimax_mo_ab_q_id()
	engine_1 := new_engine_version_1_1()
	engine_2 := new_engine_version_1_1()
	play_self(&engine_1, &engine_2, game)
}

func test_play_human() {
	game := game_from_fen(CHESS_START_POSITION)
	engine_1 := new_engine_version_1_1()
	engine_2 := new_engine_human()
	play_human(&engine_1, &engine_2, game)
}
