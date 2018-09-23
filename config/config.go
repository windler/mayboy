package config

import (
	"io/ioutil"
	"log"
	"os/user"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	GitlabHost          string            `yaml:"gitlabHost"`
	AccessToken         string            `yaml:"accessToken"`
	Max                 int               `yaml:"maxIssues"`
	Projects            map[string]int    `yaml:"projects"`
	ProjectAccessTokens map[string]string `yaml:"projectAccessTokens"`
	IncludeAll          bool              `yaml:"includeAll"`
}

func Parse() Config {
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

	return *cfg
}
