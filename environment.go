package jobs

import (
	"os"

	"github.com/RobbieMcKinstry/pipeline"
)

// EnvironmentStep reads a bunch of strings out of the environemt (variables) and stores them in the response context.
type EnvironmentStep struct {
	strings []string
	pipeline.StepContext
}

// NewEnvStep creates an environment step which looks in the environment for variables with the given names.
func NewEnvStep(vars []string) *EnvironmentStep {
	return &EnvironmentStep{
		strings: vars,
	}
}

// Cancel is a no-op
func (e *EnvironmentStep) Cancel() error {
	e.Status("cancel step")
	return nil
}

// Exec runs the step. Should be run as part of the pipeline, not directly.
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
