package jobs

type Checkstyle struct {
	File []struct {
		Name  string `xml:"name,attr"`
		Error []struct {
			Severity string `xml:"severity,attr"`
			Message  string `xml:"message,attr"`
			Source   string `xml:"source,attr"`
			Line     string `xml:"line,attr"`
		} `xml:"error"`
	} `xml:"file"`
	Version string `xml:"version,attr"`
}
