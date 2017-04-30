package jobs

import (
	"context"
	"encoding/xml"
	"errors"
	"strconv"
	"strings"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

// POST /repos/:owner/:repo/commits/:sha/comments
const githubCommentURL = "https://api.github.com/repos/%s/%s/commits/%s/comments"

type commentStep struct {
	owner, repo, sha string
	checkstyleReport *Checkstyle
	findbugsReport   *BugCollection
	pipeline.StepContext
}

func NewCommentStep(owner, repo, sha string) pipeline.Step {
	return &commentStep{
		owner: owner,
		repo:  repo,
		sha:   sha,
	}
}

func (c *commentStep) loadCheckstyle(req *pipeline.Request) error {
	var (
		str     string
		err     error
		decoder *xml.Decoder
		check   Checkstyle
	)
	str, err = extractStr(req.KeyVal, "checkstyle")
	if err != nil {
		return err
	}

	decoder = xml.NewDecoder(strings.NewReader(str))
	err = decoder.Decode(&check)
	c.checkstyleReport = &check
	return err
}

func (c *commentStep) loadFindbugs(req *pipeline.Request) error {
	var (
		str      string
		err      error
		decoder  *xml.Decoder
		findbugs BugCollection
	)
	str, err = extractStr(req.KeyVal, "findbugs")
	if err != nil {
		return err
	}

	decoder = xml.NewDecoder(strings.NewReader(str))
	err = decoder.Decode(&findbugs)
	c.findbugsReport = &findbugs
	return err
}

func (c *commentStep) logReports() {
	c.logFindbugsReport()
	c.logCheckstyleReport()
}

func (c *commentStep) logFindbugsReport() {

	output, err := xml.MarshalIndent(c.findbugsReport, "  ", "    ")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Info(string(output))
}

func (c *commentStep) logCheckstyleReport() {

	output, err := xml.MarshalIndent(c.checkstyleReport, "  ", "    ")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Info(string(output))
}

func (c *commentStep) SendComment(ctx context.Context, client *github.Client, comment *github.RepositoryComment) error {

	repoService := client.Repositories
	_, _, err := repoService.CreateComment(ctx, c.owner, c.repo, c.sha, comment)

	return err
}

func (c *commentStep) Exec(req *pipeline.Request) *pipeline.Result {
	c.init(req)

	c.logReports()

	ctx := context.Background()
	client := github.NewClient(nil)

	for _, f := range c.checkstyleReport.File {
		for _, checkError := range f.Error {

			position, _ := strconv.Atoi(checkError.Line)
			comment := &github.RepositoryComment{
				Body:     &checkError.Message,
				Path:     &f.Name,
				Position: &position,
			}
			err := c.SendComment(ctx, client, comment)
			if err != nil {
				return &pipeline.Result{Error: err}
			}
		}
	}

	/*  THIS CODE IS... MAYBE WEIRD
	for _, bug := range c.findbugsReport {
		comment := &github.RepositoryComment{
			Body:     bug.Abbrev,
			Path:     bug.SourceLineBugInstance.Sourcefile,
			Position: strings.Atoi(bug.SourceLineBugInstance.Sourcepath),
		}
		resp, err := repoService.CreateComment(ctx, c.owner, c.repo, c.sha, comment)
		if err != nil {
			return &pipeline.Result{Error: err}
		}
		log.Warn(string(resp.Body))
	}
	*/
	// POST to GitHub the comments
	// https://godoc.org/github.com/google/go-github/github#RepositoriesService.CreateComment
	// https://gocodecloud.com/blog/2016/08/13/receiving-and-processing-github-api-events/
	return nil
}

func (c *commentStep) init(req *pipeline.Request) error {
	var err error

	if err = c.loadCheckstyle(req); err != nil {
		return err
	}
	if err = c.loadFindbugs(req); err != nil {
		return err
	}

	if c.owner == "" {
		c.owner, err = extractStr(req.KeyVal, "OWNER")
	}
	if err != nil {
		return err
	}

	if c.repo == "" {
		c.repo, err = extractStr(req.KeyVal, "REPO")
	}

	if err != nil {
		return err
	}

	if c.sha == "" {
		c.sha, err = extractStr(req.KeyVal, "SHA")
	}

	return err
}

func (c *commentStep) Cancel() error {
	c.Status("cancel step...")
	return nil
}

func extractStr(keyval map[string]interface{}, key string) (string, error) {
	if keyval == nil {
		return "", errors.New("KeyVal was nil.")
	}

	val, ok := keyval[key]
	if !ok {
		return "", errors.New("Not such key...")
	}

	str, ok := val.(string)
	if !ok {
		return "", errors.New("Value at key " + key + " is not a string")
	}

	return str, nil
}
