package jobs

type Checkstyle struct {
	File    []File `xml:"file"`
	Version string `xml:"version,attr"`
}

type File struct {
	Name  string  `xml:"name,attr"`
	Error []Error `xml:"error"`
}
type Error struct {
	Severity string `xml:"severity,attr"`
	Message  string `xml:"message,attr"`
	Source   string `xml:"source,attr"`
	Line     string `xml:"line,attr"`
}
