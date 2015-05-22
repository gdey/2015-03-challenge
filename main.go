package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/jroimartin/gocui"
)

const (
	VERTICAL   = true
	HORIZONTAL = false
)

var (
	grid        = make([][]GridSquare, 16)
	shipTypes   = []string{"aircraft_carrier", "battleship", "submarine", "destroyer", "cruiser", "patrol_boat"}
	shipLengths = map[string]int{
		"aircraft_carrier": 5,
		"battleship":       4,
		"submarine":        3,
		"destroyer":        3,
		"cruiser":          3,
		"patrol_boat":      2,
	}
	shipPoints = map[string]int{
		"aircraft_carrier": 20,
		"battleship":       12,
		"submarine":        6,
		"destroyer":        6,
		"cruiser":          6,
		"patrol_boat":      2,
	}
	boatHits map[string]int

	playerScore int
)

var inputText string

func init() {
	for i, _ := range grid {
		grid[i] = make([]GridSquare, 16)
	}
	boatHits = make(map[string]int)
}

type GridSquare struct {
	HasShip  bool
	BeenHit  bool
	ShipType string
}

type Coordinate struct {
	X int
	Y int
}

func placeShips() {
	for _, ship := range shipTypes {
		placeBoat(ship)
	}
}

func isSunk(boat []GridSquare) bool {
	for _, square := range boat {
		if !square.BeenHit {
			return false
		}
	}
	return true
}

func placeBoat(boatType string) {
	orientation := randOrientation()
	if orientation == VERTICAL {
		start := Coordinate{
			X: rand.Intn(16),
			Y: rand.Intn(16 - shipLengths[boatType]),
		}
		squares := grid[start.X][start.Y : start.Y+shipLengths[boatType]]
		for _, square := range squares {
			if square.HasShip {
				placeBoat(boatType)
				return
			}
		}

		for i, _ := range squares {
			squares[i].HasShip = true
			squares[i].ShipType = boatType
		}
		return
	} else {
		start := Coordinate{
			X: rand.Intn(16 - shipLengths[boatType]),
			Y: rand.Intn(16),
		}
		rows := grid[start.X : start.X+shipLengths[boatType]]
		for _, row := range rows {
			if row[start.Y].HasShip {
				placeBoat(boatType)
				return
			}
		}

		for i, _ := range rows {
			rows[i][start.Y].HasShip = true
			rows[i][start.Y].ShipType = boatType
		}
		return
	}
}

func randOrientation() bool {
	randomInt := rand.Intn(2)
	if randomInt == 1 {
		return VERTICAL
	} else {
		return HORIZONTAL
	}
}

func shoot(x int, y int) error {
	if x > 15 || y > 15 {
		return errors.New(fmt.Sprintf("You cannot shoot at location (%d, %d): it is off the board\n", x, y))
	}
	square := grid[x][y]
	if square.BeenHit {
		return errors.New(fmt.Sprintf("You have already shot at square (%d, %d)\n", x, y))
	}
	if square.HasShip {
		fmt.Printf("You hit a ship at location (%d, %d)!\n", x, y)
		boatHits[square.ShipType] += 1
		if boatHits[square.ShipType] == shipLengths[square.ShipType] {
			fmt.Printf("You've sunk my %s!\n", square.ShipType)
			playerScore += shipPoints[square.ShipType]
		}
	} else {
		fmt.Printf("You missed at (%d, %d)\n", x, y)
		playerScore -= 1
	}
	grid[x][y].BeenHit = true

	return nil
}

func getCoordinate() (x int, y int) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Input a coordinate:")
	fmt.Printf("\tx: ")
	xStr, err := reader.ReadString([]byte("\n")[0])
	if err != nil {
		fmt.Printf("Error reading from stdin: %s", err.Error())
		return getCoordinate()
	}
	if len(xStr) > 1 {
		xStr = xStr[:len(xStr)-1]
	}
	fmt.Printf("\ty: ")
	yStr, err := reader.ReadString([]byte("\n")[0])
	if err != nil {
		fmt.Printf("Error reading from stdin: %s", err.Error())
		return getCoordinate()
	}
	if len(yStr) > 1 {
		yStr = yStr[:len(yStr)-1]
	}
	x, err = strconv.Atoi(xStr)
	if err != nil {
		fmt.Printf("Error converting `%s` to int: %s\n", xStr, err.Error())
		return getCoordinate()
	}
	y, err = strconv.Atoi(yStr)
	if err != nil {
		fmt.Printf("Error converting `%s` to int: %s\n", xStr, err.Error())
		return getCoordinate()
	}
	return x, y
}

