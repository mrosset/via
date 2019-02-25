package via

import (
	"fmt"
	"github.com/git-lfs/git-lfs/tools/humanize"
	"io"
	"strconv"
	"strings"
	"time"
)

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
	time.Sleep(2.4e6)
	return pw.w.Write(b)
}
