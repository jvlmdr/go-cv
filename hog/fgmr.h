#ifndef HOG_H__
#define HOG_H__

void size(int* dims, int sbin, int* cells, int* out);

void compute(int* dims,
             double* im,
             double* hist,
             double* norm,
             int sbin,
             int* cells,
             int* out,
             double* feat);

#endif
