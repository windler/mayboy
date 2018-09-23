package elements

import (
	"github.com/rivo/tview"
	"github.com/windler/mayboy/events"
	"github.com/windler/mayboy/gitlab"
)

//TableInfo show contextual information
type TableInfo struct {
	footer        *tview.TextView
	em            *events.EventManager
	issueRetriver SelectedIssueRetriever
}

//SelectedIssueRetriever return the current selected issue
type SelectedIssueRetriever interface {
	GetSelectedIssue() *gitlab.Issue
}

//NewTableInfo creates a TableInfo
func NewTableInfo(em *events.EventManager, issueRetriver SelectedIssueRetriever) TableInfo {
	footer := tview.NewTextView()
	footer.SetBorderPadding(1, 0, 0, 0)

	f := TableInfo{
		footer:        footer,
		em:            em,
		issueRetriver: issueRetriver,
	}

	f.registerListeners()

	return f
}

func (f *TableInfo) registerListeners() {
	f.em.Listen(events.ProjectSelected, f.showWebURL)
	f.em.Listen(events.IssueTableLineSelectionChanged, f.showWebURL)
	f.em.Listen(events.IssueTableRefreshed, f.showWebURL)

	f.em.Listen(events.IssueTableFocusLost, func() {
		f.footer.SetText("")
	})
}

func (f *TableInfo) showWebURL() {
	issue := f.issueRetriver.GetSelectedIssue()
	if issue != nil {
		f.footer.SetText(issue.WebURL)
	} else {
		f.footer.SetText("")
	}
}

//GetPrimitive returns the rivo/tview primtive
func (f *TableInfo) GetPrimitive() tview.Primitive {
	return f.footer
}
