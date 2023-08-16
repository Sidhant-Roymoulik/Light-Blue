package main

import (
	"runtime"

	"github.com/Sidhant-Roymoulik/Light-Blue/chess"
	"github.com/Sidhant-Roymoulik/Light-Blue/engine"
)

func main() {
	// print("Running main...")
	// defer print("Finished main.")

	chess.InitBitboards()
	engine.InitTables()
	engine.InitEvalBitboards()
	engine.InitZobrist()

	runtime.GOMAXPROCS(runtime.NumCPU())

	engine.RunEngine()
}
