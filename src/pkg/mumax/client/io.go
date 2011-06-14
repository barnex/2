//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package client

// This file implements some basic I/O functions.

import (
	. "mumax/common"
	"fmt"
	"io"
	"os"
	"strings"
)


// makes a fifo
// syscall.Mkfifo seems unavailable for the moment.
func mkFifo(fname string) {
	err := syscommand("mkfifo", []string{fname})
	if err != nil {
		panic(IOErr(fmt.Sprintf("mkfifo", fname, "returned", err)))
	}
}


// reads a line from the reader and splits it in words
func parseLine(in io.Reader) (words []string, eof bool) {
	str := ""
	var c byte
	c, eof = readChar(in)
	if eof {
		return
	}
	for c != '\n' {
		str += string(c)
		//Debug("str:", str)
		c, eof = readChar(in)
		if eof {
			return
		}
	}
	words = strings.Split(str, " ", -1)
	return
}

// reads one character from the reader
func readChar(in io.Reader) (char byte, eof bool) {
	var buffer [1]byte

	n := 0
	var err os.Error
	for n == 0 {
		n, err = in.Read(buffer[:])
		if err != nil {
			Debug(err)
			eof = true
			return
		}
	}
	char = buffer[0]
	return
}