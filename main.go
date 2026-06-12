package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const paddleSymbol = 0x2588
const ballSymbol = 0x25CF

const paddleHeight = 4
const initialVelocityRow = 1
const initialVelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var isGamePaused bool
var debugLog string

var gameObjects []*GameObject

func main() {

	initScreen()
	InitGameState()

	inputChan := InitUserInput()

	for !isGameOver() {
		HandleUserInput(ReadInput(inputChan))
		updateState()
		drawState()

		time.Sleep(100 * time.Millisecond)
	}

	screenWidth, screenHeight := screen.Size()
	winner := getWinner()
	PrintStringCentered(screenHeight/2-1, screenWidth/2, "Game Over...")
	PrintStringCentered(screenHeight/2, screenWidth/2, fmt.Sprintf("%s wins...", winner))
	screen.Show()

	time.Sleep(3 * time.Second)
	screen.Fini()
}

func updateState() {
	if isGamePaused {
		return
	}

	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if collidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if collidesWithPaddle(ball, player1Paddle) || collidesWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}
}

func drawState() {
	if isGamePaused {
		return
	}

	screen.Clear()
	screen.SetTitle("Ping Pong")

	PrintString(0, 0, debugLog)

	for _, obj := range gameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

func collidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func collidesWithPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesOnColumn bool
	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}

	return collidesOnColumn &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
}

func isGameOver() bool {
	return getWinner() != ""
}

func getWinner() string {
	screenWidth, _ := screen.Size()

	if ball.col < 0 {
		return "Player 1"
	} else if ball.col >= screenWidth {
		return "Player 2"
	} else {
		return ""
	}
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func HandleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < screenHeight {
		player1Paddle.row++
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < screenHeight {
		player2Paddle.row++
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	}
}

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - paddleHeight/2

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, height: paddleHeight,
		velRow: 0, velCol: 0,
		symbol: paddleSymbol,
	}

	player2Paddle = &GameObject{
		row: paddleStart, col: width - 1, width: 1, height: paddleHeight,
		velRow: 0, velCol: 0,
		symbol: paddleSymbol,
	}

	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1,
		velRow: initialVelocityRow, velCol: initialVelocityCol,
		symbol: ballSymbol,
	}

	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}

}

func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}

func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}
