package main

import (
	"github.com/alligrader/jobs"
)

func main() {

	var (
		findBugsCmd   string = "jar -jar /findbugs -textui -effort:max -output /findbugs_output.txt $SRC_DIR"
		checkstyleCmd string = ""
		pipe                 = pipeline.NewProgress("staticAnalysis", 1000, time.Minutes*10)
		preaction            = pipeline.NewStage("preaction", false, false)
		action               = pipeline.NewStage("action", false, false)
		postaction           = pipeline.NewStage("postaction", false, false)

		// Steps:
		// Inject the environment variables
		// Fetch the data from GitHub using the environment variables
		// Run the two steps below.
		// Persist the information back to GH
		findbugs   = jobs.NewStepFromCommand("findbugs", findBugsCmd)
		checkstyle = jobs.NewStepFromCommand("checkstyle", checkstyleCmd)
	)

	action.AddStep(findbugs)
	action.AddStep(checkstyle)

	pipe.AddStage(preaction)
	pipe.AddStage(action)
	pipe.AddStage(postaction)

	if err := pipe.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("timeTaken:", workpipe.GetDuration())
}
