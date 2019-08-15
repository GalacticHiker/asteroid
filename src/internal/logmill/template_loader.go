package logmill

import (
	"log"
	"text/template"
)

// The indentation of the {{end}} must not be aligned - as it isi 'in' the templat definition
var defaultTemplates = [...]string{

	`{{define "defaultKVP"}}application={{.AppName}}, pid={{.Pid}}, time={{.LogTime.Format "2006-01-02T15:04:05.000Z07:00" }}, protocol={{.Protocol}}, tag={{.Tag}}, seqno={{.Seqno}}, sourceName={{.SourceHostname}}, sourceIP={{.SourceIP}}, level={{.|level}}, size={{.LogLength}}, template={{.Template}}, padding={{.|padding}}
{{end}}`,

	`{{define "defaultCSV"}}{{.AppName}}, {{.Pid}}, {{.LogTime.Format "2006-01-02T15:04:05.000Z07:00" }}, {{.SourceIP}}, {{.Seqno}}, {{.Protocol}}, {{.|level}}, {{.LogLength}}, {{.Template}}
{{end}}`,

	`{{define "short"}}{{.AppName}}, {{.Pid}}, {{.Seqno}}, {{.ConstTag}}, {{.Template}}
{{end}}`,

}

// LogTemplates - need to access this after creation
var LogTemplates *template.Template

// LoadTemplates loads the hardcoded templates and any templates specified in the templateGlob.
// Note that hard-coded templates and templates from template files are "associated" to "allLogTemplates".
// When executing templates are referenced by name - the "associated templates" capabilities are not used. 
// TODO: don't use "associated templates" as a list of templates (it screws up testing). Log templates should
// be in thier own file.
func LoadTemplates(templateGlob string) *template.Template {

	fmap := template.FuncMap{
		"padding": Padding,
		"level":   Level,
	}

	// the "empty"  root for all associated templates	
	LogTemplates = template.New("allLogTemplates").Funcs(fmap)

	for _, t := range defaultTemplates {
		LogTemplates = template.Must(LogTemplates.Parse(t))
	}

	if templateGlob != "" {
		_, err := LogTemplates.Funcs(fmap).ParseGlob(templateGlob)
		if err != nil {
			log.Printf("Warning parsing template files:%s, Error=%v\n", templateGlob, err)
		}
	}

	return LogTemplates
}

// TemplateExists - return true if a template with the name exists 
func TemplateExists( logTemplateName string) bool {

	if LogTemplates.Lookup(logTemplateName) == nil {
		return false
	}
	return true
}