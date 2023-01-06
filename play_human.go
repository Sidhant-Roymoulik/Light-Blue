package main

import (
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

func play_human(engine Engine, human Engine, game *chess.Game) {
	engine.reset(game.Position())

	print("Starting Human vs Engine Game", "\n")

	var human_white bool = false
	var white, black Engine

	if human_white {
		white = human
		black = engine
	} else {
		white = engine
		black = human
	}

	print("Player One: " + white.getName())
	print("Player Two: " + black.getName())
	print("")
	print(game.FEN())
	print(game.Position().Board().Draw())

	for game.Outcome() == chess.NoOutcome {
		var eval int
		var move *chess.Move

		if game.Position().Turn() == chess.White {
			eval, move = white.run(game.Position())
		} else {
			eval, move = black.run(game.Position())
		}

		if move == nil {
			panic("No legal moves")
		}

		err := game.Move(move)
		if err != nil {
			panic(err)
		}

		if (game.Position().Turn() == chess.Black && engine == white) || (game.Position().Turn() == chess.White && engine == black) {
			print("Best Move:", move)
			if eval > 100000 {
				print("Eval: Mate in", CHECKMATE_VALUE-eval, "ply")
			} else if eval < -100000 {
				print("Eval: Mate in", CHECKMATE_VALUE+eval, "ply")
			} else {
				print("Eval:", float32(-1*eval*getMultiplier(game.Position().Turn() == chess.White))/100.0)
			}
			print("Time Taken:", (time.Since(start)).Round(time.Millisecond))
			print("Unique Positions Checked:", states)
			print("Q-Positions Checked:", q_states)
			print("Hashes:", hash_hits, hash_writes)
			// print(game.FEN())
		}
		print(game.Position().Board().Draw())
	}
	print(game.Outcome())
	print(game.Method())
	print(game.String())
}
