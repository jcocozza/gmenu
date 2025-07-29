package util

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

type Profiler interface {
	// must be called first
	Init(path string) error

	// actually profiles the code
	Start() error
	Stop() error

	// analysis - only after Stop() has been called
	Duration() time.Duration
}

func NewProfiler() Profiler {
	return &PProfProfiler{}
}

type PProfProfiler struct {
	f     *os.File
	start time.Time
	end   time.Time
}

func (p *PProfProfiler) Init(name string) error {
	n := fmt.Sprintf("%s.cpu.pprof", name)
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	p.f = f
	return nil
}

func (p *PProfProfiler) Start() error {
	p.start = time.Now()
	return pprof.StartCPUProfile(p.f)
}

func (p *PProfProfiler) Stop() error {
	p.end = time.Now()
	pprof.StopCPUProfile()
	return p.f.Close()
}

func (p *PProfProfiler) Duration() time.Duration {
	return p.end.Sub(p.start)
}
