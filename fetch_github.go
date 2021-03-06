package jobs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
)

const defaultArchieveFormat = "zipball"

// the parameters take the format OWNER, REPO, ARCHIEVE_FORMAT, REF
const githubURL = "https://api.github.com/repos/%s/%s/%s/%s"

// "GET /repos/:owner/:repo/:archive_format/:ref"

// GithubFetchStep will download the source code for the given **public** repo (should inject the correct client to fetch a private repo)
type GithubFetchStep struct {
	owner string
	repo  string
	ref   string
	log   *logrus.Logger
	pipeline.StepContext
}

// NewGithubStep takes what is essentially the URL of the repo and the ref to download
func NewGithubStep(owner, repo, ref string, logger *logrus.Logger) *GithubFetchStep {
	return &GithubFetchStep{
		owner: owner,
		repo:  repo,
		ref:   ref,
		log:   logger,
	}
}

// NewGithubStepFromEnvironment reads the owner, repo, and ref from the OWNER, REPO, and REF
// environment variables.
func NewGithubStepFromEnvironment() pipeline.Step {
	return NewGithubStep(os.Getenv("OWNER"), os.Getenv("REPO"), os.Getenv("REF"), nil)
}

// Exec runs the step. Should not be run directly.
func (g *GithubFetchStep) Exec(request *pipeline.Request) *pipeline.Result {
	g.Status(fmt.Sprintf("%+v", request))

	// Generate the URL to ping GitHub
	url := fmt.Sprintf(githubURL, g.owner, g.repo, defaultArchieveFormat, g.ref)
	fileUID := fmt.Sprintf("%v-%v-%v", g.owner, g.repo, g.ref)

	// Make a POST request to the server (TODO using the installation token in the future)
	g.Status("Fetching archive from GitHub...")
	resp, err := http.Get(url)
	if err != nil {
		g.Status("Failed to fetch archive from GitHub")
		return &pipeline.Result{Error: err}
	}

	// Create a temp file to store the file in
	tmpfile, err := ioutil.TempFile("", fileUID)
	defer tmpfile.Close()
	if err != nil {
		g.Status("Failed to create a temporary file.")
		return &pipeline.Result{Error: err}
	}

	/*
		stat, err := tmpfile.Stat()
		if err != nil {
			g.Status("Failed to read the file stats. Path error!")
			return &pipeline.Result{Error: err}
		}

		tmpfileName := stat.Name()
	*/
	tmpfileName := tmpfile.Name()

	// Save the zipball to the filesystem
	_, err = io.Copy(tmpfile, resp.Body)
	defer resp.Body.Close()
	if err != nil {
		g.Status("Failed to save the github archive to the temp file's body")
		return &pipeline.Result{Error: err}
	}

	// Make the temporary directory
	dir, err := ioutil.TempDir("", fileUID)
	if err != nil {
		g.Status("Failed to create a tmp dir")
		return &pipeline.Result{Error: err}
	}

	// Break open the archive
	err = archiver.Zip.Open(tmpfileName, dir)
	if err != nil {
		g.log.Warn(err)
		g.Status("Failed to unarchive the file")
		return &pipeline.Result{Error: err}
	}

	f, err := os.Open(dir)
	if err != nil {
		g.log.Warn(err)
		g.Status("Failed to open directory")
		return &pipeline.Result{Error: err}
	}
	dirs, err := f.Readdirnames(0)
	if err != nil {
		g.log.Warn(err)
		g.Status("Failed to read dir names")
		return &pipeline.Result{Error: err}
	}

	if len(dirs) != 1 {
		g.log.Warn("Expected just a single directory. Instead found %v", len(dirs))
		g.Status("Rando error")
		return nil
	}

	finalPath := filepath.Join(dir, dirs[0])

	// Finally, return the result
	return &pipeline.Result{
		Error:  nil,
		KeyVal: map[string]interface{}{"archive": finalPath},
	}
}

// Cancel is a no-op
func (g *GithubFetchStep) Cancel() error {
	g.Status("cancel step")
	return nil
}
