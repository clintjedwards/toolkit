package main

import (
	"bytes"
	"testing"
	"text/template"
)

func TestPretext(t *testing.T) {
	project := projectInfo{
		Date:     "Test Date",
		Name:     "Test",
		Version:  "1.0.0",
		Username: "someuser",
	}

	var b bytes.Buffer

	tmpl := template.Must(template.New("").Parse(pretext))
	err := tmpl.Execute(&b, project)
	if err != nil {
		t.Error("could not templat pretext")
	}
}
