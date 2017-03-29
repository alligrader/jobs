package jobs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
)

const (
	DefaultFindBugsJarLoc    = "/findbugs.jar"
	DefaultFindBugsOutputLoc = "/findbugs_output.txt"
	DefaultSrcDir            = "/src"

	cmdTmplFindBugs = "jar -jar %s -textui -effort:max -output %s %s"
)

type findbugsStep struct {
	srcDir    string
	jarLoc    string
	outputLoc string
	pipeline.StepContext
}

func NewFindbugsStep(jarLoc, outputLoc, srcDir string) pipeline.Step {
	return &findbugsStep{
		jarLoc:    jarLoc,
		outputLoc: outputLoc,
		srcDir:    srcDir,
	}
}

func (fb *findbugsStep) init(request *pipeline.Request) error {

	if err := fb.setSrcDir(request); err != nil {
		return err
	}

	if fb.jarLoc == "" {
		fb.jarLoc = DefaultFindBugsJarLoc
	}
	if fb.srcDir == "" {
		fb.srcDir = DefaultSrcDir
	}
	if fb.outputLoc == "" {
		fb.outputLoc = DefaultFindBugsOutputLoc
	}

	return nil
}

func (fb *findbugsStep) setSrcDir(request *pipeline.Request) error {

	srcDirIntf, ok := request.KeyVal["src"]
	if !ok {
		return errors.New("No source directory set.")
	}

	srcDir, ok := srcDirIntf.(string)
	if !ok {
		return errors.New("Source directory is not a string")
	}
	fb.srcDir = srcDir
	return nil
}

func (fb *findbugsStep) launchCmd() (string, error) {

	cmd := fb.Cmd()
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadFile(fb.outputLoc)
	return string(contents), err
}

func (fb *findbugsStep) Exec(request *pipeline.Request) *pipeline.Result {

	// Ensure all data is set
	if err := fb.init(request); err != nil {
		return &pipeline.Result{Error: err}
	}

	// Now, launch the command
	contents, err := fb.launchCmd()
	request.KeyVal["findbugs"] = contents

	return &pipeline.Result{
		Error:  err,
		KeyVal: request.KeyVal,
	}
}

func (fb *findbugsStep) Cancel() error {
	fb.Status("Cancel")
	return nil
}

func (fb *findbugsStep) Cmd() *exec.Cmd {
	cmd := fmt.Sprintf(cmdTmplFindBugs, fb.jarLoc, fb.outputLoc, fb.srcDir)
	return exec.Command("bash", "-c", cmd)
}
