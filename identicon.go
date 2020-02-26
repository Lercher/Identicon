package identicon

import (
	"crypto/md5"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/llgcode/draw2d/draw2dimg"
)

type point struct {
	x, y int
}

type drawingPoint struct {
	topLeft     point
	bottomRight point
}

type gridPoint struct {
	value byte
	index int
}

type Identicon struct {
	hash       [16]byte
	color      [3]byte
	grid       []byte
	gridPoints []gridPoint
	pixelMap   []drawingPoint
}

type LightBackground bool

// WritePNGImage writes the identicon image to the given writer with width and height of 5 times pixwidth
// and LightBackground(true) or LightBackground(false) as lbg
func (i Identicon) WritePNGImage(w io.Writer, pixwidth int, lbg LightBackground) error {
	var img = image.NewRGBA(image.Rect(0, 0, pixwidth*5, pixwidth*5))
	var col color.RGBA
	if bool(lbg) {
		col = color.RGBA{R: i.color[0] & 0x8f, G: i.color[1] & 0x8f, B: i.color[2] & 0x8f, A: 255}
	} else {
		col = color.RGBA{R: i.color[0] | 0x80, G: i.color[1] | 0x80, B: i.color[2] | 0x80, A: 255}
	}

	for _, pixel := range i.pixelMap {
		rect(
			img, 
			col, 
			float64(pixel.topLeft.x * pixwidth), float64(pixel.topLeft.y * pixwidth), 
			float64(pixel.bottomRight.x* pixwidth), float64(pixel.bottomRight.y* pixwidth),
		)
	}

	return png.Encode(w, img)
}

type applyFunc func(Identicon) Identicon

// Generate creates an Identicon from an arbitrary byte array.
// It is garanteed that the same byte array produces the same Identicon.
func Generate(input []byte) Identicon {
	identiconPipe := []applyFunc{
		pickColor, buildGrid, filterOddSquares, buildPixelMap,
	}
	identicon := hashInput(input)
	for _, applyFunc := range identiconPipe {
		identicon = applyFunc(identicon)
	}
	return identicon
}

func hashInput(input []byte) Identicon {
	checkSum := md5.Sum(input)
	return Identicon{
		hash: checkSum,
	}
}

func pickColor(identicon Identicon) Identicon {
	rgb := [3]byte{}
	copy(rgb[:], identicon.hash[:3])
	identicon.color = rgb
	return identicon
}

func buildGrid(identicon Identicon) Identicon {
	var grid []byte
	for i := 0; i < len(identicon.hash) && i+3 <= len(identicon.hash)-1; i += 3 {
		chunk := make([]byte, 5)
		copy(chunk, identicon.hash[i:i+3])
		chunk[3] = chunk[1]
		chunk[4] = chunk[0]
		grid = append(grid, chunk...)

	}
	identicon.grid = grid
	return identicon
}

func filterOddSquares(identicon Identicon) Identicon {
	var grid []gridPoint
	for i, code := range identicon.grid {
		if code%2 == 0 {
			point := gridPoint{
				value: code,
				index: i,
			}
			grid = append(grid, point)
		}
	}
	identicon.gridPoints = grid
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

func buildPixelMap(identicon Identicon) Identicon {
	var drawingPoints []drawingPoint

	pixelFunc := func(p gridPoint) drawingPoint {
		horizontal := (p.index % 5)
		vertical := (p.index / 5)
		topLeft := point{x: horizontal, y: vertical}
		bottomRight := point{x: horizontal + 1, y: vertical + 1}

		return drawingPoint{
			topLeft,
			bottomRight,
		}
	}

	for _, gridPoint := range identicon.gridPoints {
		drawingPoints = append(drawingPoints, pixelFunc(gridPoint))
	}
	identicon.pixelMap = drawingPoints
	return identicon
}
