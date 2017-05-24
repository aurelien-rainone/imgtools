package imgscan

import (
	"image"
	"image/color"
	"testing"

	"github.com/aurelien-rainone/imgtools/binimg"
	"github.com/aurelien-rainone/imgtools/internal/test"
)

func newBinaryFromString(ss []string) *binimg.Binary {
	w, h := len(ss[0]), len(ss)
	for i := range ss {
		if len(ss[i]) != w {
			panic("all strings should have the same length")
		}
	}

	img := binimg.New(image.Rect(0, 0, w, h))
	for y := range ss {
		for x := range ss[y] {
			if ss[y][x] == '1' {
				img.SetBit(x, y, binimg.On)
			}
		}
	}
	return img
}

func testIsWhite(t *testing.T, newScanner func(image.Image) Scanner) {
	ss := []string{
		"000",
		"100",
		"011",
	}

	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
	}{
		{0, 0, 3, 3, false},
		{1, 1, 3, 3, false},
		{0, 1, 1, 2, true},
		{0, 0, 1, 1, false},
		{1, 0, 2, 1, false},
		{1, 0, 3, 2, false},
		{1, 2, 3, 3, true},
		{2, 2, 3, 3, true},
	}

	scanner := newScanner(newBinaryFromString(ss))
	for _, tt := range testTbl {
		actual := scanner.IsUniformColor(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy), binimg.On)
		if actual != tt.expected {
			t.Errorf("%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
	}
}

func testIsBlack(t *testing.T, newScanner func(image.Image) Scanner) {
	ss := []string{
		"111",
		"011",
		"100",
	}

	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
	}{
		{0, 0, 3, 3, false},
		{1, 1, 3, 3, false},
		{0, 1, 1, 2, true},
		{0, 0, 1, 1, false},
		{1, 0, 2, 1, false},
		{1, 0, 3, 2, false},
		{1, 2, 3, 3, true},
		{2, 2, 3, 3, true},
	}

	scanner := newScanner(newBinaryFromString(ss))
	for _, tt := range testTbl {
		actual := scanner.IsUniformColor(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy), binimg.Off)
		if actual != tt.expected {
			t.Errorf("(%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
	}
}

func testIsUniform(t *testing.T, newScanner func(image.Image) Scanner) {
	ss := []string{
		"111",
		"011",
		"100",
	}
	var testTbl = []struct {
		minx, miny, maxx, maxy int
		expected               bool
		expectedColor          color.Color
	}{
		{0, 0, 3, 3, false, nil},
		{1, 1, 3, 3, false, nil},
		{0, 1, 1, 2, true, binimg.Off},
		{0, 0, 1, 1, true, binimg.On},
		{1, 0, 2, 1, true, binimg.On},
		{1, 0, 3, 2, true, binimg.On},
		{1, 2, 3, 3, true, binimg.Off},
		{2, 2, 3, 3, true, binimg.Off},
	}

	scanner := newScanner(newBinaryFromString(ss))
	for _, tt := range testTbl {
		actual, color := scanner.IsUniform(image.Rect(tt.minx, tt.miny, tt.maxx, tt.maxy))
		if actual != tt.expected {
			t.Errorf("(%d,%d|%d,%d): expected %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expected, actual)
		}
		if color != tt.expectedColor {
			t.Errorf("(%d,%d|%d,%d): expected color %v, actual %v", tt.minx, tt.miny, tt.maxx, tt.maxy, tt.expectedColor, color)
		}
	}
}

func TestLinesScannerIsWhite(t *testing.T) {
	testIsWhite(t,
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			test.Check(t, err)
			return s
		},
	)
}

func TestLinesScannerIsBlack(t *testing.T) {
	testIsBlack(t,
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			test.Check(t, err)
			return s
		},
	)
}

func TestLinesScannerIsUniform(t *testing.T) {
	testIsUniform(t,
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			test.Check(t, err)
			return s
		},
	)
}

func benchmarkScanner(b *testing.B, pngfile string, newScanner func(image.Image) Scanner) {
	img, err := test.LoadPNG(pngfile)
	test.CheckB(b, err)

	scanner := newScanner(binimg.NewFromImage(img))

	// run N times
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		scanner.IsUniformColor(img.Bounds(), color.White)
		scanner.IsUniformColor(img.Bounds(), color.Black)
		scanner.IsUniform(img.Bounds())
	}
}

func BenchmarkLinesScanner(b *testing.B) {
	benchmarkScanner(b, "./testdata/big.png",
		func(img image.Image) Scanner {
			s, err := NewScanner(img)
			test.CheckB(b, err)
			return s
		})
}
