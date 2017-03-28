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
	outputPath  = "/findbugs_output.txt"

	cmdTmpl = "jar -jar %s -textui -effort:max -output %s %s"
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
