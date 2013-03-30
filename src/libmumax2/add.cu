
#include "add.h"

#include "multigpu.h"
#include <cuda.h>
#include "gpu_conf.h"
#include "gpu_safe.h"

#ifdef __cplusplus
extern "C" {
#endif

///@internal
__global__ void addKern(float* dst, float* a, float* b, int Npart) {
	int i = threadindex;
	if (i < Npart) {
		dst[i] = a[i] + b[i];
	}
}


__export__ void addAsync(float** dst, float** a, float** b, CUstream* stream, int Npart) {
	dim3 gridSize, blockSize;
	make1dconf(Npart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		addKern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], b[dev], Npart);
	}
}

///@internal
__global__ void addmaddKern(float* dst, float* a, float* b, float* c, float mul, int Npart) {
	int i = threadindex;
	if (i < Npart) {
		dst[i] = a[i] + mul * (b[i] + c[i]);
	}
}

///@internal
__global__ void maddKern(float* dst, float* a, float* b, float mulB, int Npart) {
	int i = threadindex;
	float bMask;
	if (b == NULL){
		bMask = 1.0f;
	}else{
		bMask = b[i];
	}
	if (i < Npart) {
		dst[i] = a[i] + mulB * bMask;
	}
}


__export__ void addMaddAsync(float** dst, float** a, float** b, float** c, float mul, CUstream* stream, int NPart)
{
	dim3 gridSize, blockSize;
	make1dconf(NPart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		addmaddKern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], b[dev], c[dev], mul, NPart);
	}
	
}
__export__ void maddAsync(float** dst, float** a, float** b, float mulB, CUstream* stream, int Npart) {
	dim3 gridSize, blockSize;
	make1dconf(Npart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		maddKern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], b[dev], mulB, Npart);
	}
}

///@internal
__global__ void madd1Kern(float* a, float* b, float mulB, int Npart) {
	int i = threadindex;
	float bMask;
	if (b == NULL){
		bMask = 1.0f;
	}else{
		bMask = b[i];
	}
	if (i < Npart) {
		a[i] += mulB * bMask;
	}
}


__export__ void madd1Async(float** a, float** b, float mulB, CUstream* stream, int Npart) {
	dim3 gridSize, blockSize;
	make1dconf(Npart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		madd1Kern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (a[dev], b[dev], mulB, Npart);
	}
}

///@internal
__global__ void madd2Kern(float* a, float* b, float mulB, float* c, float mulC, int Npart) {
	int i = threadindex;

	float bMask;
	if (b == NULL){
		bMask = 1.0f;
	}else{
		bMask = b[i];
	}

	float cMask;
	if (c == NULL){
		cMask = 1.0f;
	}else{
		cMask = c[i];
	}

	if (i < Npart) {
		a[i] += mulB * bMask + mulC * cMask;
	}
}


__export__ void madd2Async(float** a, float** b, float mulB, float** c, float mulC, CUstream* stream, int Npart) {
	dim3 gridSize, blockSize;
	make1dconf(Npart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		madd2Kern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (a[dev], b[dev], mulB, c[dev], mulC, Npart);
	}
}


__global__ void cmaddKern(float* dst, float a, float b, float* kern, float* src, int NComplexPart){

  int i = threadindex; // complex index
  int e = 2 * i; // real index

  if(i < NComplexPart){

    float Sa = src[e  ];
    float Sb = src[e+1];

	float k = kern[i];

	dst[e  ] += k * (a*Sa - b*Sb);
	dst[e+1] += k * (b*Sa + a*Sb);
  }
  
  return;
}

__export__ void cmaddAsync(float** dst, float a, float b, float** kern, float** src, CUstream* stream, int NComplexPart){
	dim3 gridSize, blockSize;
	make1dconf(NComplexPart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		cmaddKern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a, b, kern[dev], src[dev], NComplexPart);
	}
}

