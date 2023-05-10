package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

const (
	bannerHt = 95.0
	xIndent  = 40.0
)

func main() {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeA4, "")
	w, h := pdf.GetPageSize()
	fmt.Printf("width=%v, height=%v\n", w, h)
	pdf.AddPage()

	// Create top and bottom banners
	pdf.SetFillColor(103, 60, 79)
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{w, 0},
		{w, bannerHt},
		{0, bannerHt * 0.8},
	}, "F")
	pdf.Polygon([]gofpdf.PointType{
		{0, h},
		{0, h - (bannerHt * 0.2)},
		{w, h - (bannerHt * 0.1)},
		{w, h},
	}, "F")

	// "Invoice"
	pdf.SetFont("arial", "", 40)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt := pdf.GetFontSize()
	pdf.Text(xIndent, bannerHt-lineHt, "Invoice")

	// Banner - phone, email, domain
	pdf.SetFont("arial", "", 8)
	pdf.SetTextColor(255,255,255)
	_, lineHt = pdf.GetFontSize()
	pdf.MoveTo(w - xIndent - 2.0*124.0, (bannerHt - (lineHt*1.5*3.0))/2.0)
	pdf.MultiCell(124.0, lineHt * 1.5, "0151 336 4700\nross1012@gmail.com\ngithub.com/ematogra", gofpdf.BorderNone, gofpdf.AlignRight, false)

	// Banner - address
	pdf.SetFont("arial", "", 8)
	pdf.SetTextColor(255,255,255)
	_, lineHt = pdf.GetFontSize()
	pdf.MoveTo(w - xIndent - 124.0, (bannerHt - (lineHt*1.5*3.0))/2.0)
	pdf.MultiCell(124.0, lineHt * 1.5, "Flat 7, 14-16 Underhill Road\nLondon\nSE22 0AH", gofpdf.BorderNone, gofpdf.AlignRight, false)

	// Grid
	// drawGrid(pdf)

	err := pdf.OutputFileAndClose("p2.pdf")
	if err != nil {
		panic(err)
	}
}

// Function to create grid to help with creating pdf
func drawGrid(pdf *gofpdf.Fpdf) {
	w, h := pdf.GetPageSize()
	pdf.SetFont("courier", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetDrawColor(200, 200, 200)

	for x := 0.0; x < w; x = x + (w / 20.0) {
		pdf.SetTextColor(200, 200, 200)
		pdf.Line(x, 0, x, h)
		_, lineHt := pdf.GetFontSize()
		pdf.Text(x, lineHt, fmt.Sprintf("%d", int(x)))
	}

	for y := 0.0; y < h; y = y + (w / 20.0) {
		if y < bannerHt*0.8 {
			pdf.SetTextColor(200, 200, 200)
		} else {
			pdf.SetTextColor(80, 80, 80)
		}
		pdf.Line(0, y, w, y)
		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
	}
}
