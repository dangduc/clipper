package clipper

import (
	"bytes"
	"fmt"
	"dox2go"
	"dox2go/pdf"
	"os"
)

type Address struct {
	L1, L2, L3, L4, L5 string
}

type Order struct {
	from, to Address
	po string
	boxes []Box
}

type Box struct {
	cartons int
	items []Item
	note string
}

type Item struct {
	name, count string
}

func NewOrder(f, t Address, p string) Order {
	return Order{from: f, to: t, po: p}
}

func (o *Order) AddBox(b Box, c int) {
	b.SetCartons(c)
	o.boxes = append(o.boxes, b)
}

func (o *Order) TotalItems() int {
	total := 0
	for i := range o.boxes {
		total += len(o.boxes[i].items)
	}
	return total
}

func (o *Order) TotalCartons() int {
	total := 0
	for i := range o.boxes {
		total += o.boxes[i].cartons
	}
	return total
}

func NewBox() Box {
	b := Box{}
	return b
}

func (b *Box) AddItem(name, count string) {
	b.items = append(b.items, Item{name, count})
}

func (b *Box) SetNote(n string) {
	b.note = n
}

func (b *Box) SetCartons(num int) {
	b.cartons = num
}

func PackingList(order Order) {

	total_pages := order.TotalItems() / 16
	if order.TotalItems() % 16 > 0 {
		total_pages++
	}
	current_page := 0

	// Write the document to this buffer.
	var b bytes.Buffer

	// Create a Pdf Document instance.
	doc := pdf.NewPdfDoc(&b)

	// The document object manages font instances.
	font6 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 6)
	font10 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 10)
	font12 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 12)

	pWidth, pHeight := dox2go.StandardSize(dox2go.PS_Letter, dox2go.U_MM) //Letter Page

  // -----------------------------------------------------------------

	lmargin := 15.0
	rmargin := 200.0
	bmargin := 20.0

	ulineh := 2.0 // underline offset

	var plisth, toh, tohoff, poh, item_start_h, ihoff float64
	var itemFont dox2go.Font

	switch n := order.TotalItems(); {
	case n <= 0:
		return
  case n <= 5:
		plisth = 235.0
		toh = 200.0
		tohoff = 15.0
		poh = 165.0
		item_start_h = 140.0
		ihoff = 20.0
		itemFont = font12 
  case n <= 10:
		plisth = 235.0
		toh = 220.0
		tohoff = 15.0
		poh = 190.0
		item_start_h = 170.0
		ihoff = 15.0
		itemFont = font12 
  case n > 10:
		plisth = 237.0
		toh = 222.0
		tohoff = 10.0
		poh = 195.0
		item_start_h = 182.0
		ihoff = 10.0
		itemFont = font6
	}

	var ih float64
	var itmpath *dox2go.Path
	var ibuffer []Item
	cnum := 0
	inum := 0
	items_per_page_limit := 16
	var page dox2go.Page
	var s dox2go.Surface
	for i := range order.boxes {
		if inum % items_per_page_limit == 0 {
		// Add a page
			page = doc.CreatePage(dox2go.U_MM, // Work in millimeters
			pWidth, pHeight,
			dox2go.PO_Portrait)               // Portrait orientation

			// Get the drawing surface of the page.
			s = page.Surface()

			s.Bg(dox2go.RGB(0, 0, 0)) // Text is drawn using the background colour

			fcw := 60.0
			fch := 260.0	
			s.Text(font6, fcw, fch, order.from.L1)
			s.Text(font6, 0.0, -6.0, order.from.L2)
			s.Text(font6, 0.0, -6.0, order.from.L3)
			fcpath := dox2go.NewPath()
			fcpath.Move(lmargin, 245.0)
			fcpath.Line(rmargin, 245.0)
			s.Stroke(fcpath)
			s.PushState()

			s.Text(font6, fcw, plisth - ulineh, "Packing List")
			plpath := dox2go.NewPath()
			plpath.Move(fcw - 2.0, plisth - ulineh)
			plpath.Line(fcw + 37.0, plisth - ulineh)
			s.Stroke(plpath)
			s.PushState()


			s.Text(font6, lmargin, toh, "To:")
			s.PushState()

			s.Text(font6, lmargin + 15.0, toh, order.to.L1)
			s.PushState()

			tcpath := dox2go.NewPath()
			tcpath.Move(lmargin + 10.0, toh - ulineh)
			tcpath.Line(rmargin, toh - ulineh)
			s.Stroke(tcpath)
			tc2path := dox2go.NewPath()
			tc2path.Move(lmargin + 10.0, toh - tohoff - ulineh)
			tc2path.Line(rmargin, toh - tohoff - ulineh)
			s.Stroke(tc2path)
			s.PushState()

			s.Text(font6, lmargin + 15.0, poh, "P.O.#:")
			s.Text(font6, 25.0, 0.0, order.po)
			popath := dox2go.NewPath()
			popath.Move(lmargin + 35.0, poh - ulineh)
			popath.Line(lmargin + 160.0, poh - ulineh)
			s.Stroke(popath)
			s.PushState()

			current_page++
			s.Text(font10, lmargin + 50.0, bmargin, fmt.Sprintf("%d", order.TotalCartons()))
			s.PushState()
			s.Text(font10, lmargin + 107.0, bmargin, fmt.Sprintf("%d", current_page))
			s.PushState()
			s.Text(font10, lmargin + 140.0, bmargin, fmt.Sprintf("%d", total_pages))
			s.PushState()



			s.Text(font6, lmargin, bmargin, "Total carton(s)")
			s.Text(font6, 85.0, 0, "Page")
			s.Text(font6, 40.0, 0, "of")
			s.Text(font6, 35.0, 0,"pages(s)")
			pgpath := dox2go.NewPath()
			pgpath.Move(lmargin + 42.0, bmargin - ulineh)
			pgpath.Line(lmargin + 80.0, bmargin - ulineh)
			s.Stroke(pgpath)
			pgpath = dox2go.NewPath()
			pgpath.Move(lmargin + 100.0, bmargin - ulineh)
			pgpath.Line(lmargin + 123.0, bmargin - ulineh)
			s.Stroke(pgpath)
			pgpath = dox2go.NewPath()
			pgpath.Move(lmargin + 132.0, bmargin - ulineh)
			pgpath.Line(lmargin + 158.0, bmargin - ulineh)
			s.Stroke(pgpath)
			s.PushState()

			ih = item_start_h
		}

		if len(ibuffer) > 0 {
			bindex := 0
			for c, item := range ibuffer {
				if c != 0 && inum != 0 && inum % items_per_page_limit == 0 {
					ibuffer = ibuffer[c:]
					break;
				}
				//print item
				s.Text(itemFont, lmargin + 72.0, ih, item.name)
				s.PushState()
				s.Text(itemFont, lmargin + 132.0, ih, item.count)
				s.PushState()

				//print template
				s.Text(font6, lmargin, ih, "Ct #:")
				s.Text(font6, 50.0, 0, "Item #:")
				s.Text(font10, 70.0, 0, "@")
				s.Text(font6, 60.0, 0, "ea")
				itmpath = dox2go.NewPath()
				itmpath.Move(lmargin + 15.0, ih - ulineh)
				itmpath.Line(lmargin + 45.0, ih - ulineh)
				s.Stroke(itmpath)
				itmpath = dox2go.NewPath()
				itmpath.Move(lmargin + 70.0, ih - ulineh)
				itmpath.Line(lmargin + 115.0, ih - ulineh)
				s.Stroke(itmpath)
				itmpath = dox2go.NewPath()
				itmpath.Move(lmargin + 130.0, ih - ulineh)
				itmpath.Line(lmargin + 175.0, ih - ulineh)
				s.Stroke(itmpath)
				s.PushState()

				ih -= ihoff
				inum++
				bindex++
			}
			if bindex == len(ibuffer) {
				ibuffer = []Item{} //empty ibuffer
			}

		}

		if len(ibuffer) == 0 {

			if 1 == order.boxes[i].cartons {
				cnum++
				s.Text(itemFont, lmargin + 20.0, ih, fmt.Sprintf("%d", cnum))
			} else {
				s.Text(itemFont, lmargin + 15.0, ih, fmt.Sprintf("%d-%d", cnum + 1, cnum + order.boxes[i].cartons))
				cnum += order.boxes[i].cartons
			}
			s.PushState()

			for c, item := range order.boxes[i].items {

				if c != 0 && inum != 0 && inum % items_per_page_limit == 0 {
					ibuffer = order.boxes[i].items[c:]
					break;
				}
				//print item
				s.Text(itemFont, lmargin + 72.0, ih, item.name)
				s.PushState()
				s.Text(itemFont, lmargin + 132.0, ih, item.count)
				s.PushState()

				//print template
				s.Text(font6, lmargin, ih, "Ct #:")
				s.Text(font6, 50.0, 0, "Item #:")
				s.Text(font10, 70.0, 0, "@")
				s.Text(font6, 60.0, 0, "ea")
				itmpath = dox2go.NewPath()
				itmpath.Move(lmargin + 15.0, ih - ulineh)
				itmpath.Line(lmargin + 45.0, ih - ulineh)
				s.Stroke(itmpath)
				itmpath = dox2go.NewPath()
				itmpath.Move(lmargin + 70.0, ih - ulineh)
				itmpath.Line(lmargin + 115.0, ih - ulineh)
				s.Stroke(itmpath)
				itmpath = dox2go.NewPath()
				itmpath.Move(lmargin + 130.0, ih - ulineh)
				itmpath.Line(lmargin + 175.0, ih - ulineh)
				s.Stroke(itmpath)
				s.PushState()

				ih -= ihoff
				inum++
			}
		}
	}

	// We're finished so flush the document to the buffer.
	doc.Close()

	// Write it to a file.
	f, err := os.Create("list.pdf")
	if err != nil {
		fmt.Println("Could not create list.pdf")
		return
	}
	n, err := b.WriteTo(f)
	fmt.Printf("Wrote %d bytes\n", n)

}

