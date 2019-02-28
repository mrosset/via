package via

import (
	"bufio"
	"fmt"
	"github.com/mrosset/util/file"
	"io"
	"os"
	"sync"
)

// Batch Plan type
type Batch struct {
	plans  PlanSlice
	config *Config
	size   int64
	ch     chan bool
	wg     *sync.WaitGroup
	keys   []string
	pm     *Progress
}

// NewBatch returns a new Batch that has been initialized
func NewBatch(conf *Config, out io.Writer) Batch {
	threads := conf.Threads
	if threads == 0 {
		elog.Println("threads are too low 0. setting it to 4. please fix config.json")
		threads = 4
	}
	return Batch{
		config: conf,
		ch:     make(chan bool, threads),
		wg:     new(sync.WaitGroup),
		pm: &Progress{
			stop: make(chan bool),
			out:  out,
		},
	}
}

// Plans returns the Batch's plans
func (b Batch) Plans() PlanSlice {
	return b.plans
}

// PruneInstalled Plans within the Batch
//
//FIXME: this is not implemented and is not currently being used
func (b *Batch) PruneInstalled() {
	for _, p := range b.plans {
		if IsInstalled(b.config, p.Name) {
			panic("Not Implemented")
		}
	}
}

// Add a 'Plan to the Batch
func (b *Batch) Add(plan *Plan) {
	b.plans = append(b.plans, plan)
	b.size += plan.Size
}

// Walk the plan's dependency tree and add each dependency to the
// batch if it does not already exist
func (b *Batch) Walk(plan *Plan) error {
	if b.config == nil {
		return fmt.Errorf("config is nil")
	}
	b.Add(plan)
	for _, d := range plan.Depends() {
		p, err := NewPlan(b.config, d)
		if err != nil {
			return err
		}
		if b.plans.Contains(p) {
			continue
		}
		b.Walk(p)
	}
	return nil
}

// ToInstall Returns a string slice of 'Plans to install
func (b *Batch) ToInstall() PlanSlice {
	var plans PlanSlice
	for _, p := range b.plans {
		if !IsInstalled(b.config, p.Name) {
			plans = append(plans, p)
		}
	}
	return plans
}

// ToDownload returns a string slice of Plans to download
func (b *Batch) ToDownload() []string {
	s := []string{}
	for _, p := range b.plans {
		if !PackageFileExists(b.config, p) && !IsInstalled(b.config, p.Name) {
			s = append(s, p.Name)
		}
	}
	return s
}

// Download plan's binary tarball for configured ipfs gateway
func (b Batch) Download(plan *Plan) error {
	var (
		pfile = PackagePath(b.config, plan)
		url   = ""
	)
	if isDocker() {
		url = "http://172.17.0.1:8080/ipfs/" + plan.Cid
	} else {
		url = b.config.Binary + "/" + plan.Cid
	}
	b.config.Repo.Ensure()
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
	defer fd.Close()
	pw := NewProgressWriter(b.pm, plan.Name, res.ContentLength, fd)
	_, err = io.Copy(pw, res.Body)
	return err
}

// PlanFunc provides a func that takes a Plan
type PlanFunc func(*Plan) error

// DownloadInstall provides a Plan function that downloads and
// installs a Plan
func (b *Batch) DownloadInstall(plan *Plan) error {
	b.pm.Update(plan.Name, "download", "", false)
	if err := b.Download(plan); err != nil {
		return err
	}
	b.pm.Update(plan.Name, "install", "", false)
	if err := NewInstaller(b.config, plan).Install(); err != nil {
		elog.Fatalf("%s: %s", plan.Name, err)
		return fmt.Errorf("%s: %s", plan.Name, err)
	}
	b.pm.Update(plan.Name, "done", "", true)
	return nil
}

// ForEach run PlanFunc on each plan in Plans.
func (b Batch) ForEach(fn PlanFunc, plans PlanSlice) (errors []error) {
	for _, p := range plans {
		b.pm.Add(p.Name, "wait")
		go func(plan *Plan) {
			b.wg.Add(1)
			b.ch <- true
			if err := fn(plan); err != nil {
				errors = append(errors, err)
			}
			b.wg.Done()
			<-b.ch
		}(p)
	}
	b.pm.Start()
	b.wg.Wait()
	b.pm.Done()
	return errors
}

// Install does the final download and installing of the batch plans
func (b *Batch) Install() (errors []error) {
	return b.ForEach(b.DownloadInstall, b.ToInstall())
}

// PromptInstall prompts user before installing
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
	//	bf := new(bytes.Buffer)
	//	st := struct {
	//		Install  []string
	//		Download []string
	//		Size     string
	//	}{b.ToInstall(), b.ToDownload(), human.ByteSize(b.size).String()}

	//	ot := `
	// Installing:
	// {{.Install}}

	// Downloading:
	// {{.Download}}

	// Install Size: {{.Size}}

	// `
	//	tmpl, err := template.New("output").Parse(ot)

	//	if err != nil {
	//		panic(err)
	//	}
	//	err = tmpl.Execute(bf, st)
	//	if err != nil {
	//		panic(err)
	//	}
	//	return bf.String()
	return "Not Implimented"
}
