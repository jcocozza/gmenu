# gmenu

gmenu is a cross platform version of [dmenu](https://tools.suckless.org/dmenu/).

It has been tested on ubuntu(i3), an arm mac and windows 11.

I haven't done a feature by feature copy over, but the essential feature, selecting from a GUI list and writing to stdout, remains the same.

I also added a niche, but perhaps useful "alias" mode that parses items just a little more intelligently.

## Why?

[dmenu](https://tools.suckless.org/dmenu/) is very good, but unfortunately, some of the time I cannot be on linux.
I also don't like when I can't use the same tool across multiple platforms. I want something that functions the same regardless of where I use it.
Moreover, I found some of the existing alternatives to dmenu a little too clunky.

## Uses

1. Not using spotlight on Macos.
If you use something like skhd:
`cmd + shift - space : open -a "/Applications/$(ls /Applications | /path/to/gmenu)"`

In the future I may figure out how to package this up nicely so it's an actual app for macos.
