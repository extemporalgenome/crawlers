// Copyright 2012 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"io"
	"os"
)

func SkipLine(r io.ByteReader) error {
	for {
		if b, err := r.ReadByte(); err != nil || b == '\n' {
			return err
		}
	}
	return nil
}

// Assumes that even partial matches may not overlap. Stops at match or newline.
func ScanUntil(r io.ByteReader, target string, ignorespace bool) (scanned string, pos int, err error) {
	var i, j, k int
	var b byte
	buf := make([]byte, 0, len(target))
	pos = -1
	for {
		b, err = r.ReadByte()
		if err != nil {
			return
		} else if b == '\n' {
			return
		}
		switch target[j] {
		case b:
			buf = append(buf, b)
			j++
		default:
			if j > 0 && (b == ' ' || b == '\t') {
				buf = append(buf, b)
				k++
			} else {
				j = 0
			}
		}
		if j >= len(target) {
			scanned = string(buf)
			pos = i - j - k + 1
			return
		}
		i++
	}
	return
}

// Assumes that even partial matches may not overlap. Stops at newline.
func ScanLine(r io.ByteReader, target string, ignorespace bool) (scanned string, pos int, err error) {
	scanned, pos, err = ScanUntil(r, target, ignorespace)
	if pos >= 0 {
		err = SkipLine(r)
	}
	return
}

func Parse(path, name string) (parsed_name, desc string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	r := bufio.NewReader(file)
	var i int
	switch _, i, err = ScanLine(r, "character file.", false); {
	case i < 0:
		err = errors.New("Unrecognized morgue format")
		fallthrough
	case err != nil:
		return
	}
	for {
		if parsed_name, i, err = ScanUntil(r, name, true); i == 0 {
			if _, i, err = ScanUntil(r, "(", false); i >= 0 {
				var s string
				if s, err = r.ReadString(')'); err == nil {
					desc = s[:len(s)-1]
					return
				}
			}
		}
		if i > 0 {
			err = SkipLine(r)
		}
		if err != nil {
			return
		}
	}
	return
}
