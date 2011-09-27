from mumax2 import *

setgridsize(512, 512, 2)
print 'gridsize', getgridsize()

setcellsize(5e-9, 5e-9, 50e-9)
print 'cellsize', getcellsize()

load('test')
savegraph("graph.dot")

setscalar('alpha', 0.05)
getscalar('alpha')
setscalar('alpha', 0.05)

#print 'alpha', getvalue('alpha'), '\n'
#print 'alphaMask', getmask('alpha'), '\n'
#print 'alpha', getfield('alpha'), '\n'

setscalar('Msat', 800e3)
print 'Msat', getvalue('Msat'), '\n'

m=[ [[[0.01]]], [[[0]]], [[[-1]]] ]
setfield('m', m)

i=3
j=2
k=1
setcell('m', i,j,k, [0,1,0])
print 'getcell', getcell('m', i,j,k)


Hx = 0 / mu0
Hy = 0 / mu0
Hz = 0.1 / mu0 #1T

setvalue('H_z', [Hx, Hy, Hz])
mask = [ [ [[0]],[[0]] ], [ [[0]], [[0]] ], [ [[1]], [[0]] ] ]
#setmask('H_z', mask)
#print 'H_z',getvalue('H_z'), '\n'
#print 'H', getfield('H'), '\n'
#setfield('torque', m) # must fail

setscalar('dt', 2e-12)
f = open('ll', 'w')
for i in range(1000):
	t = getscalar('t')
	m = getcell('m', 0,0,0)
	f.write(str(t) + "\t")
	f.write(str(m[0]) + "\t")
	f.write(str(m[1]) + "\t")
	f.write(str(m[2]) + "\n")
	step()

printstats()

f.close()





