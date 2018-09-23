package elements

import "github.com/rivo/tview"

//Grid structures primitves for ui
type Grid struct {
	grid *tview.Grid
}

//CreateGrid creates a new Grid
func CreateGrid(header, mainLeft, mainRight, mainBottomRight, footer tview.Primitive) Grid {
	g := Grid{}

	g.grid = tview.NewGrid().
		SetRows(1, 0, 2, 1).
		SetColumns(30, 0)

	g.grid.AddItem(header, 0, 0, 1, 2, 0, 100, false)

	g.grid.AddItem(mainLeft, 1, 0, 1, 1, 0, 100, true)
	g.grid.AddItem(mainRight, 1, 1, 1, 1, 0, 100, false)
	g.grid.AddItem(tview.NewTextView(), 2, 0, 1, 1, 0, 100, false)
	g.grid.AddItem(mainBottomRight, 2, 1, 1, 1, 0, 100, false)

	g.grid.AddItem(footer, 3, 0, 1, 2, 0, 100, false)

	return g
}

//GetPrimitive returns the rivo/tview primtive
func (g *Grid) GetPrimitive() tview.Primitive {
	return g.grid
}
