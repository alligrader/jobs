package jobs

import (
	"fmt"
	"reflect"

	"github.com/RobbieMcKinstry/pipeline"
)

func ExampleCommandStep() {

	const name = "test pipeline 1"

	var (
		cmd      = NewStepFromCommand("echo", "echo hello world")
		workpipe = pipeline.New(name, 10000)
		stage    = pipeline.NewStage(name, false, false)
	)

	stage.AddStep(cmd)
	workpipe.AddStage(stage)

	res := workpipe.Run()
	if res == nil {
		fmt.Println("Result was nil! Need to check the result for the cmd logs")
		return
	}

	if res.Error != nil {
		fmt.Println("experienced an error:")
		fmt.Println(res.Error)
	}

	out := res.KeyVal["stdout"]
	if stdout, ok := out.(string); ok {
		fmt.Println(stdout)
	} else {
		fmt.Println("Value at keyVal[stdout] is not a string")
		fmt.Println(reflect.TypeOf(out))
	}

	// Output:
	// hello world
}
