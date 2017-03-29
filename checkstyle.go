package jobs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
)

type checkstyleStep struct {
	srcDir    string
	jarLoc    string
	outputLoc string
	checkLoc  string
	pipeline.StepContext
}

const (
	DefaultCheckstyleJarLoc    = "/checkstyle-7.6.1.jar"
	DefaultCheckstyleOutputLoc = "/checkstyle_output.txt"
	DefaultCheckstyleConfigLoc = "/checks.xml"

	cmdTmplCheckstyle = "java -jar %s -c %s -o %s %s"
)

func NewCheckstyleStep(srcDir, jarLoc, outputLoc, checkLoc string) pipeline.Step {
	return &checkstyleStep{
		srcDir:    srcDir,
		jarLoc:    jarLoc,
		outputLoc: outputLoc,
		checkLoc:  checkLoc,
	}
}

func (checkstyle *checkstyleStep) Cancel() error {
	checkstyle.Status("Cancel")
	return nil
}

func (checkstyle *checkstyleStep) init(request *pipeline.Request) error {

	if err := checkstyle.setSrcDir(request); err != nil {
		return err
	}

	if checkstyle.jarLoc == "" {
		checkstyle.jarLoc = DefaultCheckstyleJarLoc
	}

	if checkstyle.srcDir == "" {
		checkstyle.srcDir = DefaultSrcDir
	}

	if checkstyle.outputLoc == "" {
		checkstyle.outputLoc = DefaultCheckstyleOutputLoc
	}

	if checkstyle.checkLoc == "" {
		checkstyle.checkLoc = DefaultCheckstyleConfigLoc
	}

	return nil
}

func (checkstyle *checkstyleStep) setSrcDir(request *pipeline.Request) error {

	srcDirIntf, ok := request.KeyVal["src"]
	if !ok {
		return errors.New("No source directory set.")
	}

	srcDir, ok := srcDirIntf.(string)
	if !ok {
		return errors.New("Source directory is not a string")
	}
	checkstyle.srcDir = srcDir
	return nil
}

func (checkstyle *checkstyleStep) launchCmd() (string, error) {

	cmd := checkstyle.Cmd()
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadFile(checkstyle.outputLoc)
	return string(contents), err
}

func (checkstyle *checkstyleStep) Exec(request *pipeline.Request) *pipeline.Result {

	// Ensure all data is set
	if err := checkstyle.init(request); err != nil {
		return &pipeline.Result{Error: err}
	}

	// Now, launch the command
	contents, err := checkstyle.launchCmd()
	request.KeyVal["checkstyle"] = contents

	return &pipeline.Result{
		Error:  err,
		KeyVal: request.KeyVal,
	}
}

func (checkstyle *checkstyleStep) Cmd() *exec.Cmd {
	cmd := fmt.Sprintf(
		cmdTmplCheckstyle,
		checkstyle.jarLoc,
		checkstyle.checkLoc,
		checkstyle.outputLoc,
		checkstyle.srcDir,
	)

	return exec.Command("bash", "-c", cmd)
}
