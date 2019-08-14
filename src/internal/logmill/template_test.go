package logmill

import (
	"testing"
	"text/template"
)

// usage: go test ./internal/logmill/ -v

func TestTemplateLoader(t *testing.T) {

	templateContext := NewTemplateContext("HelloTest","")
	logTemplates := LoadTemplates("../../../runtime/conf/example.template")

	tl := logTemplates.Templates()
	for i, tmpl := range tl {
		t.Logf("logTemplate[%d]:%v\n", i, tmpl.Name())
	}

	for _, tmpl := range tl {

		templateContext.LogTextBuffer.Reset()
		templateContext.Template = tmpl.Name()
		if templateContext.Template == "allLogTemplates" {
			continue // HACK ALERT produces a log of length 0, see template_loader
		}
		testOne(logTemplates, templateContext, t)

	}

}
func testOne(logTemplates *template.Template, templateContext *TemplateContext, t *testing.T) {

	err := logTemplates.ExecuteTemplate(templateContext.LogTextBuffer, templateContext.Template, templateContext)
	if err != nil {
		t.Errorf("ExecuteTemplate:%s Error:%v\n", templateContext.Template, err)
	}

	logText := (*templateContext.LogTextBuffer).String()
	if len(logText) == 0 {
		t.Errorf("ExecuteTemplate:%s Error:log length==0\n", templateContext.Template)
	}

	t.Logf("%s:%s\n", templateContext.Template, logText)
}
