//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package modules

import (
	. "mumax/common"
	. "mumax/engine"
	"mumax/gpu"
)

// Register this module
func init() {
	RegisterModule("3T", "Three-temperature model module", Load3T)
}

func Load3T(e *Engine) {
    //Load Temperatures
    LoadT(e)
    LoadGes(e)
    Te := e.Quant("Te")
    Ts := e.Quant("Ts")
    Tl := e.Quant("Tl")
    
    Q := e.AddNewQuant("Q", SCALAR, MASK, Unit("J/s"), "The power of external heater")
    gamma_e := e.AddNewQuant("gamma_e", SCALAR, MASK, Unit("J/(K^2*m^3)"), "The heat capacity of electrons")
    Cs := e.AddNewQuant("Cs", SCALAR, MASK, Unit("J/(K*m^3)"), "The heat capacity of spins")
    Cl := e.AddNewQuant("Cl", SCALAR, MASK, Unit("J/(K*m^3)"), "The heat capacity of phonons")
    
    Gel := e.AddNewQuant("Gel", SCALAR, MASK, Unit("W/(K*m^3)"), "The electron-phonon coupling coefficient")
    //Ges := e.AddNewQuant("Ges", SCALAR, MASK, Unit("W/(K*m^3)"), "The electron-spin coupling coefficient")
    Ges := e.Quant("Ges")
    Gsl := e.AddNewQuant("Gsl", SCALAR, MASK, Unit("W/(K*m^3)"), "The spin-phonon coupling coefficient")
    
    Qe := e.AddNewQuant("Qe", SCALAR, FIELD, Unit("W/(m^3)"), "The heat flux density of electrons")
    Qs := e.AddNewQuant("Qs", SCALAR, FIELD, Unit("W/(m^3)"), "The heat flux density of spins")
    Ql := e.AddNewQuant("Ql", SCALAR, FIELD, Unit("W/(m^3)"), "The heat flux density of phonons")
    
    e.Depends("Qe", "Te", "Tl", "Ts", "gamma_e", "Ges", "Gel", "Q")
    e.Depends("Qs", "Te", "Tl", "Ts", "Cs",      "Gsl", "Ges")
    e.Depends("Ql", "Te", "Tl", "Ts", "Cl",      "Gsl", "Gel")
    
    Qe.SetUpdater(&QeUpdater{Qe: Qe, Te: Te, Tl: Tl, Ts: Ts, Q: Q, gamma_e: gamma_e, Gel: Gel, Ges: Ges })
    Qs.SetUpdater(&QsUpdater{Qs: Qs, Te: Te, Tl: Tl, Ts: Ts, Cs: Cs, Gsl: Gsl, Ges: Ges })
    Ql.SetUpdater(&QlUpdater{Ql: Ql, Te: Te, Tl: Tl, Ts: Ts, Cl: Cl, Gel: Gel, Gsl: Gsl })
    
    e.AddPDE1("Te", "Qe")
    e.AddPDE1("Ts", "Qs")
    e.AddPDE1("Tl", "Ql")
    
}

type QeUpdater struct {
	Qe, Te, Tl, Ts, Q, gamma_e, Gel, Ges *Quant
}

type QsUpdater struct {
	Qs, Te, Tl, Ts, Cs, Gsl, Ges *Quant
}

type QlUpdater struct {
	Ql, Te, Tl, Ts, Cl, Gel, Gsl *Quant
}

func (u *QeUpdater) Update() {
    gpu.Qe_async(
        u.Qe.Array(),
        u.Te.Array(),
        u.Tl.Array(),
        u.Ts.Array(),
        u.Q.Array(),
        u.gamma_e.Array(),
        u.Gel.Array(),
        u.Ges.Array(),
        u.Q.Multiplier(),
        u.gamma_e.Multiplier(),
        u.Gel.Multiplier(),
        u.Ges.Multiplier())
    u.Qe.Array().Sync()
}

func (u *QsUpdater) Update() {
    gpu.Qs_async(
        u.Qs.Array(),
        u.Te.Array(),
        u.Tl.Array(),
        u.Ts.Array(),
        u.Cs.Array(),
        u.Gsl.Array(),
        u.Ges.Array(),
        u.Cs.Multiplier(),
        u.Gsl.Multiplier(),
        u.Ges.Multiplier())
    u.Qs.Array().Sync()
}

func (u *QlUpdater) Update() {
    gpu.Ql_async(
        u.Ql.Array(),
        u.Te.Array(),
        u.Tl.Array(),
        u.Ts.Array(),
        u.Cl.Array(),
        u.Gel.Array(),
        u.Gsl.Array(),
        u.Cl.Multiplier(),
        u.Gel.Multiplier(),
        u.Gsl.Multiplier())
    u.Ql.Array().Sync()
}

