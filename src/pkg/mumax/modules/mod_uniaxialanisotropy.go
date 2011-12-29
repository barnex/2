//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package modules

// This file implements the uniaxial anisotropy module
// Author: Arne Vansteenkiste

import (
	. "mumax/common"
	. "mumax/engine"
	"mumax/gpu"
)

// Register this module
func init() {
	RegisterModule("anisotropy/uniaxial", "Uniaxial magnetocrystalline anisotropy", LoadAnisUniaxial)
}

func LoadAnisUniaxial(e *Engine) {
	LoadHField(e)

	Hanis := e.AddNewQuant("H_anis", VECTOR, FIELD, Unit("A/m"), "uniaxial anisotropy field")
	ku1 := e.AddNewQuant("Ku1", SCALAR, MASK, Unit("J/m3"), "uniaxial anisotropy constant K1")
	ku2 := e.AddNewQuant("Ku2", SCALAR, MASK, Unit("J/m3"), "uniaxial anisotropy constant K2")
	anisU := e.AddNewQuant("anisU", VECTOR, MASK, Unit(""), "uniaxial anisotropy direction (unit vector)")

	hfield := e.Quant("H")
	sum := hfield.Updater().(*SumUpdater)
	sum.AddParent("H_anis")
	e.Depends("H_anis", "Ku1", "Ku2", "anisU", "MSat")

	Hanis.SetUpdater(&UniaxialAnisUpdater{e.Quant("m"), Hanis, ku1, ku2, anisU})
}

type UniaxialAnisUpdater struct {
	m, hanis, ku1, ku2, anisU *Quant
}

func (u *UniaxialAnisUpdater) Update() {
	hanis := u.hanis.Array()
	m := u.m.Array()
	ku1 := u.ku1.Array()
	ku1mul := u.ku1.Multiplier()[0]
	ku2 := u.ku2.Array()
	ku2mul := u.ku2.Multiplier()[0]
	anisU := u.anisU.Array()
	anisUMul := u.anisU.Multiplier()
	stream := u.hanis.Array().Stream

	// TODO
	msat := GetEngine().Quant("msat")
	gpu.UniaxialAnisotropyAsync(hanis, m, ku1, 2*ku1mul/(Mu0*msat.Scalar()), ku2, ku2mul, anisU, anisUMul, stream)

	u.hanis.Array().Stream.Sync()
}
