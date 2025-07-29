package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/jcocozza/gmenu/internal/app"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func usage() {
	fmt.Printf("usage: %s [options] [FILE]\n", os.Args[0])
	fmt.Println("")
	fmt.Println("read from stdin and display the results in a selectable list")
	fmt.Println("")
	fmt.Println("options:")
	flag.PrintDefaults()
}

func printVersion() {
	_, err := fmt.Fprintf(os.Stdout, "%s: built on %s, commit: %s\n", Version, Date, Commit)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	// this *NEEDS* to be in the main method, otherwise things break
	runtime.LockOSThread()
	flag.Usage = usage

	var version bool
	flag.BoolVar(&version, "v", false, "print version and exit")
	flag.BoolVar(&version, "version", false, "print version and exit")

	var alias bool
	flag.BoolVar(&alias, "a", false, "enable alias parsing mode")
	flag.BoolVar(&alias, "alias", false, "enable alias parsing mode")

	flag.Parse()

	if version {
		printVersion()
		os.Exit(0)
	}

	var in io.Reader
	if flag.NArg() > 0 {
		fpath := flag.Arg(0)
		file, err := os.Open(fpath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "error: could not open file: %v", err)
			os.Exit(1)
		}
		defer file.Close()
		in = file
	} else {
		in = os.Stdin
	}

	cfg := app.AppConfig{Alias: alias}

	a, err := app.NewGMenuApp(cfg, in)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: startup: %v", err)
		os.Exit(1)
	}
	a.Render()
}
