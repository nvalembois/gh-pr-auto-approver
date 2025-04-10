package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/sirupsen/logrus"

	"github.com/nvalembois/gh-pr-auto-approver/pkg/config"
)

func main() {
	// initialisation des logs
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.Infoln("gh-pr-approver: start")

	// traitement des arguments de la ligne de commande
	config := config.NewConfig()
	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Github Repo: ", config.GithubRepo)
		logrus.Debug("Github Base: ", config.GithubBase)
		logrus.Debug("Github Token Len:", len(config.GithubToken))
	}

	var httpClient *http.Client = nil
	if config.Debug {
		httpClient = &http.Client{
			Transport: &logTransport{http.DefaultTransport}}
	}
	client := github.NewClient(httpClient).WithAuthToken(config.GithubToken)

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
		Direction:   "asc",
		ListOptions: github.ListOptions{Page: 1, PerPage: 10}}
	for {
		if logrus.DebugLevel == logrus.GetLevel() {
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
		lastMergedTime := time.Unix(0, 0)
		for _, pr := range prlist {
			if time.Since(lastMergedTime) < (30 * time.Second) {
				time.Sleep((30 * time.Second) - time.Since(lastMergedTime))
			}
			if pr.GetRebaseable() {
				logrus.Debugf("PR #%d/%d '%s' is Rebaseable", *pr.Number, *pr.ID, *pr.Title)
				_, _, err := client.PullRequests.UpdateBranch(ctx, user, reponame, *pr.Number, nil)
				if err != nil {
					logrus.Errorf("Update PR #%d/%d '%s' failed: %s", *pr.Number, *pr.ID, *pr.Title, err)
					continue
				}
				tmppr, _, err := client.PullRequests.Get(ctx, user, reponame, *pr.Number)
				if err != nil {
					logrus.Errorf("Get PR #%d/%d '%s' failed: %s", *pr.Number, *pr.ID, *pr.Title, err)
					continue
				}
				if tmppr.GetRebaseable() {
					logrus.Infof("PR #%d/%d '%s' is still rebaseable", *pr.Number, *pr.ID, *pr.Title)
					continue
				}
			}
			if config.DryRun {
				logrus.Infof("dry-run: would merge PR #%d/%d '%s'", *pr.Number, *pr.ID, *pr.Title)
				continue
			}
			res, _, err := client.PullRequests.Merge(ctx,
				user, reponame, *pr.Number,
				"Fusion automatique de la pull request", nil)
			if err != nil {
				logrus.Errorf("Merge PR #%d/%d '%s' failed: %s", *pr.Number, *pr.ID, *pr.Title, err)
				continue
			}
			logrus.Infof("Merged PR #%d '%s': %s", *pr.ID, *pr.Title, *res.Message)
			lastMergedTime = time.Now()
		}
		logrus.Debugf("GetPR page: %d", opts.ListOptions.Page)
		opts.ListOptions.Page += 1
	}
	logrus.Infoln("gh-pr-approver: end")
}

// logTransport est une implémentation de http.RoundTripper qui journalise les requêtes sortantes.
type logTransport struct {
	Transport http.RoundTripper
}

// RoundTrip effectue la requête HTTP et journalise les détails.
func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	logrus.Debugf("Making API request to: %s %s\n", req.Method, req.URL)
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		logrus.Debugf("Error making request: %v\n", err)
		return nil, err
	}
	logrus.Debugf("Received API response with status: %s\n", resp.Status)
	return resp, nil
}
