package jobs

import (
	"errors"
	"io/ioutil"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
)

type checkstyleStep struct {
	srcDir     string
	configPath string
}

const (
	checkstyleJar = "/checkstyle-7.6.1.jar"
	checkstyleOut = "/checkstyle_output.txt"
	checkPath     = "/checks.xml"

	// TODO change this to use checkstyle instead
	// TODO do i add the text output not to a file but to stdout and store it in the ctx?
	cmdTmplCheckstyle = "java -jar %s -c %s -o %s %s"
)

func (checkstyle *checkstyleStep) Exec(request *pipeline.Request) *pipeline.Result {

	srcDirIntf, ok := request.KeyVal["src"]
	if !ok {
		return &pipeline.Result{Error: errors.New("No source directory set.")}
	}

	srcDir, ok := srcDirIntf.(string)
	if !ok {
		return &pipeline.Result{Error: errors.New("Source directory is not a string")}
	}

	cmd := fb.Cmd()
	out, err := cmd.Output()
	if err != nil {
		return &pipeline.Result{Error: err}
	}
	stdout := string(out)
	request.KeyVal["checkpathStdout"] = stdout

	contents, err := ioutil.ReadFile(checkstyleOut)
	if err != nil {
		return &pipeline.Result{Error: err}
	}
	request.KeyVal["checkstyle"] = string(contents)

	return &pipeline.Result{
		Error:  nil,
		KeyVal: request.KeyVal,
	}
}

func (fb *findbugsStep) Cmd() *exec.Cmd {
	cmd := tmplCmdCheckstyle(checkstyleJar, checkPath, checkstyleOut, fb.srcDir)
	return exec.Command("bash", "-c", cmd)
}

func tmplCmdCheckstyle(jarPath, checkPath, outputPath, srcPath string) {
	return fmt.Sprintf(cmdTmplCheckstyle, jarPath, checkPath, outputPath, srcPath)
}
