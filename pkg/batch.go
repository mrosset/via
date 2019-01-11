package via

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/human"
	"github.com/whyrusleeping/progmeter"
	"io"
	"os"
	"sync"
	"text/template"
)

// Batch Plan type
type Batch struct {
	Plans  map[string]*Plan
	config *Config
	pm     *progmeter.ProgMeter
	size   int64
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
	b.size += plan.Size
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
	ch := make(chan bool, b.config.Threads)
	for _, n := range b.ToInstall() {
		wg.Add(1)
		go func(p *Plan) {
			defer wg.Done()
			b.pm.AddTodos(1)
			ch <- true
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
			b.pm.Finish(p.Name)
			<-ch
		}(b.Plans[n])
	}
	wg.Wait()
	b.pm.MarkDone()
	return errors
}

func (b Batch) PromptInstall() []error {
	fmt.Printf("%s", b)
	fmt.Printf("Install? y/n : ")
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	switch scan.Text() {
	case "y":
		return b.Install()
	default:
		return nil

	}
	return nil
}

// Provides stringer interface
func (b Batch) String() string {
	bf := new(bytes.Buffer)
	st := struct {
		Install  []string
		Download []string
		Size     string
	}{b.ToInstall(), b.ToDownload(), human.ByteSize(b.size).String()}

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
