package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Project prepresents a gitlab project
type Project struct {
	Name string
	ID   int
}

//Client calls gitlab api
type Client struct {
	client *http.Client
	url    string
	token  string
}

//NewClient creates a new Client
func NewClient(url, token string) Client {
	return Client{
		client: &http.Client{},
		url:    url,
		token:  token,
	}
}

//Issue represents a gitlab issue
type Issue struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      struct {
		Name string `json:"name"`
	} `json:"author"`
	Assignee struct {
		Name string `json:"name"`
	} `json:"assignee"`
	WebURL string `json:"web_url"`
}

//GetIssues gets issues for a project
func (c *Client) GetIssues(projectID int, max int) []Issue {
	issues := &[]Issue{}

	perPage := max
	if perPage <= 0 {
		perPage = 20
	} else if perPage > 100 {
		perPage = 100
	}

	apiResult := c.requestAPI(fmt.Sprintf("/api/v4/projects/%d/issues?state=opened&per_page=%d", projectID, perPage))
	err := json.Unmarshal(apiResult, issues)

	if err != nil {
		log.Panicln("Error fetching issues for project " + string(projectID) + ": " + err.Error())
	}

	return *issues
}

func (c *Client) requestAPI(endpoint string) []byte {
	reqURL := c.url + endpoint
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Private-Token", c.token)
	resp, err := c.client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 401 {
		panic("Unauthorized gitlab access")
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes()
}
