package jobs

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/sirupsen/logrus"
)

func exampleFindBugs() {

	const (
		name   = "test pipeline 1"
		jarLoc = "lib/findbugs-3.0.1/lib/findbugs.jar"
		srcDir = ".test/src"
	)

	var (
		log          = logrus.New()
		outputLoc, _ = ioutil.TempFile("", "findbugs.out")
		fbgs         = NewFindbugsStep(jarLoc, outputLoc.Name(), srcDir, true, log)
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

// TODO find out why this doesn't include the text "[Needs Braces]" in the output. That seemed to work before we marshalled the text?
func exampleCheckstyle() {
	const (
		name   = "test pipeline 1"
		jarLoc = "lib/checkstyle-7.6.1-all.jar"
		srcDir = ".test/src"
		checks = ".test/checkstyle.xml"
	)

	var (
		log        = logrus.New()
		checkstyle = NewCheckstyleStep(jarLoc, srcDir, checks, srcDir, false, log)
		workpipe   = pipeline.New(name, 10000)
		stage      = pipeline.NewStage(name, false, false)
	)

	stage.AddStep(checkstyle)
	workpipe.AddStage(stage)

	res := workpipe.Run()
	if res == nil {
		log.Fatal("Result was nil! Need to check the result for the cmd logs")
	}

	if res.Error != nil {
		log.Fatal(res.Error)
	}

	out, ok := res.KeyVal["checkstyle"]
	if !ok {
		log.Fatal("No key 'checkstyle' in the response KeyVal")
	}

	if checks, ok := out.(*Checkstyle); ok {
		fmt.Println(checks.File[0].Error[0].Message)
	} else {
		log.Println("Value at keyVal[checkstyle] is not a checkstyle struct")
		log.Println(reflect.TypeOf(out))
	}

	// Output:
	// 'for' construct must use '{}'s.
}

func TestCanReadSchema(t *testing.T) {
	contents, err := ioutil.ReadFile(".test/findbugs.out")
	if err != nil {
		t.Fatal(err)
	}
	var bugs bugcollection
	err = xml.Unmarshal(contents, &bugs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCanReadSchema2(t *testing.T) {
	contents, err := ioutil.ReadFile(".test/checkstyle.out")
	if err != nil {
		t.Fatal(err)
	}
	var style Checkstyle
	err = xml.Unmarshal(contents, &style)
	if err != nil {
		t.Fatal(err)
	}
}
