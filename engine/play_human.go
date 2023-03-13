package engine

import (
	"time"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
)

func play_human(e Engine, human Engine, game *chess.Game) {
	e.reset()

	print("Starting Human vs Engine Game", "\n")

	var human_white bool = false
	var white, black Engine

	if human_white {
		white = human
		black = e
	} else {
		white = e
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
		var start = time.Now()

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

		if (game.Position().Turn() == chess.Black && e == white) || (game.Position().Turn() == chess.White && e == black) {
			print("Depth:", e.getDepth())
			print("Best Move:", move.String())
			if eval > 100000 {
				print("Eval: Mate in", (CHECKMATE_VALUE-eval+1)/2)
			} else if eval < -100000 {
				print("Eval: Mate in", (CHECKMATE_VALUE+eval+1)/2)
			} else {
				print("Eval:", float32(-1*eval*getMultiplier(game.Position().Turn() == chess.White))/100.0)
			}
			print("Time Taken:", (time.Since(start)).Round(time.Millisecond))
			print("Nodes explored:", e.getNodesSearched())
			print("Q-Nodes explored:", e.getQNodesSearched())
			print("Hashes Used:", e.getHashesUsed())
			// print(game.FEN())
		}
		print(game.Position().Board().Draw())
	}
	print(game.Outcome())
	print(game.Method())
	print(game.String())
}
