package jobs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RobbieMcKinstry/pipeline"
	"github.com/sirupsen/logrus"
)

type (
	checkstyleStep struct {
		srcDir    string
		jarLoc    string
		outputLoc string
		checkLoc  string
		repoBase  string
		text      bool
		log       *logrus.Logger
		pipeline.StepContext
	}

	findbugsStep struct {
		srcDir    string
		jarLoc    string
		outputLoc string
		text      bool
		log       *logrus.Logger
		pipeline.StepContext
	}

	javacmd interface {
		init(*pipeline.Request) error
		setSrcDir(*pipeline.Request) error
		launchCmd() (string, error)
		Cmd() *exec.Cmd
		pipeline.Step
	}
)

const (
	DefaultCheckstyleJarLoc    = "/checkstyle-7.6.1.jar"
	DefaultCheckstyleOutputLoc = "/checkstyle_output.txt"
	DefaultCheckstyleConfigLoc = "/checks.xml"
	DefaultFindBugsJarLoc      = "/findbugs.jar"
	DefaultFindBugsOutputLoc   = "/findbugs_output.txt"
	DefaultSrcDir              = "/src"

	cmdTmplFindBugs       = "java -jar %s -textui -xml:withMessages -effort:max -output %s %s"
	cmdTmplFindBugsText   = "java -jar %s -textui                   -effort:max -output %s %s"
	cmdTmplCheckstyle     = "java -jar %s -c %s -o %s -f xml %s"
	cmdTmplCheckstyleText = "java -jar %s -c %s -o %s %s"
)

// This line forces the compiler to check the method
// sets of the findbugsStep and checkstyleStep types
// to ensure that they both fulfill the javacmd interface
var _, _ javacmd = &findbugsStep{}, &checkstyleStep{}

func NewFindbugsStep(jarLoc, outputLoc, srcDir string, textoutput bool, logger *logrus.Logger) pipeline.Step {
	return &findbugsStep{
		jarLoc:    jarLoc,
		outputLoc: outputLoc,
		srcDir:    srcDir,
		text:      textoutput,
		log:       logger,
	}
}

func NewCheckstyleStep(jarLoc, outputLoc, srcDir, checkLoc, repoBase string, text bool, logger *logrus.Logger) pipeline.Step {
	return &checkstyleStep{
		srcDir:    srcDir,
		jarLoc:    jarLoc,
		outputLoc: outputLoc,
		checkLoc:  checkLoc,
		repoBase:  repoBase,
		text:      text,
		log:       logger,
	}
}

func (fb *findbugsStep) init(request *pipeline.Request) error {

	if err := fb.setSrcDir(request); err != nil {
		return err
	}

	if fb.jarLoc == "" {
		fb.jarLoc = DefaultFindBugsJarLoc
	}
	if fb.srcDir == "" {
		fb.srcDir = DefaultSrcDir
	}
	if fb.outputLoc == "" {
		fb.outputLoc = DefaultFindBugsOutputLoc
	}

	return nil
}

func (checkstyle *checkstyleStep) init(request *pipeline.Request) error {

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

	if checkstyle.outputLoc == "" {
		checkstyle.outputLoc = DefaultCheckstyleOutputLoc
	}

	if checkstyle.checkLoc == "" {
		checkstyle.checkLoc = DefaultCheckstyleConfigLoc
	}

	return nil
}

