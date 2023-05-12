package main

import (
	"flag"
	"fmt"
	"image/color"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFOption func(*gofpdf.Fpdf)

func FillColour(c color.RGBA) PDFOption {
	return func(pdf *gofpdf.Fpdf) {
		r, g, b := rgb(c)
		pdf.SetFillColor(r, g, b)
	}
}

func rgb(c color.RGBA) (int, int, int) {
	alpha := float64(c.A) / 255.0
	alphaWhite := int(255 * (1.0 - alpha))
	r := int(float64(c.R)*alpha) + alphaWhite
	g := int(float64(c.G)*alpha) + alphaWhite
	b := int(float64(c.B)*alpha) + alphaWhite
	return r, g, b
}

type PDF struct {
	fpdf *gofpdf.Fpdf
	x, y float64
}

func (p *PDF) Move(xDelta, yDelta float64) {
	p.x, p.y = p.x+xDelta, p.y+yDelta
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) MoveAbs(x, y float64) {
	p.x, p.y = x, y
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) Text(text string) {
	p.fpdf.Text(p.x, p.y, text)
}

func (p *PDF) Polygon(pts []gofpdf.PointType, opts ...PDFOption) {
	for _, opt := range opts {
		opt(p.fpdf)
	}
	p.fpdf.Polygon(pts, "F")
}

func main() {
	name := flag.String("name", "", "the name of the person who completed the course")
	flag.Parse()
	fpdf := gofpdf.New(gofpdf.OrientationLandscape, gofpdf.UnitPoint, gofpdf.PageSizeA4, "")
	w, h := fpdf.GetPageSize()
	fpdf.AddPage()
	pdf := PDF{
		fpdf: fpdf,
	}

	primary := color.RGBA{121, 50, 168, 255}
	secondary := color.RGBA{121, 50, 168, 220}

	// Top banner
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{0, h / 9.0},
		{w, 0},
	}, FillColour(secondary))
	pdf.Polygon([]gofpdf.PointType{
		{w, 0},
		{0, 0},
		{w, h / 9.0},
	}, FillColour(primary))

	// Bottom banner
	pdf.Polygon([]gofpdf.PointType{
		{w, h - h/9.0},
		{w, h},
		{0, h},
	}, FillColour(secondary))
	pdf.Polygon([]gofpdf.PointType{
		{w, h},
		{0, h - h/9.0},
		{0, h},
	}, FillColour(primary))

	// "Certificate of Completion"
	fpdf.SetFont("arial", "B", 50)
	fpdf.SetTextColor(50, 50, 50)
	pdf.MoveAbs(0, 100)
	_, lineHt := fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, "Certificate of Completion", gofpdf.AlignCenter)
	pdf.Move(0, lineHt*2.0)

	fpdf.SetFont("arial", "", 28)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, "This certificate is awarded to", gofpdf.AlignCenter)
	pdf.Move(0, lineHt*2.0)

	fpdf.SetFont("times", "I", 42)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, *name, gofpdf.AlignCenter)
	pdf.Move(0, lineHt*1.75)

	fpdf.SetFont("arial", "", 22)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt*1.5, "For successfully completing all twenty programming exercises in the Gophercises programming course for budding Gophers (Go developers)", gofpdf.AlignCenter)
	pdf.Move(0, lineHt*3.75)

	// Gopher
	fpdf.ImageOptions("images/gophercises.png", w/2.0-50.0, pdf.y, (50.0 * 2), 0, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")
	// gopher, err := gofpdf.SVGBasicFileParse("images/gophercises.svg")
	// if err != nil {
	// 	panic(err)
	// }
	// fpdf.SVGBasicWrite(&gopher, 0.8)

	pdf.Move(0, 65.0)
	fpdf.SetFillColor(100, 100, 100)
	lineLgth := 250.0
	fpdf.Rect(w/4-lineLgth/2, pdf.y, lineLgth, 1.0, "F")
	fpdf.Rect(w/4*3-lineLgth/2, pdf.y, lineLgth, 1.0, "F")

	fpdf.SetFont("arial", "", 12)
	pdf.MoveAbs(0, pdf.y)
	fpdf.SetTextColor(100, 100, 100)
	cellWidth := 200.0
	cellHeight := lineHt * 1.5
	pdf.MoveAbs(w/4-(cellWidth/2), pdf.y)
	fpdf.CellFormat(cellWidth, cellHeight, "Date", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignCenter, false, 0, "")
	pdf.MoveAbs(w/4*3-(cellWidth/2), pdf.y)
	fpdf.CellFormat(cellWidth, cellHeight, "Instructor - Jon Calhoun", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignCenter, false, 0, "")

	fpdf.SetFont("times", "I", 22)
	pdf.MoveAbs(w/4-(cellWidth/2), pdf.y-cellHeight)
	yr, mo, day := time.Now().Date()
	dateStr := fmt.Sprintf("%d/%d/%d", day, mo, yr)
	fpdf.CellFormat(cellWidth, cellHeight, dateStr, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignCenter, false, 0, "")
	pdf.MoveAbs(w/4*3-(cellWidth/2), pdf.y)
	// fpdf.CellFormat(cellWidth, cellHeight, "Jonathan Calhoun", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignCenter, false, 0, "")

	sig, err := gofpdf.SVGBasicFileParse("images/sig.svg")
	if err != nil {
		panic(err)
	}
	pdf.Move(0, -(sig.Ht*0.25 - lineHt))
	fpdf.SVGBasicWrite(&sig, 0.35)

	// Grid
	// drawGrid(pdf)

	err = fpdf.OutputFileAndClose("cert.pdf")
	if err != nil {
		panic(err)
	}
}

// // Function to create grid to help with creating pdf
// func drawGrid(pdf *gofpdf.Fpdf) {
// 	w, h := pdf.GetPageSize()
// 	pdf.SetFont("courier", "", 12)
// 	pdf.SetTextColor(80, 80, 80)
// 	pdf.SetDrawColor(200, 200, 200)

// 	for x := 0.0; x < w; x = x + (w / 20.0) {
// 		pdf.SetTextColor(200, 200, 200)
// 		pdf.Line(x, 0, x, h)
// 		_, lineHt := pdf.GetFontSize()
// 		pdf.Text(x, lineHt, fmt.Sprintf("%d", int(x)))
// 	}

// 	for y := 0.0; y < h; y = y + (w / 20.0) {
// 		if y < bannerHt*0.8 {
// 			pdf.SetTextColor(200, 200, 200)
// 		} else {
// 			pdf.SetTextColor(80, 80, 80)
// 		}
// 		pdf.Line(0, y, w, y)
// 		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
// 	}
// }
