package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var cellReg *regexp.Regexp

func init() {
	cellReg = regexp.MustCompile(`(?P<piece>\d+)\((?P<orientation>\d+)\)`)
}

var Pieces map[int]*Piece

type Piece struct {
	Up, Down, Left, Right int
	Orientation           int //0, 1, 2, 3
	Number                int
}

func (this *Piece) SetOrientation(o int) {
	for o != this.Orientation {
		oldLeft := this.Left
		this.Left = this.Down
		this.Down = this.Right
		this.Right = this.Up
		this.Up = oldLeft

		this.Orientation++
		if this.Orientation > 3 {
			this.Orientation = 0
		}
	}
}

func loadPieces() {
	raw, err := ioutil.ReadFile("assets/pieces.txt")
	if err != nil {
		panic(err)
	}
	//***
	Pieces = make(map[int]*Piece)
	//***
	txt := string(raw)
	lines := strings.Split(strings.Replace(txt, "\r\n", "\n", -1), "\n")
	//***
	for idx, line := range lines {
		values := strings.Split(line, " ")
		if len(values) != 4 {
			return
		}
		var add Piece
		add.Up, _ = strconv.Atoi(values[0])
		add.Down, _ = strconv.Atoi(values[1])
		add.Left, _ = strconv.Atoi(values[2])
		add.Right, _ = strconv.Atoi(values[3])
		add.Number = idx
		Pieces[idx] = &add
	}
}

type Board struct {
	Content [16][16]*Piece
}

func (this *Board) CountPoints() (res int) {
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			if this.Content[y][x] != nil {
				res++
			}
		}
	}
	return
}

func (this *Board) IsCorrect() bool {

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {

			p := this.Content[y][x]
			if p == nil {
				continue
			}

			// Check value
			if p.Number < 0 || p.Number > 255 {
				fmt.Println("Invalid value in", y, x, "", p.Number)
				return false
			}

			// Check compatibility
			if !this.checkPiece(x, y, p) {
				return false
			}
		}
	}
	return true
}

func (this *Board) checkPiece(x, y int, p *Piece) bool {
	// Up
	if y == 0 {
		if p.Up != 0 {
			fmt.Println("Invalid value in", y, x, " UP", p.Up)
			return false
		}
	} else {
		other := this.Content[y-1][x]
		if other != nil && other.Down != p.Up {
			fmt.Println("Invalid in", y, x, " : value UP is", p.Up, "should be", other.Down)
			return false
		}
	}
	// Down
	if y == 15 {
		if p.Down != 0 {
			fmt.Println("Invalid value in", y, x, " DOWN", p.Down)
			return false
		}
	} else {
		other := this.Content[y+1][x]
		if other != nil && other.Up != p.Down {
			fmt.Println("Invalid in", y, x, " : value DOWN is", p.Down, "should be", other.Up)
			return false
		}
	}
	// Left
	if x == 0 {
		if p.Left != 0 {
			fmt.Println("Invalid value in", y, x, " LEFT", p.Left)
			return false
		}
	} else {
		other := this.Content[y][x-1]
		if other != nil && other.Right != p.Left {
			fmt.Println("Invalid in", y, x, " : value LEFT is", p.Left, "should be", other.Right)
			return false
		}
	}
	// Right
	if x == 15 {
		if p.Right != 0 {
			fmt.Println("Invalid value in", y, x, " RIGHT", p.Right)
			return false
		}
	} else {
		other := this.Content[y][x+1]
		if other != nil && other.Left != p.Right {
			fmt.Println("Invalid in", y, x, " : value RIGHT is", p.Right, "should be", other.Left)
			return false
		}
	}
	return true
}

func (this *Board) Print() {
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			p := this.Content[y][x]
			if p == nil {
				fmt.Print("X")
			} else {
				fmt.Print(p.Number, "(", p.Orientation, ")")
			}
			if x != 15 {
				fmt.Print("-")
			}
		}
		fmt.Println("")
	}
}

func loadBoard(file string) (res Board) {

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	data := string(raw)
	data = strings.Replace(data, "\r\n", "\n", -1)
	lines := strings.Split(data, "\n")
	for y, line := range lines {
		cells := strings.Split(line, "-")
		for x, cell := range cells {
			if cell == "X" {
				continue
			}
			r := cellReg.FindStringSubmatch(cell)
			captures := make(map[string]string)
			names := cellReg.SubexpNames()
			for i, name := range names {
				v := r[i]
				captures[name] = v
			}
			res.Content[y][x] = Pieces[toInt(captures["piece"])]
			res.Content[y][x].SetOrientation(toInt(captures["orientation"]))
		}
	}

	return
}

func toInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func main() {
	fmt.Println("Start")
	loadPieces()
	fmt.Println(len(Pieces), "pieces loaded")
	board := loadBoard(os.Args[1])
	board.Print()
	if !board.IsCorrect() {
		fmt.Println("Exit : Incorrect board.")
		return
	}
	fmt.Println("You got", board.CountPoints(), "pieces correctly placed.")
}
