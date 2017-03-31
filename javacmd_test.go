package jobs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/RobbieMcKinstry/pipeline"
)

func ExampleFindBugs() {

	const (
		name   = "test pipeline 1"
		jarLoc = "lib/findbugs-3.0.1/lib/findbugs.jar"
		srcDir = ".test/src"
	)

	var (
		outputLoc, _ = ioutil.TempFile("", "findbugs.out")
		fbgs         = NewFindbugsStep(jarLoc, outputLoc.Name(), srcDir, true)
		workpipe     = pipeline.New(name, 10000)
		stage        = pipeline.NewStage(name, false, false)
	)

	defer os.Remove(outputLoc.Name())

	stage.AddStep(fbgs)
	workpipe.AddStage(stage)

	res := workpipe.Run()
	if res == nil {
		log.Fatal("Result was nil! Need to check the result for the cmd logs")
	}

	if res.Error != nil {
		log.Fatal(res.Error)
	}

	out, ok := res.KeyVal["findbugs"]
	if !ok {
		log.Fatal("No key 'findbugs' in the response KeyVal")
	}

	if checks, ok := out.(string); ok {
		fmt.Println(checks)
	} else {
		log.Println("Value at keyVal[findbugs] is not a string")
		log.Println(reflect.TypeOf(out))
	}

	// Output:
	// M C BIT: Bitwise OR of signed byte value computed in Main.main(String[])   At Main.java:[line 12]
}
