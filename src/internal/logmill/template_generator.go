package logmill

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// TemplateGenerator data
type TemplateGenerator struct {
	logTemplate     *template.Template
	templateContext *TemplateContext
}

// LogGenerator is the interface
type LogGenerator interface {
	GenerateLog(logTime time.Time) (logText string, eof bool)
	TemplateContext() *TemplateContext
	Rewind()
	io.Closer
}

// NewTemplateGenerator creation
func NewTemplateGenerator(logTemplateName string) LogGenerator {

	lg := new(TemplateGenerator)
	lg.logTemplate = LoadTemplates("../conf/example.template")

	if !TemplateExists(logTemplateName) {
		log.Fatalf("ERROR Log Template \"%s\" not loaded.", logTemplateName)
	}

	appName := "unknown"
	appPath, err := os.Executable()
	if err != nil {
		fmt.Printf("WARNING: unable to get AppName, %v\n", err)
	} else {
		appName = filepath.Base(appPath)
	}

	lg.templateContext = NewTemplateContext(appName, logTemplateName)
	return lg
}

// GenerateLog create a log
func (lg *TemplateGenerator) GenerateLog(logTime time.Time) (logText string, eof bool) {

	lg.templateContext.LogTextBuffer.Reset()
	lg.templateContext.LogTime = logTime

	err := lg.logTemplate.ExecuteTemplate(lg.templateContext.LogTextBuffer, lg.templateContext.Template, lg.templateContext)

	if err != nil {
		panic(err) /// TODO: improve?
	}

	logText = (*lg.templateContext.LogTextBuffer).String()
	lg.templateContext.Seqno++
	return logText, false

}

// Close generator -- when log are generated from a file
func (lg *TemplateGenerator) Close() error {

	return nil
}

// Rewind resets the generator to start over from the beginning
func (lg *TemplateGenerator) Rewind() {
	lg.templateContext.Seqno = 0
}

// TemplateContext returns the TemplateContext of the generator
func (lg *TemplateGenerator) TemplateContext() *TemplateContext {
	return lg.templateContext
}
