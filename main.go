package main

import (
	"github.com/windler/mayboy/events"
	"github.com/windler/mayboy/gitlab"

	"github.com/rivo/tview"
	"github.com/windler/mayboy/config"
	"github.com/windler/mayboy/elements"
)

func main() {
	cfg := config.Parse()
	issues := gitlab.GetIssues(cfg)

	em := events.NewEventManager()
	app := tview.NewApplication()

	header := elements.NewHeader()

	projectList := elements.CreateProjectList(issues, em)
	issueTable := elements.CreateIssueTable(em, &projectList)
	tableInfo := elements.NewTableInfo(em, &issueTable)

	footer := elements.NewFooter(em)

	grid := elements.CreateGrid(
		header.GetPrimitive(),
		projectList.GetPrimitive(),
		issueTable.GetPrimitive(),
		tableInfo.GetPrimitive(),
		footer.GetPrimitive(),
	)

	em.Listen(events.ProjectSelected, func() {
		app.SetFocus(issueTable.GetPrimitive())
	})
	em.Listen(events.IssueTableFocusLost, func() {
		app.SetFocus(projectList.GetPrimitive())
	})
	em.Listen(events.ExitRequested, func() {
		app.Stop()
	})

	if err := app.SetRoot(grid.GetPrimitive(), true).Run(); err != nil {
		panic(err)
	}
}
