package main

func main() {
	print("Running main...")
	defer print("Finished main.")

	test_play_self()
}

func test_play_self() {
	game := game_from_fen(CHESS_START_POSITION)
	engine_1 := new_engine_minimax_mo_ab()
	engine_2 := new_engine_minimax_mo_ab()
	play_self(&engine_1, &engine_2, game)
}