///@internal
__global__ void linearCombination2Kern(float* dst, float* a, float mulA, float* b, float mulB, int NPart) {
	int i = threadindex;
	if (i < NPart) {
		dst[i] = mulA * a[i] + mulB * b[i];
	}
}

__export__ void linearCombination2Async(float** dst, float** a, float mulA, float** b, float mulB, CUstream* stream, int NPart){
	dim3 gridSize, blockSize;
	make1dconf(NPart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		linearCombination2Kern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], mulA, b[dev], mulB, NPart);
	}
}

///@internal
__global__ void linearCombination3Kern(float* dst, float* a, float mulA, float* b, float mulB, float* c, float mulC, int NPart) {
	int i = threadindex;
	if (i < NPart) {
		dst[i] = mulA * a[i] + mulB * b[i] + mulC * c[i];
	}
}

__export__ void linearCombination3Async(float** dst, float** a, float mulA, float** b, float mulB, float** c, float mulC, CUstream* stream, int NPart){
	dim3 gridSize, blockSize;
	make1dconf(NPart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		linearCombination3Kern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], mulA, b[dev], mulB, c[dev], mulC, NPart);
	}
}

///@internal
__global__ void linearCombination4Kern(float* dst, float* a, float mulA, float* b, float mulB, float* c, float mulC, float* d, float mulD, int NPart) {
	int i = threadindex;
	if (i < NPart) {
		dst[i] = mulA * a[i] + mulB * b[i] + mulC * c[i] + mulD * d[i];
	}
}

__export__ void linearCombination4Async(float** dst, float** a, float mulA, float** b, float mulB, float** c, float mulC, float** d, float mulD, CUstream* stream, int NPart){
	dim3 gridSize, blockSize;
	make1dconf(NPart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		linearCombination4Kern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], mulA, b[dev], mulB, c[dev], mulC, d[dev], mulD, NPart);
	}
}

///@internal
__global__ void linearCombination5Kern(float* dst, float* a, float mulA, float* b, float mulB, float* c, float mulC, float* d, float mulD, float* e, float mulE, int NPart) {
	int i = threadindex;
	if (i < NPart) {
		dst[i] = mulA * a[i] + mulB * b[i] + mulC * c[i] + mulD * d[i] + mulE * e[i];
	}
}

__export__ void linearCombination5Async(float** dst, float** a, float mulA, float** b, float mulB, float** c, float mulC, float** d, float mulD, float** e, float mulE, CUstream* stream, int NPart){
	dim3 gridSize, blockSize;
	make1dconf(NPart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		linearCombination5Kern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], mulA, b[dev], mulB, c[dev], mulC, d[dev], mulD, e[dev], mulE, NPart);
	}
}

///@internal
__global__ void linearCombination6Kern(float* dst, float* a, float mulA, float* b, float mulB, float* c, float mulC, float* d, float mulD, float* e, float mulE, float* f, float mulF, int NPart) {
	int i = threadindex;
	if (i < NPart) {
		dst[i] = mulA * a[i] + mulB * b[i] + mulC * c[i] + mulD * d[i] + mulE * e[i] + mulF * f[i];
	}
}

__export__ void linearCombination6Async(float** dst, float** a, float mulA, float** b, float mulB, float** c, float mulC, float** d, float mulD, float** e, float mulE, float** f, float mulF, CUstream* stream, int NPart){
	dim3 gridSize, blockSize;
	make1dconf(NPart, &gridSize, &blockSize);
	for (int dev = 0; dev < nDevice(); dev++) {
		gpu_safe(cudaSetDevice(deviceId(dev)));
		linearCombination6Kern <<<gridSize, blockSize, 0, cudaStream_t(stream[dev])>>> (dst[dev], a[dev], mulA, b[dev], mulB, c[dev], mulC, d[dev], mulD, e[dev], mulE, f[dev], mulF, NPart);
	}
}

#ifdef __cplusplus
}
#endif
