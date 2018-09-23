package elements

import (
	"github.com/rivo/tview"
	"github.com/windler/mayboy/events"
	"github.com/windler/mayboy/gitlab"
)

//Footer show contextual information
type Footer struct {
	footer        *tview.TextView
	em            *events.EventManager
	issueRetriver SelectedIssueRetriever
}

//SelectedIssueRetriever return the current selected issue
type SelectedIssueRetriever interface {
	GetSelectedIssue() *gitlab.Issue
}

//NewFooter creates a footer
func NewFooter(em *events.EventManager, issueRetriver SelectedIssueRetriever) Footer {
	footer := tview.NewTextView()
	footer.SetBorderPadding(1, 0, 0, 0)

	f := Footer{
		footer:        footer,
		em:            em,
		issueRetriver: issueRetriver,
	}

	f.registerListeners()

	return f
}

func (f *Footer) registerListeners() {
	f.em.Listen(events.ProjectSelected, f.showWebURL)
	f.em.Listen(events.IssueTableLineSelectionChanged, f.showWebURL)
	f.em.Listen(events.IssueTableRefreshed, f.showWebURL)

	f.em.Listen(events.IssueTableFocusLost, func() {
		f.footer.SetText("")
	})
}

func (f *Footer) showWebURL() {
	issue := f.issueRetriver.GetSelectedIssue()
	if issue != nil {
		f.footer.SetText(issue.WebURL)
	} else {
		f.footer.SetText("")
	}
}

//GetPrimitive returns the rivo/tview primtive
func (f *Footer) GetPrimitive() tview.Primitive {
	return f.footer
}
