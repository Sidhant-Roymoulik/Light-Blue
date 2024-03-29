package engine

import (
	"time"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
)

func play_self(white *Engine, black *Engine, game *chess.Game) {
	white.reset()
	black.reset()

	print("Starting Engine vs Engine Game", "\n")

	print("White Player: " + white.getName())
	print("Black Player: " + black.getName())
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

		var engine *Engine = white
		if game.Position().Turn() == chess.White {
			engine = black
		}

		print("")
		print("Depth:", engine.max_ply)
		print("Best Move:", move.String())
		if eval > MATE_CUTOFF {
			print("Eval: Mate in", (CHECKMATE_VALUE-eval+1)/2)
		} else if eval < -MATE_CUTOFF {
			print("Eval: Mate in", (CHECKMATE_VALUE+eval+1)/2)
		} else {
			print("Eval:", float32(-1*eval*getMultiplier(game.Position().Turn() == chess.White))/100.0)
		}

		print("Time Taken:", (time.Since(start)).Round(time.Millisecond))
		engine.printSearchStats()
		// print(game.FEN())
		print(game.Position().Board().Draw())
	}
	print(game.Outcome())
	print(game.Method())
	print(game.String())
}
