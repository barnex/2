//  This file is part of MuMax, a high-performance micromagnetic simulator.
//  Copyright 2011  Arne Vansteenkiste and Ben Van de Wiele.
//  Use of this source code is governed by the GNU General Public License version 3
//  (as published by the Free Software Foundation) that can be found in the license.txt file.
//  Note that you are welcome to modify this code under the condition that you do not remove any 
//  copyright notices and prominently state that you modified it, giving a relevant date.

package apigen

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

type Python struct{}

func (p *Python) Filename() string {
	return "mumax2.py"
}

func (p *Python) Comment() string {
	return "#"
}

func (p *Python) WriteHeader(out io.Writer) {
	fmt.Fprintln(out, `
import os
import json
import sys
import socket

infifo = 0
outfifo = 0
m_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
initialized = 0
outputdir = ""
M_HOST, M_PORT = "localhost", 3655

## Initializes the communication with mumax2.
# @note Internal use only
def init():
	global infifo
	global outfifo
	global outputdir
	global s_sock
	global m_sock
	global initialized
	global M_HOST
	global M_PORT
	
	# get the output directory from environment
	# outputdir=os.environ["MUMAX2_OUTPUTDIR"] + "/"	
	
	m_sock.connect((M_HOST,M_PORT))
	initialized = 1
	cmsg = 'python:'+ sys.argv[0] + '\n'
	print 'Sending back to master: ' + cmsg	
	m_sock.sendall(cmsg)
	s_addr = recvall2(m_sock)
	print 'MuMax grants connection on: ' + s_addr
	s_name = s_addr.split(':')
	S_HOST = s_name[0]
	S_PORT = int(s_name[1])
	s_sock.connect((S_HOST, S_PORT))	
	print 'python frontend is initialized'

## Calls a mumax2 command and returns the result as string.
# @note Internal use only.
def call(command, args):
	if (initialized == 0):
		init()
	s_sock.sendall(json.dumps([command, args])+'\n')
	# return json.loads(s_sock.recv(4096))
	return json.loads(recvall2(s_sock))

End='<<< End of mumax message >>>'

def recvall2(the_socket):
    total_data=[];data=''
    while True:
            data=the_socket.recv(8192)
            if End in data:
                total_data.append(data[:data.find(End)])
                break
            total_data.append(data)
            if len(total_data)>1:
                #check if end_of_data was split
                last_pair=total_data[-2]+total_data[-1]
                if End in last_pair:
                    total_data[-2]=last_pair[:last_pair.find(End)]
                    total_data.pop()
                    break
    return ''.join(total_data)
	
def recvall(the_socket,timeout=''):
    #setup to use non-blocking sockets
    #if no data arrives it assumes transaction is done
    #recv() returns a string
    the_socket.setblocking(0)
    total_data=[];data=''
    begin=time.time()
    if not timeout:
        timeout=1
    while 1:
        #if you got some data, then break after wait sec
        if total_data and time.time()-begin>timeout:
            break
        #if you got no data at all, wait a little longer
        elif time.time()-begin>timeout*2:
            break
        wait=0
        try:
            data=the_socket.recv(4096)
            if data:
                total_data.append(data)
                begin=time.time()
                data='';wait=0
            else:
                time.sleep(0.1)
        except:
            pass
        #When a recv returns 0 bytes, other side has closed
    result=''.join(total_data)
    return result
`)
}

func (p *Python) WriteFooter(out io.Writer) {

}

func (p *Python) WriteFunc(out io.Writer, name string, comment []string, argNames []string, argTypes []reflect.Type, returnTypes []reflect.Type) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintln(os.Stderr, "WriteFunc ", name, comment, argNames, err)
		}
	}()

	fmt.Fprintln(out)
	fmt.Fprintf(out, pyDocComment(comment))
	fmt.Fprint(out, "def ", name, "(")

	args := ""
	for i := range argTypes {
		if i != 0 {
			args += ", "
		}
		args += argNames[i]
	}
	fmt.Fprintln(out, args, "):")

	fmt.Fprintf(out, `	ret = call("%s", [%s])`, name, args)
	fmt.Fprint(out, "\n	return ")
	for i := range returnTypes {
		if i != 0 {
			fmt.Fprint(out, ", ")
		}
		fmt.Fprintf(out, `%v(ret[%v])`, python_convert[returnTypes[i].String()], i)
	}
	fmt.Fprintln(out)
	//fmt.Fprintln(out, fmt.Sprintf(`	return %s(call("%s", [%s])[0])`, python_convert[retType], name, args)) // single return value only
}

var (
	// maps go types to python types	
	python_convert map[string]string = map[string]string{"int": "int",
		"float32": "float",
		"float64": "float",
		"string":  "str",
		"bool":    "bool",
		"":        ""}
)

// Puts python doc comment tokens in front of the comment lines.
func pyDocComment(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	str := "#"
	for _, l := range lines {
		str += "# " + l + "\n"
	}
	return str
}
