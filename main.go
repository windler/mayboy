package main

import (
	"sync"

	"github.com/windler/mayboy/events"

	"github.com/rivo/tview"
	"github.com/windler/mayboy/config"
	"github.com/windler/mayboy/elements"
	"github.com/windler/mayboy/gitlab"
)

func main() {
	cfg := config.Parse()
	issues := GetIssues(cfg)

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

type ProjectIssues struct {
	project gitlab.Project
	issues  []gitlab.Issue
}

func GetIssues(cfg config.Config) map[string][]gitlab.Issue {
	result := map[string][]gitlab.Issue{}

	projects := []gitlab.Project{}
	for name, id := range cfg.Projects {
		projects = append(projects, gitlab.Project{
			Name: name,
			ID:   id,
		})
	}

	queue := make(chan ProjectIssues)
	var wg sync.WaitGroup

	wg.Add(len(projects))
	for _, project := range projects {
		go func(p gitlab.Project) {
			token := cfg.AccessToken

			if t, found := cfg.ProjectAccessTokens[p.Name]; found {
				token = t
			}

			client := gitlab.NewClient(cfg.GitlabHost, token)
			issues := client.GetIssues(p.ID, cfg.Max)

			projectIssues := ProjectIssues{
				issues:  issues,
				project: p,
			}
			queue <- projectIssues
		}(project)
	}

	go func() {
		for projectIssues := range queue {
			result[projectIssues.project.Name] = projectIssues.issues
			wg.Done()
		}
	}()

	wg.Wait()

	return result
}
