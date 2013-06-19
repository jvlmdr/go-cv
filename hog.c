// AUTORIGHTS
// -------------------------------------------------------
// Copyright (C) 2011-2012 Ross Girshick
// Copyright (C) 2008, 2009, 2010 Pedro Felzenszwalb, Ross Girshick
// Copyright (C) 2007 Pedro Felzenszwalb, Deva Ramanan
//
// This file is part of the voc-releaseX code
// (http://people.cs.uchicago.edu/~rbg/latent/)
// and is available under the terms of an MIT-like license
// provided in COPYING. Please retain this notice and
// COPYING if you use this file (or a portion of it) in
// your project.
// -------------------------------------------------------

#include <stdio.h>
#include <stdlib.h>
#include <math.h>

// small value, used to avoid division by zero
#define eps 0.0001

// unit vectors used to compute gradient orientation
double uu[9] = {1.0000,
    0.9397,
    0.7660,
    0.500,
    0.1736,
    -0.1736,
    -0.5000,
    -0.7660,
    -0.9397};
double vv[9] = {0.0000,
    0.3420,
    0.6428,
    0.8660,
    0.9848,
    0.9848,
    0.8660,
    0.6428,
    0.3420};

double minDouble(double x, double y) { return (x <= y ? x : y); }
double maxDouble(double x, double y) { return (x <= y ? y : x); }

int minInt(int x, int y) { return (x <= y ? x : y); }
int maxInt(int x, int y) { return (x <= y ? y : x); }

void size(int* dims, int sbin, int* cells, int* out) {
  // memory for caching orientation histograms & their norms
  cells[0] = (int)round((double)dims[0]/(double)sbin);
  cells[1] = (int)round((double)dims[1]/(double)sbin);

  // memory for HOG features
  out[0] = maxInt(cells[0]-2, 0);
  out[1] = maxInt(cells[1]-2, 0);
  out[2] = 27+4;
}

