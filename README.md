# Light Blue
First attempt at a UCI-Compatible Chess Engine Built in Golang

---

## Current Version Upgrades

### Basic Requirements

 - [x] DFS
 - [x] Alpha-Beta Pruning
 - [x] Time Management
 - [x] UCI Protocol

### S-Tier Upgrades

 - [x] Piece Square Tables
 - [x] Quiescence Search
 - [x] Transposition/Hash Tables

### A-Tier Upgrades

 - [x] Iterative Deepening
 - [x] Move Picking
 - [x] Principal-Variation + Null-Window Search
 - [x] Aspiration Window
 - [x] Null Move Pruning
 - [x] Openings
 - [ ] Tablebases

### B-Tier Upgrades

 - [ ] Lazy SMP (Parallelize Engine using Golang)
 - [x] Check Extension
 - [x] Static Move Pruning
 - [x] Razoring
 - [x] Extended Futility Pruning
 - [ ] Internal Iterative Deepening
 - [x] Late Move Pruning
 - [ ] Late Move Reduction
 - [ ] Singular Extensions

### C-Tier Upgrades

 - [x] Killer Moves Heuristic

---

Challenged by Will Depue ([0hq](https://github.com/0hq)).

Check out his [Tutorial](https://www.chessengines.org/) and [Starter Code](https://github.com/0hq/starter_chess_engine) for Javascript and his [Engine in Golang](https://github.com/0hq/antikythera/tree/main).

---

### Resources

[Chessprogramming wiki](https://www.chessprogramming.org/Main_Page) is the chess engine dev's Library of Alexandria. Definitely go check it out.

This [thesis](https://www.duo.uio.no/bitstream/handle/10852/53769/1/master.pdf) covers some of the topics well, especially Lazy SMP.

---

### Libraries 

Chess package for Go: https://github.com/Sidhant-Roymoulik/chess
