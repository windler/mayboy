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

	hint := elements.NewHint(em)

	projectList := elements.CreateProjectList(issues, em)
	issueTable := elements.CreateIssueTable(em, &projectList)

	footer := elements.NewFooter(em, &issueTable)

	grid := elements.CreateGrid(
		projectList.GetPrimitive(),
		issueTable.GetPrimitive(),
		hint.GetPrimitive(),
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
