package main

import (
	"runtime"

	"github.com/Sidhant-Roymoulik/chess"
)

func main() {
	print("Running main...")
	defer print("Finished main.")

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	chess.UseNotation(chess.AlgebraicNotation{})

	test_play_self()
	// test_play_human()
}

func test_play_self() {
	game := game_from_fen(CHESS_START_POSITION)
	engine_1 := new_engine_minimax_mo_ab_id_q()
	engine_2 := new_engine_minimax_mo_ab_id_q()
	// engine_1 := new_engine_minimax_mo_ab_q()
	// engine_2 := new_engine_minimax_mo_ab_q()
	play_self(&engine_1, &engine_2, game)
}

func test_play_human() {
	game := game_from_fen(CHESS_START_POSITION)
	engine_1 := new_engine_minimax_cc()
	engine_2 := new_engine_human()
	play_human(&engine_1, &engine_2, game)
}
