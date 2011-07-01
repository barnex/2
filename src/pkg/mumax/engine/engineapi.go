//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package engine


import ()


type EngineAPI struct {
	engine *Engine
}


// For testing purposes.
func (api EngineAPI) Version() int {
	return 2
}

// For testing purposes.
func (api EngineAPI) Echo(i int) int {
	return i
}

// For testing purposes.
func (api EngineAPI) Sink(b bool, i int, f float32, d float64, s string) {
	return
}

// For testing purposes.
func (api EngineAPI) GetFloat() float32 {
	return 42.
}

// For testing purposes.
func (api EngineAPI) GetDouble() float64 {
	return 42.
}

// For testing purposes.
func (api EngineAPI) GetString() string {
	return "hello"
}


// For testing purposes.
func (api EngineAPI) GetBool() bool {
	return true
}


// For testing purposes.
func (api EngineAPI) Sum(i, j int) int {
	return i + j
}