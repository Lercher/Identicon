package identicon

import (
	"crypto/md5"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/llgcode/draw2d/draw2dimg"
)

// Point is 2d int/F5 coordinates
type Point struct {
	X, Y int
}

type gridPoint struct {
	value byte
	index int
}

// Identicon grphically represents a hash in a 5x5 matrix
type Identicon struct {
	hash       [16]byte
	color      [3]byte
	grid       []byte
	gridPoints []gridPoint
	Pixels     []Point
}

// Generate creates an Identicon from an arbitrary byte array.
// It is garanteed that the same byte array produces the same Identicon.
func Generate(input []byte) Identicon {
	identicon := hashInput(input)
	identicon = pickColor(identicon)
	identicon = buildGrid(identicon)
	identicon = filterOddSquares(identicon)
	identicon = buildPixelMap(identicon)
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
	for i := 0; i+3 < len(identicon.hash); i += 3 {
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

func buildPixelMap(identicon Identicon) Identicon {
	var points []Point

	pixelFunc := func(p gridPoint) Point {
		horizontal := (p.index % 5)
		vertical := (p.index / 5)
		return Point{X: horizontal, Y: vertical}
	}

	for _, gridPoint := range identicon.gridPoints {
		points = append(points, pixelFunc(gridPoint))
	}
	identicon.Pixels = points
	return identicon
}

// LightBackground is either true or false
type LightBackground bool

// WritePNGImage writes the identicon image to the given writer with width and height of 5 times pixwidth
// and LightBackground(true) or LightBackground(false) as lbg
func (i Identicon) WritePNGImage(w io.Writer, pixwidth int, lbg LightBackground) error {
	var img = image.NewRGBA(image.Rect(0, 0, pixwidth*5, pixwidth*5))
	col := color.YCbCr{Y: i.color[0], Cb: i.color[1], Cr: i.color[2]}
	if bool(lbg) {
		col.Y &= 0x7f // reset high bit in luma
	} else {
		col.Y |= 0x80 // set high bit in luma
	}

	gc := draw2dimg.NewGraphicContext(img)
	for _, pixel := range i.Pixels {
		rect(
			gc,
			col,
			float64(pixel.X*pixwidth), float64(pixel.Y*pixwidth),
			float64(pixel.X*pixwidth+pixwidth), float64(pixel.Y*pixwidth+pixwidth),
		)
	}
	gc.Close()
	return png.Encode(w, img)
}

func rect(gc *draw2dimg.GraphicContext, col color.Color, x1, y1, x2, y2 float64) {
	gc.SetFillColor(col)
	gc.SetLineWidth(0)
	gc.MoveTo(x1, y1)
	gc.LineTo(x1, y2)
	gc.LineTo(x2, y2)
	gc.LineTo(x2, y1)
	gc.LineTo(x1, y1)
	gc.FillStroke()
}
