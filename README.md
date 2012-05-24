# crawlers

A utility for [Dungeon Crawl Stone Soup](http://crawl.develz.org) that displays
backgrounds grouped by character name.

## Installing
crawlers is written in [Go](http://golang.org), and requires compilation.
Running `go get github.com/extemporalgenome/crawlers` on a system with a Go 1
installation should produce a `crawlers` binary, probably in `$GOROOT/bin`

## Running crawlers

`crawlers` without arguments lists each background combination by
character name. Passing the `-i` option will invert grouping--instead
displaying character names by associated background.

Non-option arguments are taken to be character names to limit output to:
`crawlers Charlie Joe` will only display the backgrounds of characters that
were named either Charlie or Joe.

While output is primarily intended for human consumption, it should also be
valid YAML.

## Platform/Version Restrictions

crawlers is only confirmed to work with the Linux port of crawl. It
specifically looks in `~/.crawl/morgue`. Any crawl distribution that stores
morgue character files in that path ought to work. Character files produced
between crawl versions 0.7 and the current 0.10.x are confirmed to work.
