package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/Sidhant-Roymoulik/chess"
)

type UCIEngine struct {
	reader  *bufio.Reader
	writer  *bufio.Writer
	name    string
	engine  Engine
	options map[string]Option
}

type Option struct {
	name         string
	typ          string
	defaultValue string
	minValue     string
	maxValue     string
}

func NewUCIEngine(name string) *UCIEngine {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	engine := new_light_blue()
	options := make(map[string]Option)
	return &UCIEngine{reader, writer, engine.getName(), &engine, options}
}

func (e *UCIEngine) AddOption(name string, typ string, defaultValue string, minValue string, maxValue string) {
	e.options[name] = Option{name, typ, defaultValue, minValue, maxValue}
}

func (e *UCIEngine) Run() {
	for {
		line, _, err := e.reader.ReadLine()
		if err != nil {
			break
		}
		command := strings.Split(string(line), " ")
		switch command[0] {
		case "uci":
			e.UCI()
		case "isready":
			e.IsReady()
		case "setoption":
			e.SetOption(command[1:])
		case "ucinewgame":
			e.UCINewGame()
		case "position":
			e.Position(command[1:])
		case "go":
			e.Go(command[1:])
		case "stop":
			e.Stop()
		case "ponderhit":
			e.PonderHit()
		case "quit":
			e.Quit()
			return
		}
	}
}

func (e *UCIEngine) UCI() {
	e.writer.WriteString("id name " + e.name + "\n")
	for _, option := range e.options {
		e.writer.WriteString("option name " + option.name + " type " + option.typ)
		if option.defaultValue != "" {
			e.writer.WriteString(" default " + option.defaultValue)
		}
		if option.minValue != "" {
			e.writer.WriteString(" min " + option.minValue)
		}
		if option.maxValue != "" {
			e.writer.WriteString(" max " + option.maxValue)
		}
		e.writer.WriteString("\n")
	}
	e.writer.WriteString("uciok\n")
	e.writer.Flush()
}

func (e *UCIEngine) IsReady() {
	e.writer.WriteString("readyok\n")
	e.writer.Flush()
}

func (e *UCIEngine) SetOption(args []string) {
	// parse option name and value from args
	// set option value in engine
	if len(args) < 4 {
		return
	}
	var name, value = args[1], args[3]
	for _, option := range e.options {
		if strings.EqualFold(option.name, name) {
			option.defaultValue = value
		}
	}
}

func (e *UCIEngine) UCINewGame() {
	// reset engine's internal state for a new game
	e.engine.reset()
}

func (e *UCIEngine) Position(args []string) {
	// parse position and moves from args
	// set position in engine
	var token = args[0]
	var fen string
	var movesIndex = findIndexString(args, "moves")
	if token == "startpos" {
		fen = game_from_opening("Start Position").FEN()
	} else if token == "fen" {
		if movesIndex == -1 {
			fen = strings.Join(args[1:], " ")
		} else {
			fen = strings.Join(args[1:movesIndex], " ")
		}
	} else {
		return
	}
	var position = game_from_fen(fen).Position()
	// e.engine.Add_Zobrist_History(Zobrist.GenHash(position))
	if movesIndex >= 0 && movesIndex+1 < len(args) {
		for _, smove := range args[movesIndex+1:] {
			move, err := chess.AlgebraicNotation{}.Decode(position, smove)
			if err != nil {
				panic(err)
			}
			position = position.Update(move)
			e.engine.Add_Zobrist_History(Zobrist.GenHash(position))
		}
	}
}

func (e *UCIEngine) Go(args []string) {
	// parse search parameters from args
	// start search in engine
	// send best move to GUI
}

func (e *UCIEngine) Stop() {
	// stop current search
}

func (e *UCIEngine) PonderHit() {
	// ponderhit command is used to inform the engine that the opponent has played the expected move
	// the engine should continue searching
}

func (e *UCIEngine) Quit() {
	// clean up and exit
}

func findIndexString(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

// func main() {
// 	engine := NewUCIEngine("MyEngine")
// 	engine.AddOption("Hash", "spin", "32", "1", "1024")
// 	engine.Run()
// }
