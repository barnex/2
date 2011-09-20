//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package engine

//
//	Here the user input (X,Y,Z) is changed to internal input (Z,Y,X)
//

import (
	. "mumax/common"
	"mumax/host"
	"fmt"
	"os"
)

// The API methods are accessible to the end-user through scripting languages.
type API struct {
	Engine *Engine
}

//________________________________________________________________________________ init

// Set the grid size.
// WARNING: convert to ZYX
func (a API) SetGridSize(x, y, z int) {
	a.Engine.SetGridSize([]int{z, y, x}) // convert to internal axes
}

// Get the grid size.
// WARNING: convert to ZYX
func (a API) GetGridSize() (x, y, z int) {
	size := a.Engine.GridSize()
	return size[Z], size[Y], size[X] // convert to internal axes
}

// Set the cell size.
// WARNING: convert to ZYX
func (a API) SetCellSize(x, y, z float64) {
	a.Engine.SetCellSize([]float64{z, y, x}) // convert to internal axes and units
}

// Set the cell size.
// WARNING: convert to ZYX, internal units
func (a API) GetCellSize() (x, y, z float64) {
	size := a.Engine.CellSize()
	return size[Z], size[Y], size[X] // convert to internal axes
}

// Load a physics module. Not aware of dependencies (yet)
// TODO: cleaner management a la modprobe
func (a API) Load(module string) {
	switch module {
	default:
		panic(InputErr(fmt.Sprint("Unknown module:", module, " Options: micromag")))
	case "test":
		a.Engine.LoadTest()
	case "micromag":
		a.Engine.LoadMicromag()
	case "micromagenergy":
		a.Engine.LoadMicromagEnergy()
	case "spintorque":
		a.Engine.LoadSpintorque()
	}
}

//________________________________________________________________________________ run

// Take one solver step
func (a API) Step() {
	a.Engine.Step()
}

//________________________________________________________________________________ set quantities

// Sets the value of a quantity. The quantity must be of type VALUE or MASK.
// If the quantity is a MASK, the value will be multiplied by a space-dependent mask
// which typically contains dimensionless numbers between 0 and 1.
func (a API) SetValue(name string, value []float64) {
	q := a.Engine.Quant(name)
	swapXYZ(value)
	//Debug("swapXYZ", value)
	q.SetValue(value)
}

// Sets a space-dependent multiplier mask for the quantity.
// The value of the quantity (set by SetValue), will be multiplied
// by the mask value in each point of space. The mask is dimensionless
// and typically contains values between 0 and 1.
func (a API) SetMask(name string, mask *host.Array) {
	q := a.Engine.Quant(name)
	qArray := q.Array()
	if !EqualSize(mask.Size3D, qArray.Size3D()) {
		Log("auto-resampling ", q.Name(), "from", Size(mask.Size3D), "to", Size(qArray.Size3D()))
		mask = Resample(mask, qArray.Size3D())
	}
	q.SetMask(mask)
}


// Sets a space-dependent field quantity, like the magnetization.
func (a API) SetField(quant string, field *host.Array) {
	q := a.Engine.Quant(quant)
	qArray := q.Array()
	if !EqualSize(field.Size3D, qArray.Size3D()) {
		Log("auto-resampling ", quant, "from", Size(field.Size3D), "to", Size(qArray.Size3D()))
		field = Resample(field, qArray.Size3D())
	}
	q.SetField(field)
}


//________________________________________________________________________________ get quantities



// Get the value of a space-independent or masked quantity.
// Returns an array with vector components or an
// array with just one element in case of a scalar quantity.
func (a API) GetValue(name string) []float64 {
	q := a.Engine.Quant(name)
	q.Update()//!
	value := make([]float64, len(q.multiplier))
	copy(value, q.multiplier)
	swapXYZ(value)
	return value
}

// Get the value of a scalar, space-independent quantity.
// Similar to GetValue, but returns a single number.
func (a API) GetScalar(name string) float64 {
	q := a.Engine.Quant(name)
	q.Update()//!
	return q.Scalar()
}


// Gets a space-dependent quantity. If the quantity uses a mask,
// value*mask is returned.
func (a API) GetField(quant string) *host.Array {
	q := a.Engine.Quant(quant)
	checkKinds(q, MASK, FIELD)
	q.Update()//!
	array := q.Array()
	buffer := q.Buffer()
	if array.IsNil(){
		set ones	
	}else{
	array.CopyToHost(buffer)
	}
	mul with mul
	return buffer
}

//________________________________________________________________________________ internal

// INTERNAL: swaps the X-Z values of the array.
func swapXYZ(array []float64) {
	if len(array) == 3 {
		array[X], array[Z] = array[Z], array[X]
	}
	return
}

// Set the value of a scalar, space-independent quantity
//func (a API) SetScalar(name string, value float32) {
//	e := a.Engine
//	q := e.Quant(name)
//	q.SetValue([]float32{value})
//}




// Get the value of a general quantity
//func (a API) Get(name string) interface{} {
//	e := a.Engine
//	q := e.Quant(name)
//	switch {
//	case q.IsScalar():
//		return q.ScalarValue()
//	}
//	panic(Bug("unimplemented case"))
//}

//________________________________________________________________________________ misc

// Save .dot file with the physics graph for plotting with graphviz.
func (a API) SaveGraph(file string) {
	f, err := os.Create(file)
	defer f.Close()
	CheckIO(err)
	a.Engine.WriteDot(f)
	Log("Wrote", file, "Run command: \"dot -Tpng", file, "> myfile.png\" to plot the physics graph (requires package graphviz).")
}
