package changelog

import (
	"bytes"
	"testing"
	"text/template"
)

func TestPretext(t *testing.T) {
	project := struct {
		Name    string
		Version string
		Date    string
	}{
		Date:    "Test Date",
		Name:    "Test",
		Version: "1.0.0",
	}

	var b bytes.Buffer

	tmpl := template.Must(template.New("").Parse(pretext))
	err := tmpl.Execute(&b, project)
	if err != nil {
		t.Errorf("could not complete templating of pretext: %v", err)
	}
}
