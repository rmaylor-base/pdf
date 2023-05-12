package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

const (
	bannerHt = 95.0
	xIndent  = 40.0
	taxRate  = 0.09
)

type LineItem struct {
	UnitName       string
	PricePerUnit   int
	UnitsPurchased int
}

func main() {
	lineItems := []LineItem{
		{
			UnitName:       "Treated Sawn Timber - 19 x 38 x 1800mm",
			PricePerUnit:   375, // in pence
			UnitsPurchased: 220,
		}, {
			UnitName:       "Plasterboard Sheet",
			PricePerUnit:   822, // in pence
			UnitsPurchased: 50,
		}, {
			UnitName:       "Paint",
			PricePerUnit:   1455, // in pence
			UnitsPurchased: 3,
		}, {
			UnitName:       "Is the word wrapping working? Not sure it is! let's see... hmmm....",
			PricePerUnit:   1455, // in pence
			UnitsPurchased: 3,
		}, {
			UnitName:       "Nails",
			PricePerUnit:   2345, // in pence
			UnitsPurchased: 4,
		},
	}

	subtotal := 0
	for _, li := range lineItems {
		subtotal += li.PricePerUnit * li.UnitsPurchased
	}
	tax := int(float64(subtotal) * taxRate)
	total := subtotal + tax
	totalGBP := toGBP(total)

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

	// Banner - "Invoice"
	pdf.SetFont("arial", "", 40)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt := pdf.GetFontSize()
	pdf.Text(xIndent, bannerHt-lineHt, "Invoice")

	// Banner - Logo
	pdf.ImageOptions("images/gophercises.png", 238.0, (bannerHt-(bannerHt/1.5))/2.0, 0, bannerHt/1.5, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")

	// Banner - phone, email, domain
	pdf.SetFont("arial", "", 8)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-2.0*124.0, (bannerHt-(lineHt*1.5*3.0))/2.0)
	pdf.MultiCell(124.0, lineHt*1.5, "0151 336 4700\nross1012@gmail.com\ngithub.com/ematogra", gofpdf.BorderNone, gofpdf.AlignRight, false)

	// Banner - address
	pdf.SetFont("arial", "", 8)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-124.0, (bannerHt-(lineHt*1.5*3.0))/2.0)
	pdf.MultiCell(124.0, lineHt*1.5, "Flat 7, 14-16 Underhill Road\nLondon\nSE22 0AH", gofpdf.BorderNone, gofpdf.AlignRight, false)

	// Summary - Billed To, Invoice #, Date of Issue
	_, sy := summaryBlock(pdf, xIndent, bannerHt+lineHt*2.0, "Billed To", "Client Name", "123 Client Street", "City County UK", "Postcode")
	summaryBlock(pdf, xIndent*2.0+lineHt*18.0, bannerHt+lineHt*2.0, "Invoice No", "000000001234")
	summaryBlock(pdf, xIndent*2.0+lineHt*18.0, bannerHt+lineHt*8.5, "Date of Issue", "11/05/2023")

	// Summary - Invoice total
	x, y := w-xIndent-124.0, bannerHt+lineHt*2.35
	pdf.MoveTo(x, y)
	pdf.SetFont("arial", "", 14)
	_, lineHt = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.CellFormat(124.0, lineHt, "Invoice Total", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x, y = x+2.0, y+lineHt*1.5
	pdf.MoveTo(x, y)
	pdf.SetFont("arial", "", 40)
	_, lineHt = pdf.GetFontSize()
	alpha := 58
	pdf.SetTextColor(72+alpha, 42+alpha, 55+alpha)
	pdf.CellFormat(124.0, lineHt, totalGBP, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x, y = x-2.0, y+lineHt*1.25

	// Divider
	if sy > y {
		y = sy
	}
	x, y = xIndent-20.0, y+30.0
	pdf.Rect(x, y, w-(xIndent*2.0)+40.0, 3.0, "F")

	// Line items - titles
	pdf.SetFont("arial", "", 14)
	_, lineHt = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	x, y = xIndent-2.0, y+lineHt
	pdf.MoveTo(x, y)
	pdf.CellFormat(w/2.65+1.8, lineHt, "Description", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100.0, lineHt, "Price per unit", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, "Quantity", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	pdf.CellFormat(119.5, lineHt, "Amount", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")

	// Line items - individual purchases
	y = y + lineHt
	for _, li := range lineItems {
		x, y = lineItem(pdf, x, y, li)
	}

	// Subtotal
	x, y = 364.13, y+lineHt*2.25
	x, y = calculateTotal(pdf, x, y, "Subtotal", subtotal)
	x, y = calculateTotal(pdf, x, y, "Tax", tax)
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(x, y, w-xIndent+10.0, y)
	y = y + lineHt/2
	x, y = calculateTotal(pdf, x, y, "Total", total)

	// Grid
	// drawGrid(pdf)

	err := pdf.OutputFileAndClose("p4.pdf")
	if err != nil {
		panic(err)
	}
}

func calculateTotal(pdf *gofpdf.Fpdf, x, y float64, label string, amount int) (float64, float64) {
	origX := x
	pdf.SetFont("arial", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, label, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(119.5, lineHt, toGBP(amount), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	y = y + lineHt*1.5
	return origX, y
}

func toGBP(pence int) string {
	penceStr := fmt.Sprintf("%d", pence%100)
	if len(penceStr) < 2 {
		penceStr = "0" + penceStr
	}
	return fmt.Sprintf("\u00a3%d.%s", pence/100, penceStr)
}

func lineItem(pdf *gofpdf.Fpdf, x, y float64, lineItem LineItem) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()
	pdf.SetFont("arial", "", 14)
	pdf.SetTextColor(50, 50, 50)
	_, lineHt := pdf.GetFontSize()
	x, y = xIndent-2.0, y+lineHt*.75
	pdf.MoveTo(x, y)
	pdf.MultiCell(w/2.65+1.8, lineHt, lineItem.UnitName, gofpdf.BorderNone, gofpdf.AlignLeft, false)
	tmp := pdf.SplitLines([]byte(lineItem.UnitName), w/2.65+1.5)
	maxY := y + float64(len(tmp)-1)*lineHt
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100.0, lineHt, toGBP(lineItem.PricePerUnit), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, fmt.Sprint(lineItem.UnitsPurchased), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	// x = w - xIndent - 2.0 - 119.5
	pdf.CellFormat(119.5, lineHt, toGBP(lineItem.PricePerUnit*lineItem.UnitsPurchased), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	if maxY > y {
		y = maxY
	}
	y = y + lineHt*1.75
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(xIndent-10.0, y, w-xIndent+10.0, y)
	return origX, y
}

func summaryBlock(pdf *gofpdf.Fpdf, x, y float64, title string, data ...string) (float64, float64) {
	pdf.SetFont("arial", "", 14)
	pdf.SetTextColor(180, 180, 180)
	_, lineHt := pdf.GetFontSize()
	y = y + lineHt
	pdf.Text(x, y, title)
	y = y + lineHt*0.25
	pdf.SetTextColor(50, 50, 50)
	for _, str := range data {
		y = y + lineHt*1.25
		pdf.Text(x, y, str)
	}
	return x, y
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
