//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package engine

// Author: Arne Vansteenkiste

import ()

// Register this module
func init() {
	RegisterModule(&ModZeeman{})
}

type ModZeeman struct{}

func (x ModZeeman) Description() string {
	return "externally applied field (A/m)"
}

func (x ModZeeman) Name() string {
	return "zeeman"
}

func (x ModZeeman) Load(e *Engine) {
	e.LoadModule("hfield")
	e.AddQuant("H_ext", VECTOR, MASK, Unit("A/m"), "ext. field")
	hfield := e.Quant("H")
	sum := hfield.updater.(*SumUpdater)
	sum.AddParent("H_ext")
	e.Depends("H_ext", "t") // EVEN IF H_ext IS NOT REALLY TIME-DEPENDENT, THINGS BREAK IF THIS IS NOT HERE. NEED TO DEBUG
}
