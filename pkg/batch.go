package via

import (
	"bytes"
	"github.com/mrosset/util/file"
	"text/template"
)

// Batch Plan type
type Batch struct {
	Plans  map[string]*Plan
	config *Config
}

// Creates a new Batch type
func NewBatch(config *Config) Batch {
	return Batch{
		Plans:  make(map[string]*Plan),
		config: config,
	}
}

// Prunes Installed Plans within the Batch
func (b *Batch) PruneInstalled() {
	for i, _ := range b.Plans {
		if IsInstalled(i) {
			delete(b.Plans, i)
		}
	}
}

// Adds 'Plane to the Batch
func (b *Batch) Add(plan *Plan) {
	b.Plans[plan.Name] = plan
}

func (b *Batch) Walk(plan *Plan) {
	b.Add(plan)
	for _, d := range plan.Depends() {
		p, _ := NewPlan(d)
		if _, ok := b.Plans[p.Name]; ok {
			continue
		}
		b.Walk(p)
	}
}

// Returns a string slice of 'Plans to install
func (b *Batch) ToInstall() []string {
	s := []string{}
	for i, _ := range b.Plans {
		if !IsInstalled(i) {
			s = append(s, i)
		}
	}
	return s
}

// Returns a string slice of 'Plans to download
func (b *Batch) ToDownload() []string {
	s := []string{}
	for i, p := range b.Plans {
		if !file.Exists(p.PackageFile()) {
			s = append(s, i)
		}
	}
	return s
}

// Provides stringer interface
func (b Batch) String() string {
	bf := new(bytes.Buffer)
	st := struct {
		Install  []string
		Download []string
		Size     int
	}{b.ToInstall(), b.ToDownload(), 100}

	ot := `
Installing:
{{.Install}}

Downloading:
{{.Download}}

Install Size: {{.Size}}
`
	tmpl, err := template.New("output").Parse(ot)

	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(bf, st)
	if err != nil {
		panic(err)
	}
	return bf.String()
}
