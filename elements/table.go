package elements

import (
	"time"

	"github.com/windler/mayboy/gitlab"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/windler/mayboy/events"
)

//ProjectIssuesRetriever return the issues for current selected project
type ProjectIssuesRetriever interface {
	GetIssuesForCurrentProject() []gitlab.Issue
}

//IssueTable shows issues for a project in a table
type IssueTable struct {
	table         *tview.Table
	em            *events.EventManager
	issueRetriver ProjectIssuesRetriever
}

//CreateIssueTable creates a new IssueTable
func CreateIssueTable(em *events.EventManager, issueRetriver ProjectIssuesRetriever) IssueTable {
	issueTable := tview.NewTable()

	table := IssueTable{
		table:         issueTable,
		em:            em,
		issueRetriver: issueRetriver,
	}

	table.registerListeners()

	issueTable.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			issueTable.SetSelectable(false, false)
			em.Fire(events.IssueTableFocusLost)

		}
	}).SetSelectionChangedFunc(table.tableSelection)

	table.refreshTable()
	issueTable.SetSelectable(false, false)

	return table
}

func (t *IssueTable) registerListeners() {
	t.em.Listen(events.ProjectSelected, func() {
		t.table.SetSelectable(true, false)

		t.refreshTable()

		issues := t.issueRetriver.GetIssuesForCurrentProject()
		if len(issues) > 0 {
			t.table.Select(1, 0)
		}
	})
}

func (t *IssueTable) refreshTable() {
	issues := t.issueRetriver.GetIssuesForCurrentProject()

	t.table.Clear()

	header := []string{"Title", "Author", "Assignee", "Creation"}
	for i := 0; i < len(header); i++ {
		t.table.SetCell(0, i, tview.NewTableCell(header[i]).SetAttributes(tcell.AttrUnderline).SetSelectable(false))
	}

	for r := 0; r < len(issues); r++ {
		creation, _ := time.Parse(time.RFC3339, issues[r].CreatedAt)
		colData := []string{
			truncateString(issues[r].Title, 50),
			issues[r].Author.Name,
			issues[r].Assignee.Name,
			creation.Format("02.01.2006 - 15:04:05"),
		}

		for c := 0; c < len(colData); c++ {
			t.table.SetCell(r+1, c, tview.NewTableCell(colData[c]).SetExpansion(1))
		}
	}

	t.table.SetSelectable(true, false)
	t.tableSelection(1, 0)

	t.em.Fire(events.IssueTableRefreshed)
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func (t *IssueTable) tableSelection(row int, column int) {
	if row > 0 {
		t.em.Fire(events.IssueTableLineSelectionChanged)
	}
}

//GetSelectedIssue implements SelectedIssueRetriever
func (t *IssueTable) GetSelectedIssue() *gitlab.Issue {
	issues := t.issueRetriver.GetIssuesForCurrentProject()
	row, _ := t.table.GetSelection()

	if len(issues) == 0 || (row+1) > (len(issues)-1) {
		return nil
	}

	return &issues[row]
}

//GetPrimitive returns the rivo/tview primtive
func (t *IssueTable) GetPrimitive() tview.Primitive {
	return t.table
}
