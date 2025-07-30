package util

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
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
	cpu   *os.File
	block *os.File
	trace *os.File
	start time.Time
	end   time.Time
}

func (p *PProfProfiler) Init(name string) error {
	runtime.SetBlockProfileRate(1)
	n := fmt.Sprintf("%s.cpu.pprof", name)
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	p.cpu = f

	bn := fmt.Sprintf("%s.block.prof", name)
	bf, err := os.Create(bn)
	if err != nil {
		return err
	}
	p.block = bf

	tn := fmt.Sprintf("%s.trace.out", name)
	tf, err := os.Create(tn)
	if err != nil {
		return err
	}
	p.trace = tf

	return nil
}

func (p *PProfProfiler) Start() error {
	p.start = time.Now()
	if err := trace.Start(p.trace); err != nil {
		return err
	}
	return pprof.StartCPUProfile(p.cpu)
}

func (p *PProfProfiler) Stop() error {
	err := pprof.Lookup("block").WriteTo(p.block, 0)
	if err != nil {
		return err
	}
	p.end = time.Now()
	trace.Stop()
	pprof.StopCPUProfile()
	err = p.cpu.Close()
	if err != nil {
		return err
	}
	err = p.trace.Close()
	if err != nil {
		return err
	}

	return p.block.Close()
}

func (p *PProfProfiler) Duration() time.Duration {
	return p.end.Sub(p.start)
}