// main function:
// takes a double color image and a bin size
// returns HOG features
void process(int* dims,
             double* im,
             int sbin,
             int* cells,
             int* out,
             double* feat) {
  int x;
  int y;
  int o;

  double *hist = (double *)malloc(cells[0]*cells[1]*18*sizeof(double));
  double *norm = (double *)malloc(cells[0]*cells[1]*sizeof(double));

  int visible[2];
  visible[0] = cells[0]*sbin;
  visible[1] = cells[1]*sbin;

	// "basis vectors" for addressing image pixels
	const int i_im = dims[2] * dims[0]; // dims[0];
	const int j_im = dims[2];           // 1;
	const int k_im = 1;                 // dims[0] * dims[1];
	const int i_hist = cells[0];
	const int j_hist = 1;
	const int k_hist = cells[0] * cells[1];
	const int i_norm = cells[0];
	const int j_norm = 1;
	const int i_feat = out[2] * out[0]; // out[0];
	const int j_feat = out[2];          // 1;
	const int k_feat = 1;               // out[0] * out[1];

  for (x = 1; x < visible[1]-1; x++) {
    for (y = 1; y < visible[0]-1; y++) {
			int a = minInt(x, dims[1]-2);
			int b = minInt(y, dims[0]-2);

      // first color channel
      double *s = im + a*i_im + b*j_im + 0*k_im;
      double dy = *(s+j_im) - *(s-j_im);
      double dx = *(s+i_im) - *(s-i_im);
      double v = dx*dx + dy*dy;

      // second color channel
      s += k_im;
      double dy2 = *(s+j_im) - *(s-j_im);
      double dx2 = *(s+i_im) - *(s-i_im);
      double v2 = dx2*dx2 + dy2*dy2;

      // third color channel
      s += k_im;
      double dy3 = *(s+j_im) - *(s-j_im);
      double dx3 = *(s+i_im) - *(s-i_im);
      double v3 = dx3*dx3 + dy3*dy3;

      // pick channel with strongest gradient
      if (v2 > v) {
        v = v2;
        dx = dx2;
        dy = dy2;
      }
      if (v3 > v) {
        v = v3;
        dx = dx3;
        dy = dy3;
      }

      // snap to one of 18 orientations
      double best_dot = 0;
      int best_o = 0;
      for (o = 0; o < 9; o++) {
        double dot = uu[o]*dx + vv[o]*dy;
        if (dot > best_dot) {
          best_dot = dot;
          best_o = o;
        } else if (-dot > best_dot) {
          best_dot = -dot;
          best_o = o+9;
        }
      }

      // add to 4 histograms around pixel using bilinear interpolation
      double xp = ((double)x+0.5)/(double)sbin - 0.5;
      double yp = ((double)y+0.5)/(double)sbin - 0.5;
      int ixp = (int)floor(xp);
      int iyp = (int)floor(yp);
      double vx0 = xp-ixp;
      double vy0 = yp-iyp;
      double vx1 = 1.0-vx0;
      double vy1 = 1.0-vy0;
      v = sqrt(v);

      if (ixp >= 0 && iyp >= 0) {
        *(hist + ixp*i_hist + iyp*j_hist + best_o*k_hist) += vx1*vy1*v;
      }

      if (ixp+1 < cells[1] && iyp >= 0) {
        *(hist + (ixp+1)*i_hist + iyp*j_hist + best_o*k_hist) += vx0*vy1*v;
      }

      if (ixp >= 0 && iyp+1 < cells[0]) {
        *(hist + ixp*i_hist + (iyp+1)*j_hist + best_o*k_hist) += vx1*vy0*v;
      }

      if (ixp+1 < cells[1] && iyp+1 < cells[0]) {
        *(hist + (ixp+1)*i_hist + (iyp+1)*j_hist + best_o*k_hist) += vx0*vy0*v;
      }
    }
  }

  // compute energy in each block by summing over orientations
  for (o = 0; o < 9; o++) {
    double *src1 = hist + o*k_hist;
    double *src2 = hist + (o+9)*k_hist;
    double *dst = norm;
    double *end = norm + cells[1]*cells[0];
    while (dst < end) {
      *(dst++) += (*src1 + *src2) * (*src1 + *src2);
      src1++;
      src2++;
    }
  }

  // compute features
  for (x = 0; x < out[1]; x++) {
    for (y = 0; y < out[0]; y++) {
      double *dst = feat + x*i_feat + y*j_feat;
      double *src, *p, n1, n2, n3, n4;

      p = norm + (x+1)*i_norm + (y+1)*j_norm;
      n1 = 1.0 / sqrt(*p + *(p+j_norm) + *(p+i_norm) + *(p+i_norm+j_norm) + eps);
      p = norm + (x+1)*i_norm + y*j_norm;
      n2 = 1.0 / sqrt(*p + *(p+j_norm) + *(p+i_norm) + *(p+i_norm+j_norm) + eps);
      p = norm + x*i_norm + (y+1)*j_norm;
      n3 = 1.0 / sqrt(*p + *(p+j_norm) + *(p+i_norm) + *(p+i_norm+j_norm) + eps);
      p = norm + x*i_norm + y*j_norm;
      n4 = 1.0 / sqrt(*p + *(p+j_norm) + *(p+i_norm) + *(p+i_norm+j_norm) + eps);

      double t1 = 0;
      double t2 = 0;
      double t3 = 0;
      double t4 = 0;

      // contrast-sensitive features
      src = hist + (x+1)*i_hist + (y+1)*j_hist;
      for (o = 0; o < 18; o++) {
        double h1 = minDouble(*src * n1, 0.2);
        double h2 = minDouble(*src * n2, 0.2);
        double h3 = minDouble(*src * n3, 0.2);
        double h4 = minDouble(*src * n4, 0.2);
        *dst = 0.5 * (h1 + h2 + h3 + h4);
        t1 += h1;
        t2 += h2;
        t3 += h3;
        t4 += h4;
        dst += k_feat;
        src += k_hist;
      }

      // contrast-insensitive features
      src = hist + (x+1)*i_hist + (y+1)*j_hist;
      for (o = 0; o < 9; o++) {
        double sum = *src + *(src + 9*k_hist);
        double h1 = minDouble(sum * n1, 0.2);
        double h2 = minDouble(sum * n2, 0.2);
        double h3 = minDouble(sum * n3, 0.2);
        double h4 = minDouble(sum * n4, 0.2);
        *dst = 0.5 * (h1 + h2 + h3 + h4);
        dst += k_feat;
        src += k_hist;
      }

      // texture features
      *dst = 0.2357 * t1;
      dst += k_feat;
      *dst = 0.2357 * t2;
      dst += k_feat;
      *dst = 0.2357 * t3;
      dst += k_feat;
      *dst = 0.2357 * t4;
    }
  }

  free(hist);
  free(norm);
}
