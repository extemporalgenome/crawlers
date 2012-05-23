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
func ScanUntil(r io.ByteReader, target string) (int, error) {
	i, j := 0, 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			return -1, err
		} else if b == '\n' {
			return -1, nil
		}
		switch target[j] {
		case b:
			j++
		default:
			j = 0
		}
		if j >= len(target) {
			return i - j + 1, nil
		}
		i++
	}
	return -1, nil
}

// Assumes that even partial matches may not overlap. Stops at newline.
func ScanLine(r io.ByteReader, target string) (int, error) {
	i, err := ScanUntil(r, target)
	if i >= 0 {
		err = SkipLine(r)
		return i, err
	}
	return -1, err
}

func Parse(path, name string) (desc string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	r := bufio.NewReader(file)
	var i int
	switch i, err = ScanLine(r, "character file."); {
	case i < 0:
		err = errors.New("Unrecognized morgue format")
		fallthrough
	case err != nil:
		return
	}
	for {
		if i, err = ScanUntil(r, name); i == 0 {
			if i, err = ScanUntil(r, "("); i >= 0 {
				if s, err := r.ReadString(')'); err == nil {
					return s[:len(s)-1], nil
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
