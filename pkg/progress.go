package via

import (
	"fmt"
	"github.com/mrosset/progmeter"
	"github.com/mrosset/util/human"
	"io"
	"strconv"
	"strings"
	"time"
)

type ProgressWriter struct {
	total int64
	w     io.Writer
	done  int64
	start time.Time
	pm    *progmeter.ProgMeter
	key   string
}

func NewProgressWriter(pm *progmeter.ProgMeter, key string, t int64, w io.Writer) *ProgressWriter {
	return &ProgressWriter{
		pm:    pm,
		key:   key,
		total: t,
		w:     w}
}

func (pw *ProgressWriter) Write(b []byte) (int, error) {
	if pw.done == 0 {
		pw.start = time.Now()
	}
	pw.done += int64(len(b))
	percent := int((pw.done * 100) / pw.total)
	width := 10
	progress := (width * percent) / 100
	bps := float64(pw.done) / time.Now().Sub(pw.start).Seconds()
	bar := fmt.Sprintf("%-*s", width, strings.Repeat("#", int(progress)))
	speed := fmt.Sprintf("%s/s %3.3s%%", human.ByteSize(bps), strconv.Itoa(percent))
	pw.pm.Working(pw.key, bar, speed)
	// time.Sleep(2.4e7)
	return pw.w.Write(b)
}

func (pw *ProgressWriter) Close() error {
	return nil
}
