//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package engine

import ()

// Loads a test module.
func (e *Engine) LoadTest() {
	e.AddQuant("m", VECTOR, FIELD)

	e.AddQuant("h_z", VECTOR, FIELD)
	e.Depends("h_z", "t")
	e.AddQuant("h", VECTOR, FIELD)
	e.Depends("h", "h_z")

	e.AddQuant("torque", VECTOR, FIELD)
	e.Depends("torque", "m")
	e.Depends("torque", "h")

	e.ODE1("m", "torque")
}
