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
	ch     chan bool
	wg     *sync.WaitGroup
}

// Creates a new Batch type
func NewBatch(conf *Config) Batch {
	threads := conf.Threads
	if threads == 0 {
		fmt.Println(conf)
		elog.Fatal("threads are too low 0. setting it to 4. please fix config.json")
		threads = 4
	}
	return Batch{
		Plans:  make(map[string]*Plan),
		config: conf,
		pm:     progmeter.NewProgMeter(false),
		ch:     make(chan bool, threads),
		wg:     new(sync.WaitGroup),
	}
}

// Prunes Installed Plans within the Batch
func (b *Batch) PruneInstalled() {
	for i, _ := range b.Plans {
		if IsInstalled(b.config, i) {
			delete(b.Plans, i)
		}
	}
}

// Adds 'Plane to the Batch
func (b *Batch) Add(plan *Plan) {
	b.Plans[plan.Name] = plan
	b.size += plan.Size
}

func (b *Batch) Walk(plan *Plan) error {
	b.Add(plan)
	for _, d := range plan.Depends() {
		p, err := NewPlan(b.config, d)
		if err != nil {
			return err
		}
		if _, ok := b.Plans[p.Name]; ok {
			continue
		}
		b.Walk(p)
	}
	return nil
}

// Returns a string slice of 'Plans to install
func (b *Batch) ToInstall() []string {
	s := []string{}
	for i, _ := range b.Plans {
		if !IsInstalled(b.config, i) {
			s = append(s, i)
		}
	}
	return s
}

// Returns a string slice of 'Plans to download
func (b *Batch) ToDownload() []string {
	s := []string{}
	for i, p := range b.Plans {
		if !file.Exists(p.PackagePath()) && !IsInstalled(b.config, p.Name) {
			s = append(s, i)
		}
	}
	return s
}

func (b Batch) Download(plan *Plan) error {
	var (
		pfile = plan.PackagePath()
		url   = ""
	)
	if isDocker() {
		url = "http://172.17.0.1:8080/ipfs/" + plan.Cid
	} else {
		url = b.config.Binary + "/" + plan.Cid
	}

	if !file.Exists(b.config.Repo) {
		os.MkdirAll(b.config.Repo, 0755)
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

type PlanFunc func(*Plan)

func (b *Batch) downloadInstall(plan *Plan) {
	b.ch <- true
	b.pm.AddEntry(plan.Name, "         "+plan.NameVersion(), "")
	defer func() { <-b.ch }()
	defer b.wg.Done()
	if err := b.Download(plan); err != nil {
		b.pm.Error(plan.Name, err.Error())
		elog.Fatal(err)
		return
	}
	b.pm.Working(plan.Name, "install       ")
	in := NewInstaller(b.config, plan)
	if err := in.Install(); err != nil {
		b.pm.Error(plan.Name, err.Error())
		elog.Fatal(err)
		return
	}
	b.pm.Finish(plan.Name)
}

func (b Batch) ForEach(fn PlanFunc) (errors []error) {
	for _, n := range b.ToInstall() {
		p, err := NewPlan(b.config, n)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		b.wg.Add(1)
		b.pm.AddTodos(1)
		go fn(p)
	}
	b.wg.Wait()
	b.pm.MarkDone()
	fmt.Println()
	return errors
}

func (b *Batch) Install() (errors []error) {
	return b.ForEach(b.downloadInstall)
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
