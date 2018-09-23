package elements

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

//Header shows app info
type Header struct {
	header *tview.TextView
}

//NewHeader creates a new Header
func NewHeader() Header {
	header := tview.NewTextView()
	header.SetBackgroundColor(tcell.ColorLightGray)
	header.SetTextColor(tcell.ColorBlack)
	header.SetText("mayboy - gitlab issue viewer")
	header.SetTextAlign(tview.AlignCenter)

	h := Header{
		header: header,
	}

	return h
}

//GetPrimitive returns the rivo/tview primtive
func (h *Header) GetPrimitive() tview.Primitive {
	return h.header
}
