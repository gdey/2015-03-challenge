//go:generate stringer -type=shipId
package main

import (
	"fmt"
	"math/rand"
)

type shipId int

const (
	Ocean shipId = iota
	AircraftCarrier
	Battleship
	Submarine
	Destroyer
	Cruiser
	PatrolBoat
)

func (s shipId) Size() int {
	switch s {
	case AircraftCarrier:
		return 5
	case Battleship:
		return 4
	case Submarine:
		return 3
	case Destroyer:
		return 3
	case Cruiser:
		return 3
	case PatrolBoat:
		return 2
	default:
		return 1
	}
}

func (s shipId) Score() int {
	switch s {
	case AircraftCarrier:
		return 20
	case Battleship:
		return 12
	case Submarine:
		return 6
	case Destroyer:
		return 6
	case Cruiser:
		return 6
	case PatrolBoat:
		return 2
	default:
		return 0
	}
}

type Ship struct {
	X, Y  int
	Type  shipId
	Id    int
	Horiz bool
	init  bool
}

func (s *Ship) String() string {
	var horzstring = "Vertical"
	if s.Horiz {
		horzstring = "Horizontal"
	}
	return fmt.Sprintf("Ship %s (%d -- %d), x %d,  y %d, %s", s.Type, s.Type.Size(), s.Id, s.X, s.Y, horzstring)
}

func RandShip() *Ship {
	s := Ship{}
	s.Type = shipId(rand.Intn(5) + 1)
	s.X = rand.Intn(16)
	s.Y = rand.Intn(16)
	s.Horiz = rand.Intn(2) == 1
	s.init = true
	return &s
}

func (s *Ship) MoveIntoBoard(width, height int) {
	if s.Horiz {
		if s.X > (width - s.Type.Size()) {
			// nudge
			s.X = width - s.Type.Size()
			return
		}
		return
	}
	if s.Y > (height - s.Type.Size()) {
		// nudge
		s.Y = height - s.Type.Size()
		return
	}
	return
}

func (s *Ship) RenderToOcean(o *[][]int) {
	oo := *o
	if s.Horiz {
		for i := 0; i < s.Type.Size(); i++ {
			oo[s.X+i][s.Y] = s.Id
		}
		return

	}
	for i := 0; i < s.Type.Size(); i++ {
		oo[s.X][s.Y+i] = s.Id
	}
	return

}

var Flotilla []*Ship

func main() {

	Flotilla = make([]*Ship, 5)
	rand.Seed(122365478)
	oceanMap := make([][]int, 16)
	for i := 0; i < 16; i++ {
		oceanMap[i] = make([]int, 16)
	}

	for i, s := range Flotilla {
		s = RandShip()
		s.Id = i + 1
		s.MoveIntoBoard(16, 16)
		s.RenderToOcean(&oceanMap)
		fmt.Println(s)
	}
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			if oceanMap[x][y] == 0 {
				fmt.Print("  ")
			} else {
				fmt.Print(oceanMap[x][y], " ")
			}
		}
		fmt.Println()
	}

}
