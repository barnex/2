from mumax2 import *
from mumax2_cmp import *

regionNameDictionary ={}

setgridsize(80, 40, 2)
setcellsize(5e-9, 5e-9, 50e-9)

load('micromagnetism')

setscalar('alpha', 0.1)
getscalar('alpha')
setscalar('Msat', 8.0e5)
setscalar('Aex', 12e-13)

m=[ [[[1]]], [[[0]]], [[[0]]] ]
setarray('m', m)

## Example of initialization of region system with a picture
imageName = os.path.abspath( os.path.dirname(__file__)) + '/engine_cmp_regions.png'
regionDic = {"M":"Blue",
			 "u":"Lime",
			 "m":"White",
			 "a":"DarkBlue",
			 "x":"Red",
			 "2":"Yellow"}
extrudeImage( imageName, regionDic)
MsatValues = {"M":1.e6,
			  "u":2.e6,
			  "m":3.e6,
			  "a":4.e6,
			  "x":5.e6,
			  "2":6.e6}
InitUniformRegionScalarQuant('Msat', MsatValues)

mValues = {"M":[ 1.0, 0.0,0.0],
		   "u":[ 1.0, 1.0,0.0],
		   "m":[ 0.0, 1.0,0.0],
		   "a":[-1.0, 1.0,0.0],
		   "x":[-1.0, 0.0,0.0],
		   "2":[-1.0,-1.0,0.0]}
InitUniformRegionVectorQuant('m', mValues)

save("Msat", "ovf", ["Text"], "regions_picture_Ms.ovf" )
save("m", "ovf", ["Text"], "regions_picture_m.ovf" )

## Example of initialization of region system with a script
def script(X, Y, Z, param):
	if Y > X * param["slope"]:
		return 'Upper'
	else:
		return 'Lower'
gridSize = getgridsize()
parameters = {"slope" : gridSize[1]/gridSize[0]}
#initRegionsScript( script , parameters)
mValues = {"Upper":1e6,
		   "Lower":2e6}
#InitUniformRegionScalarQuant('Msat', mValues)
#save("Msat", "ovf", ["Text"], "regions_script.ovf" )

Hx = 0 / mu0
Hy = 0 / mu0
Hz = 0.1 / mu0 

setvalue('H_ext', [Hx, Hy, Hz])

setscalar('dt', 1e-12)
#save("regionDefinition", "omf", ["Text"], "region.omf" )
#autosave("m", "omf", ["Text"], 10e-12)
#autosave("m", "ovf", ["Text"], 10e-12)
#autosave("m", "bin", [], 10e-12)
#autotabulate(["t", "H_ext"], "t.txt", 10e-12)
for i in range(100):
	step()

printstats()
savegraph("graph.png")





