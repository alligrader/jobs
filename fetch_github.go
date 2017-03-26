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
)

const defaultArchieveFormat = "tarball"

// the parameters take the format OWNER, REPO, ARCHIEVE_FORMAT, REF
const githubURL = "https://api.github.com/repos/%s/%s/%s/%s"

// "GET /repos/:owner/:repo/:archive_format/:ref"

type githubFetchStep struct {
	owner string
	repo  string
	ref   string
	pipeline.StepContext
}

func NewGithubStep(owner, repo, ref string) pipeline.Step {
	return &githubFetchStep{
		owner: owner,
		repo:  repo,
	}
}

// NewGithubStepFromEnvironment reads the owner, repo, and ref from the OWNER, REPO, and REF
// environment variables.
func NewGithubStepFromEnvironment() pipeline.Step {
	return NewGithubStep(os.Getenv("OWNER"), os.Getenv("REPO"), os.Getenv("REF"))
}

func (g *githubFetchStep) Exec(request *pipeline.Request) *pipeline.Result {
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
	tmpfileName := filepath.Join(os.TempDir(), fileUID)
	defer tmpfile.Close()
	if err != nil {
		g.Status("Failed to create a temporary file.")
		return &pipeline.Result{Error: err}
	}

	// Save the tarball to the filesystem
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
	err = archiver.TarGz.Open(tmpfileName, dir)
	if err != nil {
		g.Status("Failed to unarchive the file")
		return &pipeline.Result{Error: err}
	}

	// Finally, return the result
	return &pipeline.Result{
		Error:  nil,
		KeyVal: map[string]interface{}{"archive": dir},
	}
}

func (g *githubFetchStep) Cancel() error {
	g.Status("cancel step")
	return nil
}
