package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jcocozza/gmenu/menu"
	"github.com/jcocozza/gmenu/renderer"
)

// these are set on build
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func readStdin() ([]menu.Item, error) {
	var results []menu.Item
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		results = append(results, menu.BasicItem(scanner.Text()))
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
	m := menu.NewGmenu(elms)
	r := renderer.RendererFactory()
	if err := r.Init(); err != nil {
		panic(err)
	}
	defer r.Cleanup()

	m.ExecSearch()
	// init render
	if err := r.RenderFrame(m); err != nil {
		panic(err)
	}
	for !r.Done() {
		glfw.PollEvents()
		act := r.Action()
		//  we don't need to rerender if no action has happened
		if act == nil {
			continue
		}
		switch act.(type) {
		case menu.ActionSelect:
			r.ClearAction()
			item := m.SelectedItem()
			fmt.Println(item.Value())
		default:
			act.Apply(m)
			r.ClearAction()
		}
		if err := r.RenderFrame(m); err != nil {
			panic(err)
		}
	}
}
