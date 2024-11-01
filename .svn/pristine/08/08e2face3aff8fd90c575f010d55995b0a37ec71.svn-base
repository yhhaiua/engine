package job

import (
	"testing"
)

func TestJob(t *testing.T) {
	s := &TempJob{count: 2, t: t}
	Submit(s.Run)
	select {}
}

type TempJob struct {
	t     *testing.T
	count int
}

func (t *TempJob) Run() {
	t.t.Logf("Run:%d", t.count)
	t.count--
	if t.count > 0 {
		Submit(t.Run)
	}
}
