package jobs

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

func TestGitHubFetch(t *testing.T) {
	const (
		name             = "test pipeline 1"
		owner, repo, ref = "alligrader", "TestRepo", "d6a5d32f84e346574aded51404010d4ad2817641"
	)
	var (
		log   = logrus.New()
		fetch = NewGithubStep(owner, repo, ref, log)
		pipe  = pipeline.New(name, 1000)
		stage = pipeline.NewStage(name, false, false)
	)
	stage.AddStep(fetch)
	pipe.AddStage(stage)

	res := pipe.Run()
	if res == nil {
		t.Error("Result was nil!")
	}

	if res.Error != nil {
		t.Error(res.Error)
	}

	// now, check the filesystem for the output file.
	fileLoc := res.KeyVal["archive"]
	path, ok := fileLoc.(string)
	if !ok {
		t.Error(res.Error)
	}

	if _, err := os.Stat(path); err != nil {
		t.Error(err)
	}

	nestedFile := filepath.Join(path, "README.md")
	if _, err := os.Stat(nestedFile); err != nil {
		t.Error(err)
	}
}

func TestCommentStep(t *testing.T) {
	const owner, repo, sha string = "alligrader", "TestRepo", "d6a5d32f84e346574aded51404010d4ad2817641"

	httpclient := getClient()

	var log = logrus.New()
	var body, path = "Successful comment!", "README.md"
	var position = 1
	comment := &github.RepositoryComment{
		Body:     &body,
		Path:     &path,
		Position: &position,
	}
	client := github.NewClient(httpclient)
	commentStep := &CommentStep{
		owner:  owner,
		repo:   repo,
		sha:    sha,
		client: client,
		log:    log,
	}
	ctx := context.Background()

	err := commentStep.SendComment(ctx, client, comment)
	if err != nil {
		t.Error(err)
	}
}

func getClient() *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}
