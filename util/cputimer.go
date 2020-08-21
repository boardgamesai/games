package util

import (
	"syscall"
	"time"
)

type CPUTimer struct {
	start syscall.Timeval
}

func NewCPUTimer() *CPUTimer {
	return &CPUTimer{
		start: getCPUTime(),
	}
}

// Returns microseconds
func (t *CPUTimer) Elapsed() int64 {
	now := getCPUTime()

	start := time.Unix(t.start.Sec, int64(t.start.Usec*1000))
	end := time.Unix(now.Sec, int64(now.Usec*1000))
	elapsed := end.Sub(start).Microseconds()

	t.start = now
	return elapsed
}

func getCPUTime() syscall.Timeval {
	r := syscall.Rusage{}
	syscall.Getrusage(syscall.RUSAGE_SELF, &r)
	return r.Utime
}
