//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package modules

// Author: Arne Vansteenkiste

import (
	. "mumax/engine"
)

// Register this module.
// Module for the total field H. Other modules like H_demag, H_anis, H_ext have
// to add their field to the sum make by H.
// H is in units A/m and does not have a multiplier. I.e., is not normalized to Msat.
func init() {
	RegisterModule("hfield", "Total magnetic field.", LoadHField)
}

func LoadHField(e *Engine) {
	e.AddNewQuant("H", VECTOR, FIELD, Unit("A/m"), "magnetic field")
	q := e.Quant("H")
	q.SetUpdater(NewSumUpdater(q))
}
