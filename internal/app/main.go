package app

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jcocozza/gmenu/internal/menu"
	"github.com/jcocozza/gmenu/internal/profiler"
	"github.com/jcocozza/gmenu/internal/renderer"
)

type AppConfig struct {
	Alias bool
}

// GMenuApp glues together all the components
type GMenuApp struct {
	cfg AppConfig
	m   *menu.GMenu
	r   renderer.Renderer
}

func NewGMenuApp(cfg AppConfig, r io.Reader) (*GMenuApp, error) {
	a := &GMenuApp{
		cfg: cfg,
	}
	items, err := a.readAndParseItems(r)
	if err != nil {
		return nil, err
	}
	m := menu.NewGmenu(items)
	render := renderer.RendererFactory()
	a.m = m
	a.r = render
	return a, nil
}

func (a *GMenuApp) readAndParseItems(r io.Reader) ([]menu.Item, error) {
	scanner := bufio.NewScanner(r)
	var items []menu.Item
	for scanner.Scan() {
		var itm menu.Item
		if a.cfg.Alias {
			itm = menu.ParseAliasItem(scanner.Text())
		} else {
			itm = menu.BasicItem(scanner.Text())
		}
		items = append(items, itm)
	}
	return items, scanner.Err()
}

func (a *GMenuApp) Render() {
	if err := a.r.Init(); err != nil {
		panic(err)
	}
	defer a.r.Cleanup()
	a.m.ExecSearch()
	// init render
	if err := a.r.InitalRender(a.m); err != nil {
		panic(err)
	}
	profiler.StopProfiler()
	for !a.r.Done() {
		a.r.PollEvents()
		act := a.r.Action()
		//  we don't need to rerender if no action has happened
		if act == nil {
			continue
		}
		switch act.(type) {
		case menu.ActionSelect:
			a.r.ClearAction()
			item := a.m.SelectedItem()
			if item != nil { // do nothing if we don't have any results
				fmt.Println(item.Value())
				a.r.MarkClose()
			}
		default:
			act.Apply(a.m)
			a.r.ClearAction()
		}
		if err := a.r.RenderFrame(a.m); err != nil {
			panic(err)
		}
	}
}