func (fb *findbugsStep) setSrcDir(request *pipeline.Request) error {

	if fb.srcDir != "" {
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
	fb.srcDir = srcDir
	return nil
}

func (checkstyle *checkstyleStep) setSrcDir(request *pipeline.Request) error {

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

func (fb *findbugsStep) launchCmd() (string, error) {

	cmd := fb.Cmd()
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadFile(fb.outputLoc)
	return string(contents), err
}

func (checkstyle *checkstyleStep) launchCmd() (string, error) {

	cmd := checkstyle.Cmd()
	stderr, err := cmd.StdoutPipe()
	if err != nil {
		checkstyle.log.Fatal(err)
	}
	if err = cmd.Start(); err != nil {
		checkstyle.log.Fatal(err)
	}
	if err = cmd.Wait(); err != nil {
		checkstyle.log.Fatal(err)
	}

	checkstyle.log.Info("Program has finished running")
	if err != nil {
		checkstyle.log.Info("Error is not nil")
		errorMessage, err := ioutil.ReadAll(stderr)
		if err != nil {
			checkstyle.log.Warn("Failing to correctly marshal the error!")
			checkstyle.log.Fatal(err)
		}
		checkstyle.log.Info("Logging the result string")
		checkstyle.log.Warn(string(errorMessage))
		return "", err
	}

	checkstyle.listFiles()
	contents, err := ioutil.ReadFile(checkstyle.outputLoc)
	checkstyle.log.Infof("Logging output of file: %v", string(contents))
	return string(contents), err
}

// TODO delete, only used for debugging purposes.
func (checkstyle *checkstyleStep) listFiles() {
	files, err := ioutil.ReadDir(checkstyle.srcDir)
	if err != nil {
		checkstyle.log.Fatalf("Error while listing directory: %v", err)
	}

	str := ""
	for _, f := range files {
		str += f.Name() + "\n"
	}
	checkstyle.log.Info("Files:\n%v", str)
}

func (fb *findbugsStep) Exec(request *pipeline.Request) *pipeline.Result {

	// Ensure all data is set
	if err := fb.init(request); err != nil {
		return &pipeline.Result{Error: err}
	}

	// Now, launch the command
	contents, err := fb.launchCmd()
	nextMap := fromMap(request.KeyVal)
	nextMap["findbugs"] = contents

	return &pipeline.Result{
		Error:  err,
		KeyVal: nextMap,
	}
}

func (checkstyle *checkstyleStep) Exec(request *pipeline.Request) *pipeline.Result {

	// Ensure all data is set
	if err := checkstyle.init(request); err != nil {
		return &pipeline.Result{Error: err}
	}

	// Now, launch the command
	contents, err := checkstyle.launchCmd()
	if err != nil {
		return &pipeline.Result{Error: err}
	}

	// Serialize the command resuls into a struct
	// TODO this should be done by piping STDOUT into a
	// XMLDecoder, NOT by writing out to a file and then reading that file back in.
	// It's WAY less effecient to write to disk, read from disk into memory, and then decode
	// First, I need to refactor "launchCmd" to be cleaner before I can do that, though.
	ck, err := checkstyle.serialize(contents)
	// TODO next, we need to filter out the file paths into something useful
	// We can't use the absolute path because that contains the temporary directory as a base
	// Iterate through all of the files and cut out the base path.
	ck = checkstyle.filterPath(ck)

	nextMap := fromMap(request.KeyVal)
	nextMap["checkstyle"] = ck

	return &pipeline.Result{
		Error:  err,
		KeyVal: nextMap,
	}
}

// filterPath walks through each file and removes the base of the path
// from the File.Name property. This is a transformation which makes it possible to
// post-back comments to GitHub. GitHub only know the path from the base
// of the directory, not the absolute path on the machine's filesystem.
func (checkstyle *checkstyleStep) filterPath(ch *Checkstyle) *Checkstyle {
	checkstyle.log.Info("Filtering the file paths...")
	for index, f := range ch.File {
		base, err := filepath.Abs(checkstyle.repoBase)
		if err != nil {
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
	return ch
}

func (checkstyle *checkstyleStep) serialize(blob string) (*Checkstyle, error) {
	var (
		check   Checkstyle
		decoder = xml.NewDecoder(strings.NewReader(blob))
		err     = decoder.Decode(&check)
	)
	return &check, err
}

func (fb *findbugsStep) Cancel() error {
	fb.Status("Cancel")
	return nil
}

func (checkstyle *checkstyleStep) Cancel() error {
	checkstyle.Status("Cancel")
	return nil
}

func (fb *findbugsStep) Cmd() *exec.Cmd {
	var strTmpl string = cmdTmplFindBugs

	if fb.text {
		strTmpl = cmdTmplFindBugsText
	}

	cmd := fmt.Sprintf(
		strTmpl,
		fb.jarLoc,
		fb.outputLoc,
		fb.srcDir,
	)

	return exec.Command("bash", "-c", cmd)
}

func (checkstyle *checkstyleStep) Cmd() *exec.Cmd {
	var strTmpl = cmdTmplCheckstyle

	if checkstyle.text {
		strTmpl = cmdTmplCheckstyleText
	}
	cmd := fmt.Sprintf(
		strTmpl,
		checkstyle.jarLoc,
		checkstyle.checkLoc,
		checkstyle.outputLoc,
		checkstyle.srcDir,
	)

	return exec.Command("bash", "-c", cmd)
}

func fromMap(m map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	if m == nil {
		return result
	}
	for key, val := range m {
		result[key] = val
	}
	return result
}
