package elements

import (
	"github.com/rivo/tview"
	"github.com/windler/mayboy/events"
	"github.com/windler/mayboy/gitlab"
)

type ProjectList struct {
	projectList *tview.List
	em          *events.EventManager
	projects    map[string][]gitlab.Issue
}

func CreateProjectList(projects map[string][]gitlab.Issue, em *events.EventManager) ProjectList {
	projectList := tview.NewList()
	projectList.ShowSecondaryText(false)

	if _, found := projects["All"]; found {
		projectList.AddItem("All", "", '0', func() {
			em.Fire(events.ProjectSelected)
		})
	}

	i := '1'
	for project := range projects {
		if project == "All" {
			continue
		}

		projectList.AddItem(project, "", i, func() {
			em.Fire(events.ProjectSelected)
		})
		i++
	}

	projectList.AddItem("Quit", "", 'q', func() {
		em.Fire(events.ExitRequested)
	})

	return ProjectList{
		projectList: projectList,
		em:          em,
		projects:    projects,
	}
}

func (l *ProjectList) GetIssuesForCurrentProject() []gitlab.Issue {
	currentItem := l.projectList.GetCurrentItem()
	proj, _ := l.projectList.GetItemText(currentItem)

	return l.projects[proj]
}

func (l *ProjectList) GetPrimitive() tview.Primitive {
	return l.projectList
}
