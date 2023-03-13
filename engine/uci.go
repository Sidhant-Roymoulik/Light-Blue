package engine

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
)

type UCIEngine struct {
	engine light_blue
	game   *chess.Game
	moves  uint16

	OpeningBook map[uint64][]PolyglotEntry

	OptionUseBook       bool
	OptionBookPath      string
	OptionBookMoveDelay int
}

func (e *UCIEngine) reset() {
	*e = UCIEngine{}
	e.engine = new_light_blue()
}

func (e *UCIEngine) uci() {
	fmt.Printf("id name %v\n", name)
	fmt.Printf("id author %v\n", author)

	fmt.Printf("\noption name Hash type spin default 64 min 1 max 32000\n")
	fmt.Print("option name Clear Hash type button\n")
	fmt.Print("option name Clear History type button\n")
	fmt.Print("option name Clear Killers type button\n")
	// fmt.Print("option name Clear Counters type button\n")

	fmt.Print("option name UseBook type check default false\n")
	fmt.Print("option name BookPath type string default\n")
	fmt.Print("option name BookMoveDelay type spin default 2 min 0 max 10\n")

	fmt.Print("\nAvailable UCI commands:\n")

	fmt.Print("    * uci\n    * isready\n    * ucinewgame")
	fmt.Print("\n    * setoption name <NAME> value <VALUE>")
	fmt.Print("\n    * position")
	fmt.Print("\n\t* fen <FEN>")
	fmt.Print("\n\t* startpos")
	fmt.Print("\n    * go")
	fmt.Print("\n\t* wtime <MILLISECONDS>\n\t* btime <MILLISECONDS>")
	fmt.Print("\n\t* winc <MILLISECONDS>\n\t* binc <MILLISECONDS>")
	fmt.Print("\n\t* movestogo <INTEGER>\n\t* depth <INTEGER>\n\t* nodes <INTEGER>\n\t* movetime <MILLISECONDS>")
	fmt.Print("\n\t* infinite")

	fmt.Print("\n    * stop\n    * quit\n\n")
	fmt.Printf("uciok\n")
}

func (e *UCIEngine) setOption(command string) {
	fields := strings.Fields(command)
	var option, value string
	parsingWhat := ""

	for _, field := range fields {
		if field == "name" {
			parsingWhat = "name"
		} else if field == "value" {
			parsingWhat = "value"
		} else if parsingWhat == "name" {
			option += field + " "
		} else if parsingWhat == "value" {
			value += field + " "
		}
	}

	option = strings.TrimSuffix(option, " ")
	value = strings.TrimSuffix(value, " ")

	switch option {
	case "Hash":
		size, err := strconv.Atoi(value)
		if err == nil {
			e.engine.uninitializeTT()
			e.engine.resizeTT(uint64(size), SearchEntrySize)
		}
	case "Clear Hash":
		e.engine.clearTT()
	case "Clear History":
		e.engine.resetZobrist()
	case "Clear Killers":
		e.engine.resetKillerMoves()
	case "UseBook":
		if value == "true" {
			e.OptionUseBook = true
		} else if value == "false" {
			e.OptionUseBook = false
		}
	case "BookPath":
		var err error
		e.OpeningBook, err = LoadPolyglotFile(value)

		if err == nil {
			fmt.Println("Opening book loaded...")
		} else {
			fmt.Println("Failed to load opening book...")
		}
	case "BookMoveDelay":
		size, err := strconv.Atoi(value)
		if err == nil {
			e.OptionBookMoveDelay = size
		}
	}
}

