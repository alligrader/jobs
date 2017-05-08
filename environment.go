package jobs

import (
	"os"

	"github.com/RobbieMcKinstry/pipeline"
)

type EnvironmentStep struct {
	strings []string
	pipeline.StepContext
}

func NewEnvStep(vars []string) *EnvironmentStep {
	return &EnvironmentStep{
		strings: vars,
	}
}

func (e *EnvironmentStep) Cancel() error {
	e.Status("cancel step")
	return nil
}

func (e *EnvironmentStep) Exec(request *pipeline.Request) *pipeline.Result {
	keyVal := fromMap(request.KeyVal)
	for _, str := range e.strings {
		keyVal[str] = os.Getenv(str)
	}
	return &pipeline.Result{
		Error:  nil,
		KeyVal: keyVal,
	}
}
