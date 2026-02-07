package ebitenGUI

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	//Доска
	Board *ebiten.Image
	//Выбор клетки
	CellExtractors [3]*ebiten.Image
	//Фигуры
	PiecesImage [12]*ebiten.Image
}

func Assets() *Game {
	g := &Game{}

	var dirs [12]string = [12]string{
		"WhiteKing", "BlackKing",
		"WhiteQueen", "BlackQueen",
		"WhiteBishop", "BlackBishop",
		"WhiteKnight", "BlackKnight",
		"WhiteRook", "BlackRook",
		"WhitePawn", "BlackPawn",
	}

	//Импорт доски
	board, _, err := ebitenutil.NewImageFromFile("Assets/Board.png")
	if err != nil {
		log.Fatal(err)
	}
	g.Board = board

	clickImg, _, err := ebitenutil.NewImageFromFile("Assets/Click.png")
	if err != nil {
		log.Fatal(err)
	}
	g.CellExtractors[0] = clickImg

	selectImg, _, err := ebitenutil.NewImageFromFile("Assets/Select.png")
	if err != nil {
		log.Fatal(err)
	}
	g.CellExtractors[1] = selectImg

	eatImg, _, err := ebitenutil.NewImageFromFile("Assets/Eat.png")
	if err != nil {
		log.Fatal(err)
	}
	g.CellExtractors[2] = eatImg

	//Импорт всех фигур
	for i := 0; i < 12; i++ {
		piece, _, err := ebitenutil.NewImageFromFile("Assets/" + dirs[i] + ".png")
		if err != nil {
			log.Fatal(err)
		}
		g.PiecesImage[i] = piece
	}

	return g
}

// Логическая карта шахматов, значения XY это верхняя-левая точка каждой клетки
var mapX [8]float64 = [8]float64{6, 138, 270, 402, 534, 666, 798, 930}
var mapY [8]float64 = [8]float64{6, 138, 270, 402, 534, 666, 798, 930}

func Click() (int, int, bool) {
	cursorX, cursorY := ebiten.CursorPosition()
	var clickX int
	var clickY int
	for i := 0; i < 8; i++ {
		if cursorX >= int(mapX[i]-6) && cursorX <= int(mapX[i])+128 { // 6 это половина толшины стенок между клетками, 122 это размер клетки, 128 это вся клетка + пол стены
			clickX = i
			for i := 0; i < 8; i++ {
				if cursorY >= int(mapY[i]-6) && cursorY <= int(mapY[i])+128 {
					clickY = i
				}
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return clickX, clickY, true
	}
	return clickX, clickY, false
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//Доска
	bOard := &ebiten.DrawImageOptions{}
	bOard.GeoM.Translate(0, 0)
	screen.DrawImage(g.Board, bOard)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1056, 1056 // логический размер окна
}

func Main() {

	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("ШАХМАТЫ")

	if err := ebiten.RunGame(Assets()); err != nil {
		log.Fatal(err)
	}

}
