package jobs

import (
	"errors"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
)

type findbugsStep struct {
	srcDir string
}

const (
	findbugsJar = "/findbugs.jar"
	findbugsOut = "/findbugs_output.txt"

	cmdTmplFindBugs = "jar -jar %s -textui -effort:max -output %s %s"
)

func (fb *findbugsStep) Exec(request *pipeline.Request) *pipeline.Result {

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
	contents, err := ioutil.ReadFile(findbugsOut)
	if err != nil {
		return &pipeline.Result{Error: err}
	}

	request.KeyVal["findbugs"] = string(contents)
	request.KeyVal["findbugsStdout"] = stdout

	return &pipeline.Result{
		Error:  nil,
		KeyVal: request.KeyVal,
	}
}

func (fb *findbugsStep) Cmd() *exec.Cmd {
	cmd := tmplCmdFindbugs(findbugsJar, findbugsOut, fb.srcDir)
	return exec.Command("bash", "-c", cmd)
}

func tmplCmdFindbugs(jarPath, outputPath, srcPath string) {
	return fmt.Sprintf(cmdTmplFindBugs, jarPath, outputPath, srcPath)
}
