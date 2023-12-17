package main

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/sirupsen/logrus"

	"github.com/nvalembois/gh-pr-approver/internal"
)

func main() {
	// initialisation des logs
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.Infoln("gh-pr-approver: start")

	// traitement des arguments de la ligne de commande
	config := internal.NewConfig()
	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Github Repo: ", config.GithubRepo)
		logrus.Debug("Github Base: ", config.GithubBase)
		logrus.Debug("Github Token Len:", len(config.GithubToken))
	}

	client := github.NewClient(nil).WithAuthToken(config.GithubToken)

	user, reponame, _ := strings.Cut(config.GithubRepo, "/")
	ctx := context.Background()
	if config.GithubBase == "" {
		repo, _, err := client.Repositories.Get(ctx, user, reponame)
		if err != nil {
			logrus.Errorf("failed to get repo '%s': %s", config.GithubRepo, err)
			os.Exit(1)
		}
		logrus.Infof("Repo '%s' ID=%d", config.GithubRepo, repo.ID)
		config.GithubBase = *repo.DefaultBranch
		logrus.Infoln("Base:", config.GithubBase)
	}

	opts := github.PullRequestListOptions{
		State:       "open",
		Head:        "",
		Base:        config.GithubBase,
		Sort:        "created",
		Direction:   "desc",
		ListOptions: github.ListOptions{Page: 0, PerPage: 10}}
	for {
		if config.Debug {
			logrus.Debugf("Query PR for '%s' with options:", config.GithubRepo)
			logrus.Debugln("  State:", opts.State)
			logrus.Debugln("  Head:", opts.Head)
			logrus.Debugln("  Base:", opts.Base)
			logrus.Debugln("  Sort:", opts.Sort)
			logrus.Debugln("  Direction:", opts.Direction)
			logrus.Debugln("  Page:", opts.ListOptions.Page)
			logrus.Debugln("  PerPage:", opts.ListOptions.PerPage)
		}
		prlist, _, err := client.PullRequests.List(ctx, user, reponame, &opts)
		if err != nil {
			logrus.Errorf("failed to get opened pull request for '%s': %s", config.GithubRepo, err)
			os.Exit(1)
		}
		if len(prlist) == 0 {
			break
		}
		for _, pr := range prlist {
			if pr.Mergeable == nil || !*pr.Mergeable {
				logrus.Infof("PR #%d/'%s' is not mergable: %s", *pr.ID, *pr.Title, *pr.MergeableState)
				continue
			}
			res, _, err := client.PullRequests.Merge(ctx,
				user, reponame, int(*pr.ID),
				"Fusion automatique de la pull request", nil)
			if err != nil {
				logrus.Errorf("Merge PR #%d/'%s' failed: %s", *pr.ID, *pr.Title, err)
				continue
			}
			logrus.Infof("Merged PR #%d/'%s': %s", *pr.ID, *pr.Title, *res.Message)
		}
		logrus.Debugf("GetPR page: %d", opts.ListOptions.Page)
		opts.ListOptions.Page += 1
	}
}
