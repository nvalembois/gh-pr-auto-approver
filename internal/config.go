package internal

import (
	"flag"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	GithubToken string
	GithubRepo  string
	GithubBase  string
	Debug       bool
}

func NewConfig() *Config {
	var c Config
	flag.StringVar(&(c.GithubToken), "token", os.Getenv("GITHUB_TOKEN"), "gihub token")
	flag.StringVar(&(c.GithubRepo), "repo", os.Getenv("GITHUB_REPO"), "github repo")
	flag.StringVar(&(c.GithubBase), "base", os.Getenv("GITHUB_BASE"), "github base")
	flag.BoolVar(&(c.Debug), "debug", boolVarOrDefault("LOGLEVEL_DEBUG", false), "debug")
	flag.Parse()
	assertString("^[\\w-]+$", c.GithubToken, "invalid or missing github token")
	assertString("^[\\w-]+(/[\\w-]+)+$", c.GithubRepo, "invalid or missing github repo")
	assertString("^[\\w-]*$", c.GithubBase, "invalid github base")
	return &c
}

func assertString(regex string, value string, message string) {
	res, err := regexp.MatchString(regex, value)
	if err != nil {
		logrus.Errorln(message, err)
		os.Exit(1)
	} else if !res {
		logrus.Errorln(message)
		os.Exit(1)
	}
}

func boolVarOrDefault(envVar string, defaultValue bool) bool {
	result, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	return strings.EqualFold(result, "true")
}
