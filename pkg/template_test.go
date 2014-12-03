package via

import (
	"os"
	"testing"
	"text/template"
)

var turl = "http://stuff/plan-{{.Version}}.tar.gz"

func TestTemplate(t *testing.T) {
	tmpl, err := template.New("").Parse(turl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, testPlan)
	if err != nil {
		panic(err)
	}
}
