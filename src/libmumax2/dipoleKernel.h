/**
  * @file
  *
  * @author Ben Van de Wiele
  */

#ifndef _DIPOLEKERNEL_H_
#define _DIPOLEKERNEL_H_

#include <cuda.h>

#ifdef __cplusplus
extern "C" {
#endif

/// Initialization of the elements of the dipole kernel
/// @param data: float arrray to store the kernel data
/// @param co1: Component of the field (0, 1, or 2)
/// @param co2: Component of the source dipole (0, 1, or 2)
/// @param N0: x-size of the data array
/// @param N1: y-size of the data array
/// @param N2: z-size of the data array
/// @param N1part: y-size of the data array, stored on 1 device
/// @param per0: periodic repetitions in x-direction
/// @param per1: periodic repetitions in y-direction
/// @param per2: periodic repetitions in z-direction
/// @param cellX: size of the FD cell in x-direction
/// @param cellY: size of the FD cell in y-direction
/// @param cellZ: size of the FD cell in z-direction
/// @param dev_qd_P_10: float array containing the Gauss quadrature points for integration (10th order)
/// @param dev_qd_W_10: float array containing the Gauss quadrature weight for integration (10th order)
/// @param streams: used streams
void initFaceKernel6ElementAsync(float **data, 
                                 int co1,
                                 int co2,
                                 int N0, int N1, int N2,          
                                 int N1part,                                /// size of the kernel
                                 int per0, int per1, int per2,              /// periodicity
                                 float cellX, float cellY, float cellZ,     /// cell size
                                 float **dev_qd_P_10, float **dev_qd_W_10,  /// quadrature points and weights
                                 CUstream *streams
                                );

#ifdef __cplusplus
}
#endif
#endif