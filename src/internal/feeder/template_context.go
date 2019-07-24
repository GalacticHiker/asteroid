package feeder

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// TemplateContext contains the data item available in the templates
type TemplateContext struct {
	LogTime        time.Time
	Seqno          int64
	AppName        string
	SourceHostname string
	SourceIP       string
	Pid            int
	LogLength      int
	Protocol       string
	Template       string
	LogTextBuffer  *bytes.Buffer
}

// NewTemplateContext creates a new TemplatContext and initalizes it with configuration values
func NewTemplateContext(appName, logTemplate string) *TemplateContext {

	templateContext := new(TemplateContext)

	templateContext.LogTime = time.Now()

	var logBuffer bytes.Buffer
	templateContext.LogTextBuffer = &logBuffer

	templateContext.Pid = os.Getpid()

	templateContext.AppName = appName

	srchostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("WARNING: %v\n", err)
		srchostname = "unknown"
	}
	templateContext.SourceHostname = srchostname

	templateContext.SourceIP = getLocalIP()

	templateContext.Template = logTemplate

	templateContext.LogLength = 300
	templateContext.Protocol = "logdna"

	return templateContext
}

// Padding template value
func Padding(templateContext TemplateContext) (string, error) {

	var b bytes.Buffer
	b = *templateContext.LogTextBuffer

	s := ""
	padLength := templateContext.LogLength - b.Len()
	if padLength > 0 {
		s = pad(b.Len(), int(templateContext.Seqno), padLength)
	}

	return s, nil
}

// pad processes {{.padding}}. padding is ' Ar-nnn' where r is the row, nnn is the character offset in the log buffer,guarantee the log terminates with ' Z...\n'
func pad(padoffset, seqno, padLength int) (padding string) {

	alphachars := "ABCDEFGHIJKLMNNOPQRSTUVWXYZ"
	rowAlpha := alphachars[seqno%len(alphachars) : seqno%len(alphachars)+1]

	nextWritePosition := 0
	padWord := rowAlpha + strconv.FormatInt(int64(seqno), 10) + "-" + strconv.FormatInt(int64(padoffset), 10) // omit leading space on first word

	var padbuf bytes.Buffer
	for {

		padwordlen := len(padWord)

		residue := padLength - nextWritePosition

		if residue == 0 {
			break
		}

		// the last character should never be a space
		trailChars := " XXXXXXXXXXXXXXXXXXXXXXXXXXX"
		if residue <= padwordlen+2 {
			t := trailChars[:residue]
			if padLength-residue == 0 {
				t = trailChars[1:residue] // skip leading blank.
			}
			padbuf.WriteString(t)
			break
		}

		padbuf.WriteString(padWord)

		// account for missing leading ' ' in first pad word.
		if nextWritePosition == 0 {
			nextWritePosition += padwordlen + 1
			padoffset += padwordlen + 1
		} else {
			nextWritePosition += padwordlen
			padoffset += padwordlen
		}

		padWord = " " + rowAlpha + strconv.FormatInt(int64(seqno), 10) + "-" + strconv.FormatInt(int64(padoffset), 10)

	}

	return padbuf.String()

}

// Level provides a random log level
func Level(templateContext TemplateContext) (string, error) {

	levels := [...]string{"INFO", "WARNING", "ERROR", "DEBUG"}

	level := levels[rand.Intn(len(levels))]

	return level, nil
}
