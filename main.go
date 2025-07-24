package main

import (
	"bufio"
	"os"

	"github.com/jcocozza/gmenu/menu"
	"github.com/jcocozza/gmenu/renderer"
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

func main() {
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
