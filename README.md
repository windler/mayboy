# mayboy
CLI Tool to list gitlab issues for multiple projects.

![Screenshot](screen.png)

## Installation
```shell
go get -u github.com/windler/mayboy
```

Create a config-file in `~/.mayboy` with the following contents:

```yaml
gitlabHost: "https://gitlab-host.com" #gitlab host name
accessToken: "secret-api-access-token-for-general-use" #api access token if neccessary
maxIssues: 50 #defaults to 20
includeAll: true #defaults to false
projects: #projects to show. Name->id
    Project1: <id>
    Project2: <id>
projectAccessTokens: #optional. Define specific acces tokens for some projects
    Project1: "token-for-project-1"
```