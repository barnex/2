//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

// This package implements automated mumax API generation.
// Based on the exported methods of engine.API, an API
// library in any of the supported programming languages is
// automatically generated.
//
// Author: Arne Vansteenkiste
package apigen

import (
	."mumax/common"
	"io/ioutil"
)

// Auto-generate API libraries for all languages.
func APIGen2() {
	_, err := ioutil.ReadFile(GetExecDir()+"../src/pkg/mumax/engine/api.go")
	CheckIO(err)
}