func (e *UCIEngine) position(command string) {
	// parse position and moves from args
	// set position in engine
	args := strings.TrimPrefix(command, "position ")
	fen := ""

	if strings.HasPrefix(args, "startpos") {
		args = strings.TrimPrefix(args, "startpos ")
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	} else if strings.HasPrefix(args, "fen ") {
		args = strings.TrimPrefix(args, "startpos ")
		remaining_args := strings.Fields(args)
		fen = strings.Join(remaining_args[0:6], " ")
		args = strings.Join(remaining_args[6:], " ")
	}

	e.game = game_from_fen(fen)
	e.engine.zobristHistoryPly = e.moves
	e.engine.Add_Zobrist_History(Zobrist.GenHash(e.game.Position()))

	if strings.HasPrefix(args, "moves ") {
		args = strings.TrimSuffix(strings.TrimPrefix(args, "moves"), " ")
		if args != "" {
			for _, smove := range strings.Fields(args) {
				move, err := chess.UCINotation{}.Decode(
					e.game.Position(), smove,
				)
				if err != nil {
					panic(err)
				}
				e.game.Move(move)
				e.engine.Add_Zobrist_History(Zobrist.GenHash(e.game.Position()))
			}
		}
	}

	e.moves++
}

func (e *UCIEngine) search(command string) {
	if e.OptionUseBook {
		if entries, ok := e.OpeningBook[GenPolyglotHash(e.game.Position())]; ok {
			// To allow opening variety, randomly select a move from an entry matching
			// the current position.
			entry := entries[rand.Intn(len(entries))]
			move, err := chess.UCINotation{}.Decode(e.game.Position(), entry.Move)
			if err != nil {
				panic(err)
			}

			// if inter.Search.Pos.MoveIsPseduoLegal(move) {
			time.Sleep(time.Duration(e.OptionBookMoveDelay) * time.Second)
			fmt.Printf("bestmove %v\n", move)
			return
			// }
		}
	}

	command = strings.TrimPrefix(command, "go")
	command = strings.TrimPrefix(command, " ")
	fields := strings.Fields(command)

	colorPrefix := "b"
	if e.game.Position().Turn() == chess.White {
		colorPrefix = "w"
	}

	// Parse the go command arguments.
	timeLeft := int(InfiniteTime)
	increment := int(NoValue)
	movesToGo := int(NoValue)
	maxDepth := uint64(MAX_DEPTH)
	maxNodeCount := uint64(math.MaxUint64)
	moveTime := uint64(NoValue)

	for index, field := range fields {
		if strings.HasPrefix(field, colorPrefix) {
			if strings.HasSuffix(field, "time") {
				timeLeft, _ = strconv.Atoi(fields[index+1])
			} else if strings.HasSuffix(field, "inc") {
				increment, _ = strconv.Atoi(fields[index+1])
			}
		} else if field == "movestogo" {
			movesToGo, _ = strconv.Atoi(fields[index+1])
		} else if field == "depth" {
			maxDepth, _ = strconv.ParseUint(fields[index+1], 10, 8)
		} else if field == "nodes" {
			maxNodeCount, _ = strconv.ParseUint(fields[index+1], 10, 64)
		} else if field == "movetime" {
			moveTime, _ = strconv.ParseUint(fields[index+1], 10, 64)
		}
	}

	// Setup the timer with the go command time control information.
	e.engine.timer.Setup(
		int64(timeLeft),
		int64(increment),
		int64(moveTime),
		int16(movesToGo),
		uint8(maxDepth),
		maxNodeCount,
	)

	// Report the best move found by the engine to the GUI.
	_, bestMove := e.engine.run(e.game.Position())
	fmt.Printf("bestmove %v\n", bestMove)
}

func (e *UCIEngine) quit() {
	e.engine.uninitializeTT()
}

func (e *UCIEngine) loop() {
	reader := bufio.NewReader(os.Stdin)

	e.uci()
	e.reset()

	e.engine.resizeTT(DefaultTTSize, SearchEntrySize)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		command = strings.Replace(command, "\r\n", "\n", -1)

		if command == "uci\n" {
			e.uci()
		} else if command == "isready\n" {
			print("readyok")
		} else if strings.HasPrefix(command, "setoption") {
			e.setOption(command)
		} else if strings.HasPrefix(command, "ucinewgame") {
			e.moves = 0
			e.engine.reset()
		} else if strings.HasPrefix(command, "position") {
			e.position(command)
		} else if strings.HasPrefix(command, "go") {
			go e.search(command)
		} else if strings.HasPrefix(command, "stop") {
			e.engine.timer.ForceStop()
		} else if command == "quit\n" {
			e.quit()
			break
		}
	}
}
