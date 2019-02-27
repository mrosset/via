package via

import (
	"fmt"
	"github.com/git-lfs/git-lfs/tools/humanize"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ProgressItem provides item for Progress type
type ProgressItem struct {
	key    string
	status string
	state  string
	done   bool
	sync.Mutex
}

// Update items state status and done state
func (pi *ProgressItem) Update(state, status string, done bool) {
	pi.Lock()
	pi.state = state
	pi.status = status
	pi.done = done
	pi.Unlock()
}

// Progress provides type for displaying concurrent
type Progress struct {
	items []*ProgressItem
	stop  chan bool
	sync.Mutex
	out io.Writer
}

func (p *Progress) print() {
	var (
		tdone = 0
	)
	if len(p.items) == 0 {
		return
	}
	fmt.Fprintf(p.out, "\033[H\033[2J")
	for _, it := range p.items {
		if it.done {
			tdone++
		}
		fmt.Fprintf(p.out, "%-10s %-13s %s\n", it.state, it.status, it.key)
	}
	fmt.Fprintf(p.out, "[%d/%d]\n", tdone, len(p.items))
	os.Stdout.Sync()
}

// Update updates item by key setting status and done stat
func (p Progress) Update(key, state, status string, done bool) {
	for _, i := range p.items {
		if i.key == key {
			i.Update(state, status, done)
		}
	}
}

// Start starts printing progress
func (p Progress) Start() {
	time.Sleep(time.Second / 4)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			p.print()
			select {
			case <-p.stop:
				return
			default:
			}
		}
	}()
}

// Done stops printing progress
func (p *Progress) Done() {
	p.stop <- true
}

// Add a progress item
func (p *Progress) Add(key string, state string) {
	item := &ProgressItem{key: key, state: state}
	p.Lock()
	p.items = append(p.items, item)
	p.Unlock()
}

// ProgressWriter provides a writer interface that updates speed and
// progress information for Progress type
type ProgressWriter struct {
	total int64
	w     io.Writer
	done  int64
	start time.Time
	pm    *Progress
	key   string
}

// NewProgressWriter returns a new ProgressWriter that has been initialized
func NewProgressWriter(pm *Progress, key string, t int64, w io.Writer) *ProgressWriter {
	return &ProgressWriter{
		pm:    pm,
		key:   key,
		total: t,
		w:     w}
}

// Write provides writer interface write method
func (pw *ProgressWriter) Write(b []byte) (int, error) {
	if pw.done == 0 {
		pw.start = time.Now()
	}
	pw.done += int64(len(b))
	percent := int((pw.done * 100) / pw.total)
	width := 8
	progress := (width * percent) / 100
	bar := fmt.Sprintf("[%-*s]", width, strings.Repeat("#", int(progress)))
	bps := humanize.FormatByteRate(uint64(pw.done), time.Now().Sub(pw.start))
	speed := fmt.Sprintf(" %s %s%%", bps, strconv.Itoa(percent))
	pw.pm.Update(pw.key, bar, speed, false)
	// time.Sleep(2.4e6)
	return pw.w.Write(b)
}
