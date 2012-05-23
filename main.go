// Copyright 2012 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	var path string
	if u, err := user.Current(); err != nil {
		log.Fatalln("Could not determine home directory")
	} else {
		path = filepath.Join(u.HomeDir, ".crawl", "morgue")
	}

	var invertGrouping bool
	flag.BoolVar(&invertGrouping, "i", false, "invert grouping (display names by type)")
	flag.Parse()

	filenames, err := filepath.Glob(filepath.Join(path, "morgue-*.txt"))
	if err != nil {
		log.Fatalln(err)
	}
	var (
		filter   bool
		restrict map[string]bool
	)
	if args := flag.Args(); len(args) > 0 {
		filter = true
		restrict = make(map[string]bool, len(args))
		for _, name := range args {
			restrict[strings.ToLower(name)] = true
		}
	}
	groups := make(map[string][]string, len(filenames))
	for _, filename := range filenames {
		name := ExtractName(filename)
		if filter && !restrict[strings.ToLower(name)] {
			continue
		}
		desc, err := Parse(filename, name)
		if err != nil {
			continue
		}
		var key, value string
		if invertGrouping {
			key, value = desc, name
		} else {
			key, value = name, desc
		}
		groups[key] = SetAdd(groups[key], value)
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		set := groups[key]
		sort.Strings(set)
		fmt.Printf("%s:\n", key)
		for _, value := range set {
			fmt.Println("-", value)
		}
		fmt.Println()
	}
}

func SetAdd(s []string, v string) []string {
	for _, x := range s {
		if v == x {
			return s
		}
	}
	return append(s, v)
}

func ExtractName(filename string) string {
	parts := strings.Split(filename, "-")
	return strings.Join(parts[1:len(parts)-2], "-")
}
