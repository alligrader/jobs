package main

import (
	"github.com/alligrader/jobs"
)

func main() {

	var (
		t   *jobs.Tree = jobs.NewTree()
		cmd string     = "jar -jar /findbugs -textui -effort:max -output /findbugs_output.txt $SRC_DIR"
	)

	j := jobs.CmdVisitable(cmd)
	t.NewNode("Findbugs", "Root", j)

	t.Execute()
}
