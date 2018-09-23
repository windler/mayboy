package elements

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/windler/mayboy/events"
)

//Footer shows key hints
type Footer struct {
	footer *tview.TextView
	em     *events.EventManager
}

//NewFooter creates a new Header
func NewFooter(em *events.EventManager) Footer {
	footer := tview.NewTextView()
	footer.SetBackgroundColor(tcell.ColorLightGray)
	footer.SetTextColor(tcell.ColorBlack)
	footer.SetTextAlign(tview.AlignRight)

	f := Footer{
		footer: footer,
		em:     em,
	}

	f.registerListeners()
	f.showProjectListText()

	return f
}

func (f *Footer) showProjectListText() {
	f.footer.SetText(getCommandString([]cmd{
		cmd{
			cmd:  "UP",
			desc: "Move up",
		},
		cmd{
			cmd:  "DOWN",
			desc: "Move down",
		},
		cmd{
			cmd:  "ENTER",
			desc: "Select project",
		},
		cmd{
			cmd:  "<n>",
			desc: "Select project (n)",
		},
		cmd{
			cmd:  "ESC",
			desc: "Quit application",
		},
	}))
}

func (f *Footer) registerListeners() {
	f.em.Listen(events.IssueTableFocusLost, f.showProjectListText)
	f.em.Listen(events.ProjectSelected, func() {
		f.footer.SetText(getCommandString([]cmd{
			cmd{
				cmd:  "UP",
				desc: "Move up",
			},
			cmd{
				cmd:  "DOWN",
				desc: "Move down",
			},
			cmd{
				cmd:  "ESC",
				desc: "Go back to project list",
			},
		}))
	})
}

type cmd struct {
	cmd, desc string
}

func getCommandString(commands []cmd) string {
	result := ""
	for _, c := range commands {
		result += fmt.Sprintf("%s: %s\t", c.cmd, c.desc)
	}
	return result
}

//GetPrimitive returns the rivo/tview primtive
func (f *Footer) GetPrimitive() tview.Primitive {
	return f.footer
}
