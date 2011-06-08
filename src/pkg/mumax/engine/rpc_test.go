//  This file is part of MuMax, a high-perfomrance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package engine

// Author: Arne Vansteenkiste

import (
	//. "mumax/common"
	"testing"
)

type Tester struct {
	Value int
}

func (t *Tester) SetValue(i int) {
	t.Value = i
}

func (t *Tester) GetValue() int {
	return t.Value
}

func TestRPC(t *testing.T) {
	//	server := NewRPC()
	//	client := NewRPC()
	//
	//	server.Register(new(Tester))
	//
	//	end1, end2 := Pipe2()
	//	go server.ServeConn(end1)
	//	go client.ServeConn(end2)
	//
	//	client.Call("SetValue", []interface{}{42})
	//	ret := client.Call("GetValue", []interface{}{})[0].(int)
	//
	//	if ret != 42 {
	//		t.Fail()
	//	}
}