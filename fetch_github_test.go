package jobs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RobbieMcKinstry/pipeline"
)

func TestGitHubFetch(t *testing.T) {
	const (
		name             = "test pipeline 1"
		owner, repo, ref = "alligrader", "TestRepo", "d6a5d32f84e346574aded51404010d4ad2817641"
	)
	var (
		fetch = NewGithubStep(owner, repo, ref)
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
