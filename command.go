package jobs

import (
	"fmt"
	"os/exec"

	"github.com/RobbieMcKinstry/pipeline"
)

// Make a function that takes a command and returns a pipeline.Step from the string.
type commandStep struct {
	name string
	cmd  *exec.Cmd
	pipeline.StepContext
}

func NewStepFromCommand(name, command string) pipeline.Step {
	return &commandStep{
		name: name,
		cmd:  exec.Command("bash", "-c", command),
	}
}

func (c *commandStep) Exec(request *pipeline.Request) *pipeline.Result {
	c.Status(fmt.Sprintf("%+v", request))
	out, err := c.cmd.Output()
	if err != nil {
		c.Cancel()
		fmt.Println("FAILED")
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

func (c *commandStep) Cancel() error {
	c.Status("cancel step")
	fmt.Println("Cancellation request made!")
	return nil
}
