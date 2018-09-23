package elements

import "github.com/rivo/tview"

//Grid structures primitves for ui
type Grid struct {
	grid *tview.Grid
}

//CreateGrid creates a new Grid
func CreateGrid(mainLeft, mainRight, footerLeft, footerRight tview.Primitive) Grid {
	g := Grid{}
	g.grid = tview.NewGrid().
		SetRows(0, 2).
		SetColumns(30, 0)

	g.grid.AddItem(mainLeft, 0, 0, 1, 1, 0, 100, true)
	g.grid.AddItem(mainRight, 0, 1, 1, 1, 0, 100, false)
	g.grid.AddItem(footerLeft, 1, 0, 1, 1, 0, 100, false)
	g.grid.AddItem(footerRight, 1, 1, 1, 1, 0, 100, false)

	return g
}

//GetPrimitive returns the rivo/tview primtive
func (g *Grid) GetPrimitive() tview.Primitive {
	return g.grid
}
