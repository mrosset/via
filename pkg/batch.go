package via

import (
	"bytes"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/file"
	"os"
	"sync"
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
		if !file.Exists(p.PackageFilePath(config)) {
			s = append(s, i)
		}
	}
	return s
}

func (b *Batch) Download() []error {
	rdir := join(b.config.Repo, "repo")
	if !file.Exists(rdir) {
		os.MkdirAll(rdir, 0755)
	}
	errors := []error{}
	for _, p := range b.ToDownload() {
		plan := b.Plans[p]
		err := gurl.NameDownload(rdir, b.config.Binary+"/"+plan.Cid, plan.PackageFile())
		if err != nil {
			errors = append(errors, err)
		}

	}
	return errors
}

func (b *Batch) Install() (errors []error) {
	wg := new(sync.WaitGroup)
	for _, n := range b.ToInstall() {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			if err := Install(name); err != nil {
				errors = append(errors, err)
			}
		}(n)
	}
	wg.Wait()
	return errors
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
