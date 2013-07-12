clipper
======

Clipper is a Go program for generating packing slips and box outside labels in pdf.

Example
-------

```go
package main

import (
	"clipper"
)

func main() {

  from := clipper.Address{
    "Strickland Propane",
    "Propane and Propane Acessories",
    "135 Los Gatos Road", "Arlen, Texas 78701",
    "",
  }

  to := clipper.Address{
		"Mega-Lo-Mart",
		"16 E Grand Ave",
		"Englewood, CO 80113",
		"(303) 761-6474","",
  }

  // specify a purchase order ident
  po := "4501500012"

  // create a new order
  order := clipper.NewOrder(from, to, po)

  var box clipper.Box

  // add a new box to order
  box = clipper.NewBox()
	box.AddItem("Propane", "20 Gal") // specify a box with item Propane with amount 20 gal
	box.SetNote("Do not spill")      // set a note for this box
  order.AddBox(box, 2)             // set the number of boxes of the same.

  box = clipper.NewBox()
	box.AddItem("Propane", "30 Gal")
	box.SetNote("Do not shake")
  order.AddBox(box, 1)

  box = clipper.NewBox()
	box.AddItem("Nozzle", "45")
  order.AddBox(box, 5)

  
	clipper.PackingList(order) // generate packing list
  clipper.Labels(order)      // generate labels
  // files will appear in the working directory
}
```

