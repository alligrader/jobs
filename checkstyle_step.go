package jobs

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/sirupsen/logrus"
)

// NewCheckstyleStep creates a new pipeline step for running Checkstyle over a source directory.
// If any of the arguments are left as "", then they will use the package defaults instead.
func NewCheckstyleStep(jarLoc, srcDir, checkLoc, repoBase string, text bool, logger *logrus.Logger) *CheckstyleStep {
	return &CheckstyleStep{
		srcDir:   srcDir,
		jarLoc:   jarLoc,
		checkLoc: checkLoc,
		repoBase: repoBase,
		text:     text,
		log:      logger,
	}
}

func (checkstyle *CheckstyleStep) init(request *pipeline.Request) error {

	if err := checkstyle.setSrcDir(request); err != nil {
		return err
	}

	if checkstyle.jarLoc == "" {
		checkstyle.jarLoc = DefaultCheckstyleJarLoc
	}

	if checkstyle.srcDir == "" {
		checkstyle.srcDir = DefaultSrcDir
	}

	// Populate checkstyle.srcDir before populating repoBase instance variable
	if checkstyle.repoBase == "" {
		checkstyle.repoBase = checkstyle.srcDir
	}

	if checkstyle.checkLoc == "" {
		checkstyle.checkLoc = DefaultCheckstyleConfigLoc
	}

	return nil
}

func (checkstyle *CheckstyleStep) setSrcDir(request *pipeline.Request) error {
	if checkstyle.srcDir != "" {
		return nil
	}

	srcDirIntf, ok := request.KeyVal["archive"]
	if !ok {
		return errors.New("No source directory set.")
	}

	srcDir, ok := srcDirIntf.(string)
	if !ok {
		return errors.New("Source directory is not a string")
	}
	checkstyle.srcDir = srcDir
	checkstyle.log.Infof("Setting source directory to %v", srcDir)
	return nil
}

func (checkstyle *CheckstyleStep) launchCmd() (*Checkstyle, error) {

	log := checkstyle.log
	cmd := checkstyle.Cmd()
	stdout, stderr, err := checkstyle.captureOutput(cmd)
	if err != nil {
		log.Warn("Could not collect stdout/stderr")
		return nil, err
	}
	tee, stream2 := checkstyle.teeStream(stdout)

	if err = checkstyle.startCmd(cmd); err != nil {
		return nil, err
	}

	var check = &Checkstyle{}
	if err = xml.NewDecoder(tee).Decode(&check); err != nil {
		log.Warn("Decoding failed!")
		log.Warn("Dumping Stdout")
		dumpStream(stream2, log)
		log.Warn("Dumping Stderr")
		dumpStream(stderr, log)
		return nil, err
	}

	if err = cmd.Wait(); err != nil {
		checkstyle.log.Warn("Failed to wait for command completion")
		return nil, err
		// Do not call Fatal because we need to dump stderr first.
	}

	checkstyle.log.Info("Program has finished running")
	if err != nil {
		log.Info("Dumping STDOUT")
		dumpStream(stream2, log)
		log.Info("Dumping STDERR")
		dumpStream(stderr, log)
		return nil, err
	}

	checkstyle.listFiles()
	checkstyle.log.Info("Completed the marshalling of the Checkstyle struct.")
	return check, err
}

func dumpStream(stream io.Reader, log *logrus.Logger) {
	b, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Warn(err)
	}

	log.Infof("Stdout is: %s", b)
}

func (checkstyle *CheckstyleStep) teeStream(stream io.Reader) (io.Reader, io.Reader) {
	var buffer bytes.Buffer
	var tee = io.TeeReader(stream, &buffer)
	return tee, stream
}

func (checkstyle *CheckstyleStep) captureOutput(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		return nil, nil, err1
	}
	if err2 != nil {
		return nil, nil, err2
	}
	return stdout, stderr, nil
}

// TODO delete, only used for debugging purposes.
func (checkstyle *CheckstyleStep) listFiles() {
	files, err := ioutil.ReadDir(checkstyle.srcDir)
	if err != nil {
		checkstyle.log.Fatalf("Error while listing directory: %v", err)
	}

	str := ""
	for _, f := range files {
		str += f.Name() + "\n"
	}
	checkstyle.log.Infof("Files:\n%v", str)
}

func (checkstyle *CheckstyleStep) startCmd(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		checkstyle.log.Warn("Failed to start the command")
		return err
	}
	return nil
}

// Exec will run the step with the provided *pipeline.Request. Should be run by the pipeline, not directly.
func (checkstyle *CheckstyleStep) Exec(request *pipeline.Request) *pipeline.Result {

	// Ensure all data is set
	if err := checkstyle.init(request); err != nil {
		return &pipeline.Result{Error: err}
	}

	// Now, launch the command
	check, err := checkstyle.launchCmd()
	if err != nil {
		return &pipeline.Result{Error: err}
	}

	check = checkstyle.filterPath(check)

	nextMap := fromMap(request.KeyVal)
	nextMap["checkstyle"] = check

	return &pipeline.Result{
		Error:  err,
		KeyVal: nextMap,
	}
}

// filterPath walks through each file and removes the base of the path
// from the File.Name property. This is a transformation which makes it possible to
// post-back comments to GitHub. GitHub only know the path from the base
// of the directory, not the absolute path on the machine's filesystem.
func (checkstyle *CheckstyleStep) filterPath(ch *Checkstyle) *Checkstyle {
	checkstyle.log.Info("Filtering the file paths...")
	checkstyle.log.Infof("There are %v files with errors\n", len(ch.File))
	for index, f := range ch.File {
		base, err := filepath.Abs(checkstyle.repoBase)
		if err != nil {
			checkstyle.log.Warn("Error in determining the absolute path. Exiting")
			checkstyle.log.Fatal(err)
		}
		regexDescriptor := fmt.Sprintf("^%s", base)
		r := regexp.MustCompile(regexDescriptor)
		fileName := f.Name
		if loc := r.FindStringIndex(fileName); loc != nil {
			checkstyle.log.Info("Found a match in the filename.")
			start := loc[1] + 1
			ch.File[index].Name = fileName[start:]
			checkstyle.log.Warnf("Locations are: (%v, %v)", loc[0], loc[1])
			checkstyle.log.Infof("File name is now %v\n and file.Name is now %v", ch.File[index].Name, f.Name)
		} else {
			checkstyle.log.Warnf("Found a match in the filename.\nRegex Descriptor: '%s', filename: %s", regexDescriptor, fileName)
		}
	}
	checkstyle.log.Info("Finished filtering paths.")
	return ch
}

func (checkstyle *CheckstyleStep) serialize(blob string) (*Checkstyle, error) {
	var (
		check   Checkstyle
		decoder = xml.NewDecoder(strings.NewReader(blob))
		err     = decoder.Decode(&check)
	)
	return &check, err
}

// Cancel will cancel the step. Does nothing.
// Could be rewritten to kill the subprocess
func (checkstyle *CheckstyleStep) Cancel() error {
	checkstyle.Status("Cancel")
	return nil
}

// Cmd returns a *exec.Cmd configued to run Checkstyle over the source code referenced in the CheckstyleStep struct.
func (checkstyle *CheckstyleStep) Cmd() *exec.Cmd {
	var strTmpl = cmdTmplCheckstyle

	if checkstyle.text {
		strTmpl = cmdTmplCheckstyleText
	}
	cmd := fmt.Sprintf(
		strTmpl,
		checkstyle.jarLoc,
		checkstyle.checkLoc,
		checkstyle.srcDir,
	)

	return exec.Command("bash", "-c", cmd)
}
