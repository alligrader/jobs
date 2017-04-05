package jobs

import (
	"os"

	"github.com/RobbieMcKinstry/pipeline"
)

type environmentStep struct {
	strings []string
	pipeline.StepContext
}

func NewEnvStep(vars []string) pipeline.Step {
	return &environmentStep{
		strings: vars,
	}
}

func (e *environmentStep) Cancel() error {
	e.Status("cancel step")
	return nil
}

func (e *environmentStep) Exec(request *pipeline.Request) *pipeline.Result {
	keyVal := fromMap(request.KeyVal)
	for _, str := range e.strings {
		keyVal[str] = os.Getenv(str)
	}
	return &pipeline.Result{
		Error:  nil,
		KeyVal: keyVal,
	}
}
