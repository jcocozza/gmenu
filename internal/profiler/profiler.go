//go:build profile

package profiler


import (
	"os"
	"fmt"
	"github.com/jcocozza/goprof"
)

func StartProfiler() {
	if err := goprof.Start("gmenu"); err != nil {
		fmt.Fprintf(os.Stdout, "error: startup: %v", err)
		os.Exit(1)
	}
}

func StopProfiler() {
	if err := goprof.Stop(); err != nil {
		fmt.Fprintf(os.Stdout, "error: stopping: %v", err)
		os.Exit(1)
	}
	fmt.Println("init render time")
	goprof.Summarize()
}
