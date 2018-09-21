# mayboy
CLI Tool to list gitlab issues for multiple projects.

![Screenshot](screen.png)

## Installation
```shell
go get -u github.com/windler/mayboy
```

Create a config-file in `~/.mayboy` with the following contents:

```yaml
gitlabHost: "https://gitlab-host.com"
accessToken: "secret-api-access-token-for-general-use"
maxIssues: 50 #defaults to 20
projects:
    Project1: <id>
    Project2: <id>
projectAccessTokens:
    Project1: "token-for-project-1"
```