package clipper

import (
	"testing"
)

func TestGenerate(t *testing.T) {
  from := Address{
    "Strickland Propane",
    "Propane and Propane Acessories",
    "135 Los Gatos Road", "Arlen, Texas 78701",
    "",
  }

  to := Address{
		"Mega-Lo-Mart",
		"16 E Grand Ave",
		"Englewood, CO 80113",
		"(303) 761-6474","",
  }

  po := "4501500012"
	tsu := "This side up"
  dos := "Do not shake"

  order := NewOrder(from, to, po)
  var box Box

  box = NewBox()
	box.AddItem("Propane", "20 Gal")
	box.SetNote(dos)
  order.AddBox(box, 2)

  box = NewBox()
	box.AddItem("Propane", "30 Gal")
	box.SetNote(tsu)
  order.AddBox(box, 1)

  box = NewBox()
	box.AddItem("Nozzle", "45")
  order.AddBox(box, 5)


	PackingList(order)
  Labels(order)
}
