package main

import (
	"io/ioutil"
	"log"
	"os/user"
	"sync"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/windler/mayboy/gitlab"
	yaml "gopkg.in/yaml.v2"
)

var (
	app           *tview.Application
	projectList   *tview.List
	footer        *tview.TextView
	hint          *tview.TextView
	issueTable    *tview.Table
	grid          *tview.Grid
	projectIssues map[string][]gitlab.Issue
)

type Project struct {
	Name string
	ID   int
}

type Config struct {
	GitlabHost          string            `yaml:"gitlabHost"`
	AccessToken         string            `yaml:"accessToken"`
	Max                 int               `yaml:"maxIssues"`
	Projects            map[string]int    `yaml:"projects"`
	ProjectAccessTokens map[string]string `yaml:"projectAccessTokens"`
}

func main() {
	initIssues()

	app = tview.NewApplication()
	footer = tview.NewTextView()
	footer.SetBorderPadding(1, 0, 0, 0)
	hint = tview.NewTextView()
	hint.SetBorderPadding(1, 0, 0, 0)

	initProjectList()
	initIssueTable()
	initGrid()

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}

func initGrid() {
	grid = tview.NewGrid().
		SetRows(0, 2).
		SetColumns(30, 0)

	grid.AddItem(projectList, 0, 0, 1, 1, 0, 100, true)
	grid.AddItem(issueTable, 0, 1, 1, 1, 0, 100, false)
	grid.AddItem(hint, 1, 0, 1, 1, 0, 100, false)
	grid.AddItem(footer, 1, 1, 1, 1, 0, 100, false)
}

func initIssueTable() {
	issueTable = tview.NewTable()

	issueTable.Select(1, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == 'q' {
			app.Stop()
		}

		if key == tcell.KeyEsc {
			app.SetFocus(projectList)
			issueTable.SetSelectable(false, false)
			hint.SetText("")
			footer.SetText("")
		}
	}).SetSelectionChangedFunc(tableSelection)

	fillTable()
	issueTable.SetSelectable(false, false)
	hint.SetText("")
}

func initProjectList() {
	projectList = tview.NewList()
	projectList.ShowSecondaryText(false)

	i := '1'
	for project := range projectIssues {
		projectList.AddItem(project, "", i, fillTable)
		i++
	}

	projectList.AddItem("Quit", "", 'q', func() {
		app.Stop()
	})
}

func fillTable() {
	proj, _ := projectList.GetItemText(projectList.GetCurrentItem())
	issues := projectIssues[proj]

	issueTable.Clear()

	header := []string{"Title", "Author", "Assignee", "Creation"}
	for i := 0; i < len(header); i++ {
		issueTable.SetCell(0, i, tview.NewTableCell(header[i]).SetAttributes(tcell.AttrUnderline).SetSelectable(false))
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
			issueTable.SetCell(r+1, c, tview.NewTableCell(colData[c]).SetExpansion(1))
		}
	}

	issueTable.SetSelectable(true, false)
	app.SetFocus(issueTable)
	tableSelection(1, 0)
	hint.SetText("Press ESC to exit list")
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

func tableSelection(row int, column int) {
	if row > 0 {
		proj, _ := projectList.GetItemText(projectList.GetCurrentItem())
		issues := projectIssues[proj]

		if len(issues) > 0 {
			issue := issues[row-1]
			footer.SetText(issue.WebURL)
		}
	}
}

type ProjectIssues struct {
	project Project
	issues  []gitlab.Issue
}

func initIssues() {
	result := map[string][]gitlab.Issue{}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}

	dat, err := ioutil.ReadFile(usr.HomeDir + "/.mayboy")
	if err != nil {
		log.Fatal("Can not open " + usr.HomeDir + "/.mayboy")
	}

	err = yaml.Unmarshal(dat, cfg)
	if err != nil {
		log.Fatal(err)
	}

	projects := []Project{}
	for name, id := range cfg.Projects {
		projects = append(projects, Project{
			Name: name,
			ID:   id,
		})
	}

	queue := make(chan ProjectIssues)
	var wg sync.WaitGroup

	wg.Add(len(projects))
	for _, project := range projects {
		go func(p Project) {
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
	projectIssues = result
}
