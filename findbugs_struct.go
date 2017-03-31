package jobs

type BugCollection struct {
	Project           Project         `xml:"Project"`
	Version           string          `xml:"version,attr"`
	BugInstance       []BugInstance   `xml:"BugInstance"`
	Timestamp         string          `xml:"timestamp,attr"`
	Release           string          `xml:"release,attr"`
	FindBugsSummary   FindBugsSummary `xml:"FindBugsSummary"`
	ClassFeatures     string          `xml:"ClassFeatures"`
	History           string          `xml:"History"`
	AnalysisTimestamp string          `xml:"analysisTimestamp,attr"`
	Errors            Errors          `xml:"Errors"`
	Sequence          string          `xml:"sequence,attr"`
}

type FindBugsSummary struct {
	Total_bugs         string         `xml:"total_bugs,attr"`
	Referenced_classes string         `xml:"referenced_classes,attr"`
	Num_packages       string         `xml:"num_packages,attr"`
	Peak_mbytes        string         `xml:"peak_mbytes,attr"`
	Gc_seconds         string         `xml:"gc_seconds,attr"`
	Vm_version         string         `xml:"vm_version,attr"`
	Cpu_seconds        string         `xml:"cpu_seconds,attr"`
	Alloc_mbytes       string         `xml:"alloc_mbytes,attr"`
	Timestamp          string         `xml:"timestamp,attr"`
	Total_size         string         `xml:"total_size,attr"`
	Priority_2         string         `xml:"priority_2,attr"`
	Total_classes      string         `xml:"total_classes,attr"`
	Java_version       string         `xml:"java_version,attr"`
	Clock_seconds      string         `xml:"clock_seconds,attr"`
	Priority_1         string         `xml:"priority_1,attr"`
	PackageStats       []PackageStats `xml:"PackageStats"`
	ClassProfile       []ClassProfile `xml:"FindBugsProfile>ClassProfile"`
}
type SourceLine struct {
	Sourcepath string `xml:"sourcepath,attr"`
	Classname  string `xml:"classname,attr"`
	Start      string `xml:"start,attr"`
	End        string `xml:"end,attr"`
	Sourcefile string `xml:"sourcefile,attr"`
}
type Method struct {
	Role       string     `xml:"role,attr"`
	IsStatic   string     `xml:"isStatic,attr"`
	Classname  string     `xml:"classname,attr"`
	Name       string     `xml:"name,attr"`
	Signature  string     `xml:"signature,attr"`
	SourceLine SourceLine `xml:"SourceLine"`
}
type SourceLineBugInstance struct {
	Role          string `xml:"role,attr"`
	StartBytecode string `xml:"startBytecode,attr"`
	EndBytecode   string `xml:"endBytecode,attr"`
	Sourcefile    string `xml:"sourcefile,attr"`
	Sourcepath    string `xml:"sourcepath,attr"`
	Classname     string `xml:"classname,attr"`
	Start         string `xml:"start,attr"`
	End           string `xml:"end,attr"`
}
type Field struct {
	Classname       string     `xml:"classname,attr"`
	Name            string     `xml:"name,attr"`
	Signature       string     `xml:"signature,attr"`
	SourceSignature string     `xml:"sourceSignature,attr"`
	Role            string     `xml:"role,attr"`
	IsStatic        string     `xml:"isStatic,attr"`
	SourceLine      SourceLine `xml:"SourceLine"`
}
type Type struct {
	Role           string `xml:"role,attr"`
	TypeParameters string `xml:"typeParameters,attr"`
	Descriptor     string `xml:"descriptor,attr"`
}
type SourceLineTypeBugInstance struct {
	Start      string `xml:"start,attr"`
	End        string `xml:"end,attr"`
	Sourcefile string `xml:"sourcefile,attr"`
	Sourcepath string `xml:"sourcepath,attr"`
	Classname  string `xml:"classname,attr"`
}
type PackageStats struct {
	Total_bugs  string       `xml:"total_bugs,attr"`
	Total_types string       `xml:"total_types,attr"`
	Priority_1  string       `xml:"priority_1,attr"`
	Priority_2  string       `xml:"priority_2,attr"`
	Total_size  string       `xml:"total_size,attr"`
	Package     string       `xml:"package,attr"`
	ClassStats  []ClassStats `xml:"ClassStats"`
}
type BugInstance struct {
	Category              string        `xml:"category,attr"`
	Type                  string        `xml:"type,attr"`
	Priority              string        `xml:"priority,attr"`
	Rank                  string        `xml:"rank,attr"`
	Abbrev                string        `xml:"abbrev,attr"`
	Method                Method        `xml:"Method"`
	SourceLineBugInstance SourceLine    `xml:"SourceLine"`
	Int                   Int           `xml:"Int"`
	Property              Property      `xml:"Property"`
	Class                 Class         `xml:"Class"`
	String                String        `xml:"String"`
	LocalVariable         LocalVariable `xml:"LocalVariable"`
	Field                 Field         `xml:"Field"`
}
type SourceLineMethodBugInstance struct {
	Start         string `xml:"start,attr"`
	End           string `xml:"end,attr"`
	StartBytecode string `xml:"startBytecode,attr"`
	EndBytecode   string `xml:"endBytecode,attr"`
	Sourcefile    string `xml:"sourcefile,attr"`
	Sourcepath    string `xml:"sourcepath,attr"`
	Classname     string `xml:"classname,attr"`
}
type LocalVariable struct {
	Name     string `xml:"name,attr"`
	Register string `xml:"register,attr"`
	Pc       string `xml:"pc,attr"`
	Role     string `xml:"role,attr"`
}
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
type Project struct {
	ProjectName string `xml:"projectName,attr"`
	Jar         string `xml:"Jar"`
}
type ClassProfile struct {
	TotalMilliseconds                          string `xml:"totalMilliseconds,attr"`
	Invocations                                string `xml:"invocations,attr"`
	AvgMicrosecondsPerInvocation               string `xml:"avgMicrosecondsPerInvocation,attr"`
	MaxMicrosecondsPerInvocation               string `xml:"maxMicrosecondsPerInvocation,attr"`
	StandardDeviationMircosecondsPerInvocation string `xml:"standardDeviationMircosecondsPerInvocation,attr"`
	Name                                       string `xml:"name,attr"`
}
type Class struct {
	Classname  string     `xml:"classname,attr"`
	Role       string     `xml:"role,attr"`
	SourceLine SourceLine `xml:"SourceLine"`
}
type SourceLineFieldBugInstance struct {
	Classname  string `xml:"classname,attr"`
	Sourcefile string `xml:"sourcefile,attr"`
	Sourcepath string `xml:"sourcepath,attr"`
}
type String struct {
	Value string `xml:"value,attr"`
	Role  string `xml:"role,attr"`
}
type Int struct {
	Value string `xml:"value,attr"`
	Role  string `xml:"role,attr"`
}
type Errors struct {
	StackTrace     []string `xml:"Error>StackTrace"`
	ErrorMessage   string   `xml:"Error>ErrorMessage"`
	Exception      string   `xml:"Error>Exception"`
	MissingClass   string   `xml:"MissingClass"`
	Errors         string   `xml:"errors,attr"`
	MissingClasses string   `xml:"missingClasses,attr"`
}
type ClassStats struct {
	Priority_2 string `xml:"priority_2,attr"`
	Priority_1 string `xml:"priority_1,attr"`
	SourceFile string `xml:"sourceFile,attr"`
	Interface  string `xml:"interface,attr"`
	Size       string `xml:"size,attr"`
	Bugs       string `xml:"bugs,attr"`
	Class      string `xml:"class,attr"`
}
