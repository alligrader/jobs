package jobs

import (
	"log"
	"os"
	"testing"

	"github.com/RobbieMcKinstry/pipeline"
)

func TestEnvironmentStep(t *testing.T) {
	const (
		name  = "test pipeline 1"
		value = "xyz123"
	)

	var (
		vars     = []string{"abcdefg"}
		env      = NewEnvStep(vars)
		workpipe = pipeline.New(name, 10000)
		stage    = pipeline.NewStage(name, false, false)
	)

	os.Setenv(vars[0], value)
	defer os.Unsetenv(vars[0])

	stage.AddStep(env)
	workpipe.AddStage(stage)

	res := workpipe.Run()
	if res == nil {
		log.Fatal("Result was nil! Need to check the result for the cmd logs")
	}

	if res.Error != nil {
		log.Fatal(res.Error)
	}

	out, ok := res.KeyVal[vars[0]]
	if !ok {
		log.Fatalf("No key '%v' in the response KeyVal", vars[0])
	}

	if _, ok := out.(string); !ok {
		log.Fatal("result is not a string")
	}

	if result := out.(string); result != value {
		log.Fatal("expected %v, observed %v", value, result)
	}

}
