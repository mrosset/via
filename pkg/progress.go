package via

import (
	"fmt"
	"github.com/mrosset/util/human"
	"github.com/whyrusleeping/progmeter"
	"io"
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
	bps := float64(pw.done) / time.Now().Sub(pw.start).Seconds()
	info := fmt.Sprintf("%03d%%%s/s", percent, human.ByteSize(bps))
	pw.pm.Working(pw.key, info)
	return pw.w.Write(b)
}

func (pw *ProgressWriter) Close() error {
	return nil
}
