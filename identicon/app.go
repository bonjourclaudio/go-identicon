package identicon

import (
	"crypto/md5"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
)

type Identicon struct {
	Title string
	Hash [16]byte
	Color [3]byte
	Grid []byte
	GridPoints []GridPoint
	PixelMap []DrawingPoint
}

type GridPoint struct {
	Value byte
	Index int
}

type Point struct {
	X int
	Y int
}

type DrawingPoint struct {
	TopLetft Point
	BottomRight Point
}

type Apply func(Identicon) Identicon

func Pipe(identicon Identicon, funcs ...Apply) Identicon {
	for _, applyer := range funcs {
		identicon = applyer(identicon)
	}
	return identicon
}

func HashInput(input []byte) Identicon {
	checkSum := md5.Sum(input)

	return Identicon{
		Title: string(input),
		Hash:  checkSum,
	}
}

func PickColor(identicon Identicon) Identicon {

	rgb := [3]byte{}

	copy(rgb[:], identicon.Hash[:3])

	identicon.Color = rgb

	return identicon
}

func BuildGrid(identicon Identicon) Identicon {

	grid := []byte{}

	for i := 0; i < len(identicon.Hash) && i+3 <= len(identicon.Hash) -1; i += 3 {

		chunk := make([]byte, 5)

		copy(chunk, identicon.Hash[i:i+3])
		chunk[3] = chunk[1]
		chunk[4] = chunk[0]
		grid = append(grid, chunk ...)

	}
	identicon.Grid = grid
	return identicon
}

func FilterOddSquares(identicon Identicon) Identicon {

	grid := []GridPoint{}

	for i, code := range identicon.Grid {

		if code%2 == 0 {

			point := GridPoint{
				Value: code,
				Index: i,
			}

			grid = append(grid, point)

		}

	}
	identicon.GridPoints = grid
	return identicon
}

func BuildPixelMap(identicon Identicon) Identicon {
	drawingPoints := []DrawingPoint{}

	pixelFunc := func(p GridPoint) DrawingPoint {
		horizontal := (p.Index % 5) * 50
		vertical := (p.Index / 5) * 50
		topLeft := Point{horizontal, vertical}
		bottomRight := Point{horizontal + 50, vertical + 50}

		return DrawingPoint{
			topLeft,
			bottomRight,
		}
	}

	for _, gridPoint := range identicon.GridPoints {
		drawingPoints = append(drawingPoints, pixelFunc(gridPoint))
	}
	identicon.PixelMap = drawingPoints
	return identicon
}

func rect(img *image.RGBA, col color.Color, x1, y1, x2, y2 float64) {
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(col)
	gc.MoveTo(x1, y1)

	gc.LineTo(x1, y1)
	gc.LineTo(x1, y2)
	gc.MoveTo(x2, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)

	gc.SetLineWidth(0)
	gc.FillStroke()
}

func DrawRectangle(identicon Identicon) error {
	// We create our default image containing a 250x250 rectangle
	var img = image.NewRGBA(image.Rect(0, 0, 250, 250))
	// We retrieve the color from the color property on the identicon
	col := color.RGBA{identicon.Color[0], identicon.Color[1], identicon.Color[2], 255}

	// Loop over the pixelmap and call the rect function with the img, color and the dimensions
	for _, pixel := range identicon.PixelMap {
		rect(
			img,
			col,
			float64(pixel.TopLetft.X),
			float64(pixel.TopLetft.Y),
			float64(pixel.BottomRight.X),
			float64(pixel.BottomRight.Y),
		)
	}
	// Finally save the image to disk
	return draw2dimg.SaveToPngFile(identicon.Title+".png", img)
}