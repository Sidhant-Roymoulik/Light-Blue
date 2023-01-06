package main

import (
	"time"

	"github.com/Sidhant-Roymoulik/chess"
)

func play_self(white Engine, black Engine, game *chess.Game) {
	white.reset(game.Position())
	black.reset(game.Position())

	print("Starting Engine vs Engine Game", "\n")

	print("White Player: " + white.getName())
	print("Black Player: " + black.getName())
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

		print("Best Move:", move.String())
		if eval > 100000 {
			print("Eval: Mate in", (CHECKMATE_VALUE-eval+1)/2)
		} else if eval < -100000 {
			print("Eval: Mate in", (CHECKMATE_VALUE+eval+1)/2)
		} else {
			print("Eval:", float32(-1*eval*getMultiplier(game.Position().Turn() == chess.White))/100.0)
		}
		print("Time Taken:", (time.Since(start)).Round(time.Millisecond))
		print("Unique Positions Checked:", states)
		print("Q-Positions Checked:", q_states)
		print("Hashes Used:", hash_hits)
		print("Total:", states+q_states-hash_hits)
		// print(game.FEN())
		print(game.Position().Board().Draw())
	}
	print(game.Outcome())
	print(game.Method())
	print(game.String())
}
