package jobs

import (
	"fmt"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
)

// CommandStep is a pipeline step for running a given string as a Bash command.
// Make a function that takes a command and returns a pipeline.Step from the string.
type CommandStep struct {
	name string
	cmd  *exec.Cmd
	pipeline.StepContext
}

// NewStepFromCommand creates a new CommandStep from the given string, using running `bash -C <string>`
func NewStepFromCommand(name, command string) *CommandStep {
	return &CommandStep{
		name: name,
		cmd:  exec.Command("bash", "-c", command),
	}
}

// Exec runs the command step, should be run by the pipeline, not directly.
func (c *CommandStep) Exec(request *pipeline.Request) *pipeline.Result {
	c.Status(fmt.Sprintf("%+v", request))
	out, err := c.cmd.Output()
	if err != nil {
		c.Cancel()
		return &pipeline.Result{
			Error: err,
		}
	}
	stdout := string(out)

	return &pipeline.Result{
		Error:  nil,
		Data:   struct{ msg string }{msg: stdout},
		KeyVal: map[string]interface{}{"stdout": stdout},
	}
}

// Cancel is a no-op, required to implement the pipeline step interface
func (c *CommandStep) Cancel() error {
	c.Status("cancel step")
	return nil
}
