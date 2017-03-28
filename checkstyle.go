package jobs

import (
	"errors"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
)

type checkstyleStep struct {
	srcDir     string
	configPath string
}

const (
	checkstyleJar = "/checkstyle-7.6.1.jar"
	outputPath    = "/checkstyle_output.txt"

	// TODO change this to use checkstyle instead
	// TODO do i add the text output not to a file but to stdout and store it in the ctx?
	cmdTmpl = "java -jar checkstyle-7.6.1.jar -c /sun_checks.xml MyClass.java"
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

	request.KeyVal["findbugs"] = stdout

	return &pipeline.Result{
		Error:  nil,
		KeyVal: request.KeyVal,
	}
}

func (fb *findbugsStep) Cmd() *exec.Cmd {
	cmd := tmplCmd(findbugsJar, outputPath, fb.srcDir)
	return exec.Command("bash", "-c", cmd)
}

func tmplCmd(jarPath, outputPath, srcPath string) {
	return fmt.Sprintf(cmdTmpl, jarPath, outputPath, srcPath)
}
