# gmenu

gmenu is a cross platform version of [dmenu](https://tools.suckless.org/dmenu/).

It has been tested on ubuntu(i3), an arm mac and windows 11.

I haven't done a feature by feature copy over, but the essential feature, selecting from a GUI list and writing to stdout, remains the same.

I also added a niche, but perhaps useful "alias" mode that parses items just a little more intelligently.


## (Known) Limitations

- Doesn't behave quite as expected when using i3, but still works


## Developement

For profiling build the binary with `profile` tag: `go build -tags profile ./cmd/gmenu/`.

