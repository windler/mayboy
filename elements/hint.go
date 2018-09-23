package elements

import (
	"github.com/rivo/tview"
	"github.com/windler/mayboy/events"
)

//Hint shows key hints
type Hint struct {
	hint *tview.TextView
	em   *events.EventManager
}

//NewHint creates a new Hint
func NewHint(em *events.EventManager) Hint {
	hint := tview.NewTextView()
	hint.SetBorderPadding(1, 0, 0, 0)

	h := Hint{
		hint: hint,
		em:   em,
	}
	h.registerListeners()

	return h
}

func (h *Hint) registerListeners() {
	h.em.Listen(events.ProjectSelected, func() {
		h.hint.SetText("Press ESC to exit")
	})

	h.em.Listen(events.IssueTableFocusLost, func() {
		h.hint.SetText("")
	})
}

//GetPrimitive returns the rivo/tview primtive
func (h *Hint) GetPrimitive() tview.Primitive {
	return h.hint
}
