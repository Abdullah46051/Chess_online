package chessEngine

import (
	"encoding/json"
	"fmt"
	"os"
)

type Piece struct {
	X, Y int
}

type PieceType struct {
	Type  string
	Color bool
	Piece []Piece
}

type Board struct {
	Pieces []PieceType
}

type Moves struct {
	movesX, movesY   []int
	eatX, eatY       []int
	attackX, attackY []int
}

var history []Board
var board Board

func SearchPiece(X int, Y int) (int, int) {
	for i := 0; i <= 11; i++ {
		for j := 0; j < len(board.Pieces[i].Piece); j++ {
			x := board.Pieces[i].Piece[j].X
			y := board.Pieces[i].Piece[j].Y
			if x == X && y == Y {
				return i, j
			}
		}
	}
	return -1, -1
}

func SearchMoves(piece int, num int) Moves {
	Type := board.Pieces[piece].Type
	color := board.Pieces[piece].Color
	x := board.Pieces[piece].Piece[num].X
	y := board.Pieces[piece].Piece[num].Y
	var movesX []int
	var movesY []int
	var eatX []int
	var eatY []int
	var attackX []int
	var attackY []int
	var result Moves

	if Type == "King" {
		kingMovesX := [8]int{+1, -1, 0, 0, +1, +1, -1, -1}
		kingMovesY := [8]int{0, 0, +1, -1, +1, -1, -1, +1}
		for i := 0; i <= 7; i++ {
			movesX = append(movesX, x+kingMovesX[i])
			movesY = append(movesY, y+kingMovesY[i])
			eatX = append(eatX, -1)
			eatY = append(eatY, -1)
		}
	}

	if Type == "Knight" {
		knightMovesX := [8]int{-2, +2, +2, -2, -1, +1, +1, -1}
		knightMovesY := [8]int{-1, +1, -1, +1, -2, +2, -2, +2}
		for i := 0; i <= 7; i++ {
			movesX = append(movesX, x+knightMovesX[i])
			movesY = append(movesY, y+knightMovesY[i])
			eatX = append(eatX, -1)
			eatY = append(eatY, -1)
		}
	}

	if Type == "Pawn" {
		var v, cV, enV, enP int
		if color {
			enV = 3
			enP = 11
			v = -1
			cV = 6
		} else {
			enV = 5
			enP = 10
			v = +1
			cV = 1
		}

		//Обычный ход
		if ok, _ := SearchPiece(x, y+v); ok == -1 {
			movesX = append(movesX, x)
			movesY = append(movesY, y+v)
			eatX = append(eatX, -1)
			eatY = append(eatY, -1)
		}

		//Длинный ход
		if ok, _ := SearchPiece(x, y+(v*2)); y == cV && ok == -1 {
			if ok, _ := SearchPiece(x, y+v); ok == -1 {
				movesX = append(movesX, x)
				movesY = append(movesY, y+(v*2))
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
		}

		//Взятие
		if ok, _ := SearchPiece(x-1, y+v); ok != -1 && board.Pieces[ok].Color != color {
			movesX = append(movesX, x-1)
			movesY = append(movesY, y+v)
			eatX = append(eatX, x-1)
			eatY = append(eatY, y+v)
		}
		if ok, _ := SearchPiece(x+1, y+v); ok != -1 && board.Pieces[ok].Color != color {
			movesX = append(movesX, x+1)
			movesY = append(movesY, y+v)
			eatX = append(eatX, x+1)
			eatY = append(eatY, y+v)
		}

		//Бой
		if ok, _ := SearchPiece(x-1, y+v); (ok != -1 && board.Pieces[ok].Color != color) || ok == -1 {
			attackX = append(attackX, x-1)
			attackY = append(attackY, y+v)
		}
		if ok, _ := SearchPiece(x+1, y+v); (ok != -1 && board.Pieces[ok].Color != color) || ok == -1 {
			attackX = append(attackX, x+1)
			attackY = append(attackY, y+v)
		}

		//Взятие на проходе
		if y == enV {
			t, p := SearchPiece(x-1, y)
			if t == enP {
				xB := history[len(history)-1].Pieces[t].Piece[p].X
				yB := history[len(history)-1].Pieces[t].Piece[p].Y
				if (xB == x-1) && (yB == y-2) {
					movesX = append(movesX, x-1)
					movesY = append(movesY, y+v)
					eatX = append(eatX, x-1)
					eatY = append(eatY, y)
				}
			}

			t, p = SearchPiece(x+1, y)
			if t == enP {
				xB := history[len(history)-1].Pieces[t].Piece[p].X
				yB := history[len(history)-1].Pieces[t].Piece[p].Y
				if (xB == x+1) && (yB == y-2) {
					movesX = append(movesX, x+1)
					movesY = append(movesY, y+v)
					eatX = append(eatX, x+1)
					eatY = append(eatY, y)
				}
			}
		}
	}

	if Type == "Queen" || Type == "Bishop" {
		m1, m2, m3, m4 := true, true, true, true
		for i := 1; i <= 8; i++ {
			if m1 {
				if ok, _ := SearchPiece(x+i, y+i); ok != -1 {
					m1 = false
				}
				movesX = append(movesX, x+i)
				movesY = append(movesY, y+i)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
			if m2 {
				if ok, _ := SearchPiece(x-i, y-i); ok != -1 {
					m2 = false
				}
				movesX = append(movesX, x-i)
				movesY = append(movesY, y-i)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
			if m3 {
				if ok, _ := SearchPiece(x-i, y+i); ok != -1 {
					m3 = false
				}
				movesX = append(movesX, x-i)
				movesY = append(movesY, y+i)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
			if m4 {
				if ok, _ := SearchPiece(x+i, y-i); ok != -1 {
					m4 = false
				}
				movesX = append(movesX, x+i)
				movesY = append(movesY, y-i)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
		}
	}

	if Type == "Queen" || Type == "Rook" {
		m1, m2, m3, m4 := true, true, true, true
		for i := 1; i <= 8; i++ {
			if m1 {
				if ok, _ := SearchPiece(x+i, y); ok != -1 {
					m1 = false
				}
				movesX = append(movesX, x+i)
				movesY = append(movesY, y)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
			if m2 {
				if ok, _ := SearchPiece(x-i, y); ok != -1 {
					m2 = false
				}
				movesX = append(movesX, x-i)
				movesY = append(movesY, y)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
			if m3 {
				if ok, _ := SearchPiece(x, y+i); ok != -1 {
					m3 = false
				}
				movesX = append(movesX, x)
				movesY = append(movesY, y+i)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
			if m4 {
				if ok, _ := SearchPiece(x, y-i); ok != -1 {
					m4 = false
				}
				movesX = append(movesX, x)
				movesY = append(movesY, y-i)
				eatX = append(eatX, -1)
				eatY = append(eatY, -1)
			}
		}
	}

	for i := len(movesX) - 1; i >= 0; i-- {
		if (movesX[i] < 0 || movesX[i] > 7) || (movesY[i] < 0 || movesY[i] > 7) { //Проверка хода на пределы координат
			movesX = append(movesX[:i], movesX[i+1:]...)
			movesY = append(movesY[:i], movesY[i+1:]...)

			eatX = append(eatX[:i], eatX[i+1:]...)
			eatY = append(eatY[:i], eatY[i+1:]...)
			continue
		}

		ok, _ := SearchPiece(movesX[i], movesY[i])
		if ok != -1 {
			if board.Pieces[ok].Color == color {
				movesX = append(movesX[:i], movesX[i+1:]...)
				movesY = append(movesY[:i], movesY[i+1:]...)

				eatX = append(eatX[:i], eatX[i+1:]...)
				eatY = append(eatY[:i], eatY[i+1:]...)
				continue
			}
		}

		if ok != -1 {
			if (board.Pieces[ok].Color != color) && Type != "Pawn" {
				eatX[i] = movesX[i]
				eatY[i] = movesY[i]
			}
		}

		if ok != -1 {
			if ((board.Pieces[ok].Color != color) || ok == -1) && Type != "Pawn" {
				attackX = append(attackX, movesX[i])
				attackY = append(attackY, movesY[i])
			}
		} else {
			attackX = append(attackX, movesX[i])
			attackY = append(attackY, movesY[i])
		}
	}

	result = Moves{
		movesX: movesX, movesY: movesY,
		eatX: eatX, eatY: eatY,
		attackX: attackX, attackY: attackY,
	}

	return result
}

func MoveTo(X int, Y int, Xto int, Yto int) bool {
	tI, pI := SearchPiece(X, Y)
	if (tI != -1) && (pI != -1) {
		moves := SearchMoves(tI, pI)
		history = append(history, board)

		for i := 0; i <= len(moves.movesX)-1; i++ {
			if Xto == moves.movesX[i] && Yto == moves.movesY[i] {
				board.Pieces[tI].Piece[pI].X = Xto
				board.Pieces[tI].Piece[pI].Y = Yto

				if (moves.eatX[i] != -1) && (moves.eatY[i] != -1) {
					eatTi, eatPi := SearchPiece(moves.eatX[i], moves.eatY[i])
					board.Pieces[eatTi].Piece =
						append(
							board.Pieces[eatTi].Piece[:eatPi],
							board.Pieces[eatTi].Piece[eatPi+1:]...,
						)
				}
				return true
			}
		}
	}

	return false
}

func Main() {
	data, _ := os.ReadFile("chessEngine/Board.json")
	json.Unmarshal(data, &board)
	g := SearchMoves(6, 0)
	fmt.Println("dd", g.movesX, g.movesY)
	MoveTo(1, 7, 0, 5)
	g = SearchMoves(6, 0)
	fmt.Println("ss", g.movesX, g.movesY)
	fmt.Println(board.Pieces[6].Piece[0])
}
