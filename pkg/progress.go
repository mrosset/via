package via

import (
	"fmt"
	// "github.com/mrosset/util/human"
	"github.com/whyrusleeping/progmeter"
	"io"
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
	width := 15
	progress := (width * percent) / 100
	bar := strings.Repeat("#", int(progress))
	// bps := float64(pw.done) / time.Now().Sub(pw.start).Seconds()
	// speed := fmt.Sprintf("%03d%% %s/s", percent, human.ByteSize(bps))
	info := fmt.Sprintf("%-*s", width, bar)
	// out := speed + info[len(speed):]
	// info := fmt.Sprintf("%s/s", human.ByteSize(bps))
	pw.pm.Working(pw.key, info)
	// time.Sleep(10000000)
	return pw.w.Write(b)
}

func (pw *ProgressWriter) Close() error {
	return nil
}
