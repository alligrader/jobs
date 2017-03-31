package jobs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
	log "github.com/sirupsen/logrus"
)

type (
	checkstyleStep struct {
		srcDir    string
		jarLoc    string
		outputLoc string
		checkLoc  string
		pipeline.StepContext
	}

	findbugsStep struct {
		srcDir    string
		jarLoc    string
		outputLoc string
		text      bool
		pipeline.StepContext
	}

	javacmd interface {
		init(*pipeline.Request) error
		setSrcDir(*pipeline.Request) error
		launchCmd() (string, error)
		Cmd() *exec.Cmd
		pipeline.Step
	}
)

const (
	DefaultCheckstyleJarLoc    = "/checkstyle-7.6.1.jar"
	DefaultCheckstyleOutputLoc = "/checkstyle_output.txt"
	DefaultCheckstyleConfigLoc = "/checks.xml"
	DefaultFindBugsJarLoc      = "/findbugs.jar"
	DefaultFindBugsOutputLoc   = "/findbugs_output.txt"
	DefaultSrcDir              = "/src"

	cmdTmplFindBugs     = "java -jar %s -textui -xml:withMessages -effort:max -output %s %s"
	cmdTmplFindBugsText = "java -jar %s -textui                   -effort:max -output %s %s"
	cmdTmplCheckstyle   = "java -jar %s -c %s -o %s %s"
)

// This line forces the compiler to check the method
// sets of the findbugsStep and checkstyleStep types
// to ensure that they both fulfill the javacmd interface
var _, _ javacmd = &findbugsStep{}, &checkstyleStep{}

func NewFindbugsStep(jarLoc, outputLoc, srcDir string, textoutput bool) pipeline.Step {
	return &findbugsStep{
		jarLoc:    jarLoc,
		outputLoc: outputLoc,
		srcDir:    srcDir,
		text:      textoutput,
	}
}

func NewCheckstyleStep(srcDir, jarLoc, outputLoc, checkLoc string) pipeline.Step {
	return &checkstyleStep{
		srcDir:    srcDir,
		jarLoc:    jarLoc,
		outputLoc: outputLoc,
		checkLoc:  checkLoc,
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

func (fb *findbugsStep) setSrcDir(request *pipeline.Request) error {

	if fb.srcDir != "" {
		return nil
	}

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

func (fb *findbugsStep) launchCmd() (string, error) {

	cmd := fb.Cmd()
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadFile(fb.outputLoc)
	return string(contents), err
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

func (fb *findbugsStep) Exec(request *pipeline.Request) *pipeline.Result {

	// Ensure all data is set
	if err := fb.init(request); err != nil {
		return &pipeline.Result{Error: err}
	}

	// Now, launch the command
	contents, err := fb.launchCmd()

	nextMap := fromMap(request.KeyVal)
	nextMap["findbugs"] = contents

	return &pipeline.Result{
		Error:  err,
		KeyVal: nextMap,
	}
}

func (checkstyle *checkstyleStep) Exec(request *pipeline.Request) *pipeline.Result {

	// Ensure all data is set
	if err := checkstyle.init(request); err != nil {
		return &pipeline.Result{Error: err}
	}

	// Now, launch the command
	contents, err := checkstyle.launchCmd()

	nextMap := fromMap(request.KeyVal)
	nextMap["checkstyle"] = contents

	return &pipeline.Result{
		Error:  err,
		KeyVal: nextMap,
	}
}

func (fb *findbugsStep) Cancel() error {
	fb.Status("Cancel")
	return nil
}

func (checkstyle *checkstyleStep) Cancel() error {
	checkstyle.Status("Cancel")
	return nil
}

func (fb *findbugsStep) Cmd() *exec.Cmd {
	var cmd string
	if fb.text {
		cmd = fmt.Sprintf(cmdTmplFindBugsText, fb.jarLoc, fb.outputLoc, fb.srcDir)
	} else {
		cmd = fmt.Sprintf(cmdTmplFindBugs, fb.jarLoc, fb.outputLoc, fb.srcDir)
	}
	log.Println(cmd)

	return exec.Command("bash", "-c", cmd)
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

func fromMap(m map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	if m == nil {
		return result
	}
	for key, val := range m {
		result[key] = val
	}
	return result
}
