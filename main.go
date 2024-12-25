package main

import (
	"errors"
	"fmt"
	"hash/maphash"
	"math/rand"
	"os"
	"time"

	//"github.com/pkg/profile"
)

type Point int

const (
	Pin Point = iota
	Empty
	Invalid
)

type Board [7][7]Point

func printBoard(board *Board) {
	for _, a := range board {
		for _, b := range a {
			switch b {
			case Pin:
				fmt.Print("*")
			case Empty:
				fmt.Print(" ")
			case Invalid:
				fmt.Print("I")
			default:
				panic("fail")
			}
		}

		fmt.Println()
	}
	fmt.Println()
}

func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(int64(Rand64())))
}

func Rand64() uint64 {
	return new(maphash.Hash).Sum64()
}

// Returns the x and y position of a random pin on the board
func randomPin(board *Board) (int, int) {
	for {
		x := rand.Intn(7)
		y := rand.Intn(7)

		//x := r.Intn(7)
		//y := r.Intn(7)
		if board[x][y] == Pin {
			return x, y
		}
	}
}

// X axis is 0, Y axis is 1
func randomEmpty(board *Board, x, y, axis int) (int, int, error) {
	//for j := 0; j < 25; j++ {
	for i := 0; i < 7; i++ {
		//i := rand.Intn(7)

		theX := x
		theY := y

		if axis == 0 { // X-axis
			theX = i
		} else if axis == 1 { // Y-axis
			theY = i
		} else {
			panic("failure")
		}

		point := board[theX][theY]
		if point == Empty {
			return theX, theY, nil
		}
	}

	return 0, 0, errors.New("fail!!")
}

func isValidMove(board *Board, fromX, fromY, toX, toY, axis int) bool {
	if fromX < 0 || fromX >= 7 {
		return false
	}
	if toX < 0 || toX >= 7 {
		return false
	}
	if fromY < 0 || fromY >= 7 {
		return false
	}
	if toY < 0 || toY >= 7 {
		return false
	}

	from := board[fromX][fromY]
	to := board[toX][toY]

	if from != Pin {
		panic("fail 1")
	}

	if to != Empty {
		return false
		//panic("fail 2")
	}

	abs := func(n int) int {
		if n < 0 {
			return n * -1
		}
		return n
	}

	xDiff := abs(fromX - toX)
	yDiff := abs(fromY - toY)

	// Jump length guarding
	if !(xDiff == 0 || xDiff == 2) || !(yDiff == 0 || yDiff == 2) {
		return false
	}

	// Jump axis guarding
	if xDiff == 2 && yDiff != 0 {
		return false
	}

	if yDiff == 2 && xDiff != 0 {
		return false
	}

	midX := fromX
	midY := fromY
	if axis == 0 {
		if toX > fromX {
			midX++
		} else {
			midX--
		}
	} else {
		if toY > fromY {
			midY++
		} else {
			midY--
		}
	}

	mid := board[midX][midY]
	if mid != Pin {
		return false
	}

	return true
}

func makeMove(board *Board, fromX, fromY, toX, toY, axis int) {
	midX := fromX
	midY := fromY
	if axis == 0 {
		if toX > fromX {
			midX++
		} else {
			midX--
		}
	} else {
		if toY > fromY {
			midY++
		} else {
			midY--
		}
	}

	board[fromX][fromY] = Empty
	board[midX][midY] = Empty
	board[toX][toY] = Pin
}

func tryMove(board *Board, moves *[]Move) bool {
	for i := 0; i < 90; i++ {
		x, y := randomPin(board)
		axis := rand.Int() & 1
		wx := x
		wy := y
		if axis == 0 {
			if rand.Int() & 1 == 0 {
				wx += 2
			} else {
				wx -= 2
			}
			//wx += rand.Intn(5) - 2
		} else {
			if rand.Int() & 1 == 0 {
				wy += 2
			} else {
				wy -= 2
			}
			//wy += rand.Intn(5) - 2
		}
		/*wx, wy, err := randomEmpty(board, x, y, axis)
		if err != nil {
			continue
		}*/

		if isValidMove(board, x, y, wx, wy, axis) {
			makeMove(board, x, y, wx, wy, axis)
			*moves = append(*moves, Move{
				fromX: x,
				fromY: y,
				toX: wx,
				toY: wy,
			})
			return true
		}
	}

	return false
}

func isSolved(board *Board) bool {
	count := 0

	for _, a := range board {
		for _, b := range a {
			if b == Pin {
				count++
			}

			if count > 1 {
				return false
			}
		}
	}

	return true
}

func countPins(board *Board) int {
	count := 0

	for _, a := range board {
		for _, b := range a {
			if b == Pin {
				count++
			}
		}
	}

	return count
}

type Move struct {
	fromX int
	fromY int
	toX int
	toY int
}

func printMoves(moves []Move) {
	for _, move := range moves {
		msg := ""
		if move.toX > move.fromX {
			msg = "to the right"
		} else if move.toX < move.fromX {
			msg = "to the left"
		} else if move.toY > move.fromY {
			msg = "down"
		} else if move.toY < move.fromY {
			msg = "up"
		}
		fmt.Println("Move", move.fromX+1, move.fromY+1, msg)
	}
}

func trySolve() (Board, error) {
	var board Board
	board[3][3] = Empty
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			board[x * 5][y * 5] = Invalid
			board[x * 5][y * 5 + 1] = Invalid
			board[x * 5 + 1][y * 5] = Invalid
			board[x * 5 + 1][y * 5 + 1] = Invalid
		}
	}

	moves := make([]Move, 0)

	for {
		if !tryMove(&board, &moves) {
			nPins := countPins(&board)
			if nPins < 3 {
				fmt.Println("game ended at", countPins(&board), "pins")
				printBoard(&board)
				if nPins < 2 {
					printMoves(moves)
					return board, nil
				}
			}
			return Board{}, errors.New("nothing")
		}

		//printBoard(&board)

		/*if isSolved(&board) {
			return board, nil
		}*/
	}
}

var r *rand.Rand

func main() {
	//defer profile.Start().Stop()

	r = NewRand()

	//rand.Seed(time.Now().UnixNano())
	r.Seed(time.Now().UnixNano())

	for {
		b, err := trySolve()
		if err != nil {
			continue
		}

		printBoard(&b)
		os.Exit(0)
	}
}