func Labels(order Order) {

	// Write the document to this buffer.
	var b bytes.Buffer

	// Create a Pdf Document instance.
	doc := pdf.NewPdfDoc(&b)

	// The document object manages font instances.
	font5 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 5)
	font6 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 6)
	font8 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 8)
	font10 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 10)
	font12 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 12)
	font16 := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 16)
	fontBig := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 24)

	pWidth, pHeight := dox2go.StandardSize(dox2go.PS_Letter, dox2go.U_MM) //Letter Page
	lWidth := pHeight
	lHeight := pWidth // Landscape sizing

  // -----------------------------------------------------------------

	lmargin := 15.0
	rmargin := 260.0
	bmargin := 15.0


	ulineh := 2.0 // underline offset

	var toh, tohoff, poh float64

	toh = 160.0
	tohoff = 12.0
	poh = 95.0



	var page dox2go.Page
	var s dox2go.Surface
	cnum := 0
	for b := range order.boxes {
		for d := 0 ; d < order.boxes[b].cartons ; d++ {
			cnum++
			// Add a page
			page = doc.CreatePage(dox2go.U_MM, // Work in millimeters
			lWidth, lHeight,
			dox2go.PO_Portrait)               // Portrait orientation

			// Get the drawing surface of the page.
			s = page.Surface()

			s.Bg(dox2go.RGB(0, 0, 0)) // Text is drawn using the background colour

			fcw := 100.0
			fch := 200.0	
			fclineh := 190.0
			s.Text(font6, fcw, fch, order.from.L1)
			s.Text(font6, 0.0, -6.0, order.from.L2)
			if order.from.L3 != "" {
				s.Text(font6, 0.0, -6.0, order.from.L3)
				fclineh -= 6.0
			}

			if order.from.L4 != "" {
				s.Text(font6, 0.0, -6.0, order.from.L4)
				fclineh -= 6.0
			}

			if order.from.L5 != "" {
				s.Text(font6, 0.0, -6.0, order.from.L5)
				fclineh -= 6.0
			}
			fcpath := dox2go.NewPath()
			fcpath.Move(lmargin, fclineh)
			fcpath.Line(rmargin, fclineh)
			s.Stroke(fcpath)
			s.PushState()

			s.Text(font12, lmargin, toh, "To:")
			s.PushState()


			tclineh := 118.0
			s.Text(font10, lmargin + 25.0, toh, order.to.L1)
			s.Text(font10, 0.0, -tohoff, order.to.L2)
			s.Text(font10, 0.0, -tohoff, order.to.L3)
			if order.to.L4 != "" {
				s.Text(font10, 0.0, -tohoff, order.to.L4)
			}
			if order.to.L5 != "" {
				s.Text(font10, 0.0, -tohoff, order.to.L5)
				tclineh -= tohoff
			}
			tcpath := dox2go.NewPath()
			tcpath.Move(lmargin, tclineh)
			tcpath.Line(rmargin, tclineh)
			s.Stroke(tcpath)
			s.PushState()


			s.Text(font12, lmargin, poh, "PO #:")
			s.Text(font12, 40.0, 0.0, order.po)
			popath := dox2go.NewPath()
			popath.Move(lmargin + 35.0, poh - ulineh)
			popath.Line(lmargin + 160.0, poh - ulineh)
			s.Stroke(popath)
			s.PushState()

			noteh := poh + 12.0
			s.Text(font12, lmargin + 40.0, noteh, order.boxes[b].note)
			s.PushState()

			imh := 40.0
			s.Text(font8, lmargin, imh, "Item")
			s.PushState()
			s.Text(font8, lmargin + 18.0, imh, "#")
			s.Text(font8, 105.0, 0.0, "@")
			s.Text(font8, 105.0, 0.0, "ea.")
			impath := dox2go.NewPath()
			impath.Move(lmargin + 25.0, imh - ulineh)
			impath.Line(lmargin + 120.0, imh - ulineh)
			s.Stroke(impath)
			impath = dox2go.NewPath()
			impath.Move(lmargin + 135.0, imh - ulineh)
			impath.Line(lmargin + 225.0, imh - ulineh)
			s.Stroke(impath)
			s.PushState()


			s.Text(font8, lmargin + 18.0, bmargin, "#:")
			s.Text(font8, 105.0, 0.0, "of")
			s.Text(font8, 105.0, 0.0, "Ct(s)")
			tbpath := dox2go.NewPath()
			tbpath.Move(lmargin + 25.0, bmargin - ulineh)
			tbpath.Line(lmargin + 120.0, bmargin - ulineh)
			s.Stroke(tbpath)
			tbpath = dox2go.NewPath()
			tbpath.Move(lmargin + 135.0, bmargin - ulineh)
			tbpath.Line(lmargin + 225.0, bmargin - ulineh)
			s.Stroke(tbpath)
			s.PushState()

			s.Text(fontBig, lmargin + 60.0, bmargin, fmt.Sprintf("%d", cnum))
			s.PushState()
			s.Text(fontBig, lmargin + 165.0, bmargin, fmt.Sprintf("%d", order.TotalCartons()))
			s.PushState()


			ihoff := 0.0
			var iwoff float64
			var ihdelta float64
			var itemFont dox2go.Font
			switch n := len(order.boxes[b].items); {
			case n <= 0:
				return
		  case n <= 2:
		  	itemFont = fontBig
		  	ihdelta = 26.0
		  case n <= 3:
		  	itemFont = font16
		  	ihdelta = 15.0
		  	iwoff = 20.0
		  case n <= 5:
		  	itemFont = font10
		  	ihdelta = 11.0
		  	iwoff = 25.0
		  case n <= 8:
		  	itemFont = font6
		  	ihdelta = 6.8
		  	iwoff = 30.0
		  case n <= 10:
		  	itemFont = font5
		  	ihdelta = 5.4
		  	iwoff = 40.0
			}
			for _, item := range order.boxes[b].items {
				s.Text(itemFont, lmargin + 25.0 + iwoff, imh + ihoff, item.name)
				s.Text(itemFont, 115.0, 0.0, item.count)
				s.PushState()
				ihoff += ihdelta
			}
		}
	}

	// We're finished so flush the document to the buffer.
	doc.Close()

	// Write it to a file.
	f, err := os.Create("labels.pdf")
	if err != nil {
		fmt.Println("Could not create labels.pdf")
		return
	}
	n, err := b.WriteTo(f)
	fmt.Printf("Wrote %d bytes\n", n)


}
