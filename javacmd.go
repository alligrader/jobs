package jobs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/sirupsen/logrus"
)

type (
	CheckstyleStep struct {
		srcDir   string
		jarLoc   string
		checkLoc string
		repoBase string
		text     bool
		log      *logrus.Logger
		pipeline.StepContext
	}

	findbugsStep struct {
		srcDir    string
		jarLoc    string
		outputLoc string
		text      bool
		log       *logrus.Logger
		pipeline.StepContext
	}

	javacmd interface {
		init(*pipeline.Request) error
		setSrcDir(*pipeline.Request) error
		Cmd() *exec.Cmd
		pipeline.Step
	}
)

const (
	DefaultCheckstyleJarLoc    = "/checkstyle-7.6.1.jar"
	DefaultCheckstyleConfigLoc = "/checks.xml"
	DefaultFindBugsJarLoc      = "/findbugs.jar"
	DefaultFindBugsOutputLoc   = "/findbugs_output.txt"
	DefaultSrcDir              = "/src"

	cmdTmplFindBugs       = "java -jar %s -textui -xml:withMessages -effort:max -output %s %s"
	cmdTmplFindBugsText   = "java -jar %s -textui                   -effort:max -output %s %s"
	cmdTmplCheckstyle     = "java -jar %s -c %s -f xml %s"
	cmdTmplCheckstyleText = "java -jar %s -c %s %s"
)

// This line forces the compiler to check the method
// sets of the findbugsStep and checkstyleStep types
// to ensure that they both fulfill the javacmd interface
var _, _ javacmd = &findbugsStep{}, &CheckstyleStep{}

func NewFindbugsStep(jarLoc, outputLoc, srcDir string, textoutput bool, logger *logrus.Logger) pipeline.Step {
	return &findbugsStep{
		jarLoc:    jarLoc,
		outputLoc: outputLoc,
		srcDir:    srcDir,
		text:      textoutput,
		log:       logger,
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

	if fb.srcDir != "" {
		return nil
	}

	srcDirIntf, ok := request.KeyVal["archive"]
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
	nextMap := fromMap(request.KeyVal)
	nextMap["findbugs"] = contents

	return &pipeline.Result{
		Error:  err,
		KeyVal: nextMap,
	}
}

func (fb *findbugsStep) Cancel() error {
	fb.Status("Cancel")
	return nil
}

func (fb *findbugsStep) Cmd() *exec.Cmd {
	var strTmpl string = cmdTmplFindBugs

	if fb.text {
		strTmpl = cmdTmplFindBugsText
	}

	cmd := fmt.Sprintf(
		strTmpl,
		fb.jarLoc,
		fb.outputLoc,
		fb.srcDir,
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
