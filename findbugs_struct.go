package jobs

type bugcollection struct {
	project           project         `xml:"project"`
	version           string          `xml:"version,attr"`
	buginstance       []buginstance   `xml:"buginstance"`
	timestamp         string          `xml:"timestamp,attr"`
	release           string          `xml:"release,attr"`
	findbugssummary   findbugssummary `xml:"findbugssummary"`
	classfeatures     string          `xml:"classfeatures"`
	history           string          `xml:"history"`
	analysistimestamp string          `xml:"analysistimestamp,attr"`
	errors            errorList       `xml:"errors"`
	sequence          string          `xml:"sequence,attr"`
}

type findbugssummary struct {
	total_bugs         string         `xml:"total_bugs,attr"`
	referenced_classes string         `xml:"referenced_classes,attr"`
	num_packages       string         `xml:"num_packages,attr"`
	peak_mbytes        string         `xml:"peak_mbytes,attr"`
	gc_seconds         string         `xml:"gc_seconds,attr"`
	vm_version         string         `xml:"vm_version,attr"`
	cpu_seconds        string         `xml:"cpu_seconds,attr"`
	alloc_mbytes       string         `xml:"alloc_mbytes,attr"`
	timestamp          string         `xml:"timestamp,attr"`
	total_size         string         `xml:"total_size,attr"`
	priority_2         string         `xml:"priority_2,attr"`
	total_classes      string         `xml:"total_classes,attr"`
	java_version       string         `xml:"java_version,attr"`
	clock_seconds      string         `xml:"clock_seconds,attr"`
	priority_1         string         `xml:"priority_1,attr"`
	packagestats       []packagestats `xml:"packagestats"`
	classprofile       []classprofile `xml:"findbugsprofile>classprofile"`
}
type sourceline struct {
	sourcepath string `xml:"sourcepath,attr"`
	classname  string `xml:"classname,attr"`
	start      string `xml:"start,attr"`
	end        string `xml:"end,attr"`
	sourcefile string `xml:"sourcefile,attr"`
}
type method struct {
	role       string     `xml:"role,attr"`
	isstatic   string     `xml:"isstatic,attr"`
	classname  string     `xml:"classname,attr"`
	name       string     `xml:"name,attr"`
	signature  string     `xml:"signature,attr"`
	sourceline sourceline `xml:"sourceline"`
}
type sourcelinebuginstance struct {
	role          string `xml:"role,attr"`
	startbytecode string `xml:"startbytecode,attr"`
	endbytecode   string `xml:"endbytecode,attr"`
	sourcefile    string `xml:"sourcefile,attr"`
	sourcepath    string `xml:"sourcepath,attr"`
	classname     string `xml:"classname,attr"`
	start         string `xml:"start,attr"`
	end           string `xml:"end,attr"`
}
type field struct {
	classname       string     `xml:"classname,attr"`
	name            string     `xml:"name,attr"`
	signature       string     `xml:"signature,attr"`
	sourcesignature string     `xml:"sourcesignature,attr"`
	role            string     `xml:"role,attr"`
	isstatic        string     `xml:"isstatic,attr"`
	sourceline      sourceline `xml:"sourceline"`
}
type bugtype struct {
	role           string `xml:"role,attr"`
	typeparameters string `xml:"typeparameters,attr"`
	descriptor     string `xml:"descriptor,attr"`
}
type sourcelinetypebuginstance struct {
	start      string `xml:"start,attr"`
	end        string `xml:"end,attr"`
	sourcefile string `xml:"sourcefile,attr"`
	sourcepath string `xml:"sourcepath,attr"`
	classname  string `xml:"classname,attr"`
}
type packagestats struct {
	total_bugs  string       `xml:"total_bugs,attr"`
	total_types string       `xml:"total_types,attr"`
	priority_1  string       `xml:"priority_1,attr"`
	priority_2  string       `xml:"priority_2,attr"`
	total_size  string       `xml:"total_size,attr"`
	pkg         string       `xml:"package,attr"`
	classstats  []classstats `xml:"classstats"`
}
type buginstance struct {
	category              string        `xml:"category,attr"`
	bugtype               string        `xml:"type,attr"`
	priority              string        `xml:"priority,attr"`
	rank                  string        `xml:"rank,attr"`
	abbrev                string        `xml:"abbrev,attr"`
	method                method        `xml:"method"`
	sourcelinebuginstance sourceline    `xml:"sourceline"`
	integer               int           `xml:"int"`
	property              property      `xml:"property"`
	class                 class         `xml:"class"`
	stringtype            string        `xml:"string"`
	localvariable         localvariable `xml:"localvariable"`
	field                 field         `xml:"field"`
}
type sourcelinemethodbuginstance struct {
	start         string `xml:"start,attr"`
	end           string `xml:"end,attr"`
	startbytecode string `xml:"startbytecode,attr"`
	endbytecode   string `xml:"endbytecode,attr"`
	sourcefile    string `xml:"sourcefile,attr"`
	sourcepath    string `xml:"sourcepath,attr"`
	classname     string `xml:"classname,attr"`
}
type localvariable struct {
	name     string `xml:"name,attr"`
	register string `xml:"register,attr"`
	pc       string `xml:"pc,attr"`
	role     string `xml:"role,attr"`
}
type property struct {
	name  string `xml:"name,attr"`
	value string `xml:"value,attr"`
}
type project struct {
	projectname string `xml:"projectname,attr"`
	jar         string `xml:"jar"`
}
type classprofile struct {
	totalmilliseconds                          string `xml:"totalmilliseconds,attr"`
	invocations                                string `xml:"invocations,attr"`
	avgmicrosecondsperinvocation               string `xml:"avgmicrosecondsperinvocation,attr"`
	maxmicrosecondsperinvocation               string `xml:"maxmicrosecondsperinvocation,attr"`
	standarddeviationmircosecondsperinvocation string `xml:"standarddeviationmircosecondsperinvocation,attr"`
	name                                       string `xml:"name,attr"`
}
type class struct {
	classname  string     `xml:"classname,attr"`
	role       string     `xml:"role,attr"`
	sourceline sourceline `xml:"sourceline"`
}
type sourcelinefieldbuginstance struct {
	classname  string `xml:"classname,attr"`
	sourcefile string `xml:"sourcefile,attr"`
	sourcepath string `xml:"sourcepath,attr"`
}
type stringtype struct {
	value string `xml:"value,attr"`
	role  string `xml:"role,attr"`
}
type integer struct {
	value string `xml:"value,attr"`
	role  string `xml:"role,attr"`
}
type errorList struct {
	stacktrace     []string `xml:"error>stacktrace"`
	errormessage   string   `xml:"error>errormessage"`
	exception      string   `xml:"error>exception"`
	missingclass   string   `xml:"missingclass"`
	errors         string   `xml:"errors,attr"`
	missingclasses string   `xml:"missingclasses,attr"`
}
type classstats struct {
	priority_2    string `xml:"priority_2,attr"`
	priority_1    string `xml:"priority_1,attr"`
	sourcefile    string `xml:"sourcefile,attr"`
	interfacetype string `xml:"interface,attr"`
	size          string `xml:"size,attr"`
	bugs          string `xml:"bugs,attr"`
	class         string `xml:"class,attr"`
}
