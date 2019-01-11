package via

import (
	"bytes"
	"github.com/mrosset/util/file"
	"github.com/whyrusleeping/progmeter"
	"io"
	"os"
	"runtime"
	"sync"
	"text/template"
)

// Batch Plan type
type Batch struct {
	Plans  map[string]*Plan
	config *Config
	pm     *progmeter.ProgMeter
}

// Creates a new Batch type
func NewBatch(config *Config) Batch {
	return Batch{
		Plans:  make(map[string]*Plan),
		config: config,
		pm:     progmeter.NewProgMeter(false),
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
		if !file.Exists(p.PackageFilePath(config)) && !IsInstalled(p.Name) {
			s = append(s, i)
		}
	}
	return s
}

func (b Batch) Download(plan *Plan) error {
	var (
		rdir  = join(config.Repo, "repo")
		pfile = plan.PackageFilePath(config)
		url   = config.Binary + "/" + plan.Cid
	)
	if !file.Exists(rdir) {
		os.MkdirAll(rdir, 0755)
	}
	if file.Exists(pfile) {
		return nil
	}
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	fd, err := os.Create(pfile)
	if err != nil {
		return err
	}
	pw := NewProgressWriter(b.pm, plan.Name, res.ContentLength, fd)
	defer fd.Close()
	_, err = io.Copy(pw, res.Body)
	pw.Close()
	return err
}

func (b *Batch) Install() (errors []error) {
	wg := new(sync.WaitGroup)
	ch := make(chan bool, runtime.NumCPU())
	b.pm.AddTodos(len(b.ToInstall()))
	for _, n := range b.ToInstall() {
		wg.Add(1)
		go func(p *Plan) {
			ch <- true
			defer wg.Done()

			b.pm.AddEntry(p.Name, p.Name, "          "+p.Cid)
			if err := b.Download(p); err != nil {
				b.pm.Error(p.Name, err.Error())
				errors = append(errors, err)
				return
			}
			b.pm.Working(p.Name, "install        ")
			if err := Install(p.Name); err != nil {
				errors = append(errors, err)
			}
			<-ch
			b.pm.Finish(p.Name)
		}(b.Plans[n])
	}
	wg.Wait()
	b.pm.MarkDone()
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
