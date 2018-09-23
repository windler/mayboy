package gitlab

import (
	"sort"
	"sync"
	"time"

	"github.com/windler/mayboy/config"
)

//ProjectIssues maps a Project to its issues
type ProjectIssues struct {
	project Project
	issues  []Issue
}

//GetIssues gets all issues for all prpjects defined in Config
func GetIssues(cfg config.Config) map[string][]Issue {
	result := map[string][]Issue{}

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

			client := NewClient(cfg.GitlabHost, token)
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

	if cfg.IncludeAll {
		allIssues := []Issue{}
		for _, issues := range result {
			for _, issue := range issues {
				allIssues = append(allIssues, issue)
			}
		}

		sort.Slice(allIssues, func(i, j int) bool {
			creationI, _ := time.Parse(time.RFC3339, allIssues[i].CreatedAt)
			creationJ, _ := time.Parse(time.RFC3339, allIssues[j].CreatedAt)

			return creationI.After(creationJ)
		})

		result["All"] = allIssues
	}

	return result
}