func takeTurn() {
	coordinates := make([]Coordinate, 5)
	for i := 0; i < 5; i++ {
		x, y := getCoordinate()
		coordinates[i] = Coordinate{
			X: x,
			Y: y,
		}
	}

	var fails int
	for _, point := range coordinates {
		err := shoot(point.X, point.Y)
		if err != nil {
			fmt.Printf("Error shooting at (%d, %d): %s", point.X, point.Y, err.Error())
			fails += 1
		}
		if haveYouWon() {
			win()
			return
		}
	}

	for fails > 0 {
		coordinates := make([]Coordinate, fails)
		for i := 0; i < fails; i++ {
			x, y := getCoordinate()
			coordinates[i] = Coordinate{
				X: x,
				Y: y,
			}
		}

		fails = 0
		for _, point := range coordinates {
			err := shoot(point.X, point.Y)
			if err != nil {
				fmt.Printf("Error shooting at (%d, %d): %s\n", point.X, point.Y, err.Error())
				fails += 1
			}
		}
	}
}

func win() {
	fmt.Printf("You won!\n")
	fmt.Printf("Your final score was: %d\n", playerScore)
	return
}

func lose() {
	fmt.Printf("Unfortunately, you lost.\n")
	fmt.Printf("Your final score was: %d\n", playerScore)
}

func haveYouWon() bool {
	for _, boat := range shipTypes {
		if boatHits[boat] < shipLengths[boat] {
			return false
		}
	}
	return true
}

func playGame(g *gocui.Gui) {
	/*
		iv, err := g.View("input")
		if err != nil {
			log.Panicln("input", err)
		}
	*/
	ov, err := g.View("output")
	if err != nil {
		if err != gocui.ErrorUnkView {
			return
		}
		log.Panicln("output", err)
	}
	for i := 0; i < 6; i++ {
		/*
			takeTurn()
			if haveYouWon() {
				win()
				return
			}
		*/
		//		fmt.Fprintf(iv, "After turn %d, your score is %d\n", i+1, playerScore)
		ov.Clear()
		fmt.Fprintln(ov, "Hello World")
		//printBoardAndCheat(ov)
		g.Flush()
	}
	/*
		if haveYouWon() {
			win()
			return
		} else {
			lose()
			return
		}
	*/
}

/*
func printBoard() {
	for i, row := range grid {
		for j, cell := range row {
			var c string
			if cell.BeenHit {
				c = "X"
			} else {
				c = fmt.Sprintf("%d,%d", i, j)
			}
			fmt.Printf("%v\t", c)
		}
		fmt.Println()
	}
}
*/
func printBoardAndCheat(outputView *gocui.View) {
	const ship rune = 's'
	const blank rune = 'b'
	const hitShip rune = 'h'
	for i, row := range grid {
		for j, cell := range row {
			var c rune
			switch {
			case cell.BeenHit:
				c = hitShip
			case cell.HasShip:
				c = ship
			default:
				//TODO: Do this right
				c = blank
			}
			outputView.SetCursor(i, j)
			outputView.EditWrite(c)
		}
	}
}

func main() {
	/*
		placeShips()
		for i, row := range grid {
			for j, cell := range row {
				var c string
				if cell.ShipType != "" {
					c = fmt.Sprintf("X")
				} else {
					c = fmt.Sprintf("%d,%d", i, j)
				}
				fmt.Printf("%v ", c)
			}
			fmt.Printf("\n")
		}
	*/
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetLayout(layout)
	inputText = "Enter the x coord: "

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SetCurrentView("input")

	if err := g.MainLoop(); err != nil && err != gocui.Quit {
		log.Panicln(err)
	}

}
func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
		return err
	}

	for _, a := range []rune("1234567890,") {
		fmt.Println(a)
		aa := a
		if err := g.SetKeybinding("", a, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			iv, _ := g.View("input")
			iv.SetCursor(len(inputText), 0)
			inputText += string(a)
			iv.EditWrite(aa)
			ov, _ := g.View("output")
			ov.EditWrite('F')

			return nil
		}); err != nil {
			log.Panicln(a, err)
			return err
		}

	}

	/*
		if err := g.SetKeybinding("stdin", gocui.KeyArrowUp, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				scrollView(v, -1)
				return nil
			}); err != nil {
			return err
		}
			if err := g.SetKeybinding("stdin", gocui.KeyArrowDown, gocui.ModNone,
				func(g *gocui.Gui, v *gocui.View) error {
					scrollView(v, 1)
					return nil
				}); err != nil {
				return err
			}
	*/
	return nil
}
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("output", 0, 0, maxX-1, maxY-20); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		printBoardAndCheat(v)
	}
	if v, err := g.SetView("input", 0, maxY-20, maxX-1, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, inputText)
	}
	g.ShowCursor = true

	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}
