package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/jcocozza/gmenu/menu"
	"github.com/jcocozza/gmenu/renderer"
)

// these are set on build
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func readStdin() ([]string, error) {
	var results []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func usage() {
	fmt.Printf("usage: %s [options]\n", os.Args[0])
	fmt.Println("")
	fmt.Println("read from stdin and display the results in a selectable list")
	fmt.Println("")
	fmt.Println("options:")
	fmt.Println("-v    show version")
	fmt.Println("-h    show this help message")

}

func main() {
	flag.Usage = usage
	version := flag.Bool("v", false, "print version and exit")
	flag.Parse()

	if *version {
		_, err := fmt.Fprintf(os.Stdout, "%s: built on %s, commit: %s\n", Version, Date, Commit)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	elms, err := readStdin()
	if err != nil {
		panic(err)
	}
	m := menu.MenuFactory(elms)
	r := renderer.RendererFactory(m)
	if err := r.Init(); err != nil {
		panic(err)
	}
	defer r.Cleanup()
	err = r.Render()
	if err != nil {
		panic(err)
	}
}
