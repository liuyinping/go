// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/e_log.c
// and came with this notice.  The go code is a simpler
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_log(x)
// Return the logrithm of x
//
// Method :
//   1. Argument Reduction: find k and f such that
//			x = 2^k * (1+f),
//	   where  sqrt(2)/2 < 1+f < sqrt(2) .
//
//   2. Approximation of log(1+f).
//	Let s = f/(2+f) ; based on log(1+f) = log(1+s) - log(1-s)
//		 = 2s + 2/3 s**3 + 2/5 s**5 + .....,
//	     	 = 2s + s*R
//      We use a special Reme algorithm on [0,0.1716] to generate
// 	a polynomial of degree 14 to approximate R The maximum error
//	of this polynomial approximation is bounded by 2**-58.45. In
//	other words,
//		        2      4      6      8      10      12      14
//	    R(z) ~ Lg1*s +Lg2*s +Lg3*s +Lg4*s +Lg5*s  +Lg6*s  +Lg7*s
//  	(the values of Lg1 to Lg7 are listed in the program)
//	and
//	    |      2          14          |     -58.45
//	    | Lg1*s +...+Lg7*s    -  R(z) | <= 2
//	    |                             |
//	Note that 2s = f - s*f = f - hfsq + s*hfsq, where hfsq = f*f/2.
//	In order to guarantee error in log below 1ulp, we compute log
//	by
//		log(1+f) = f - s*(f - R)	(if f is not too large)
//		log(1+f) = f - (hfsq - s*(hfsq+R)).	(better accuracy)
//
//	3. Finally,  log(x) = k*ln2 + log(1+f).
//			    = k*ln2_hi+(f-(hfsq-(s*(hfsq+R)+k*ln2_lo)))
//	   Here ln2 is split into two floating point number:
//			ln2_hi + ln2_lo,
//	   where n*ln2_hi is always exact for |n| < 2000.
//
// Special cases:
//	log(x) is NaN with signal if x < 0 (including -INF) ;
//	log(+INF) is +INF; log(0) is -INF with signal;
//	log(NaN) is that NaN with no signal.
//
// Accuracy:
//	according to an error analysis, the error is always less than
//	1 ulp (unit in the last place).
//
// Constants:
// The hexadecimal values are the intended ones for the following
// constants. The decimal values may be used, provided that the
// compiler will convert from decimal to binary accurately enough
// to produce the hexadecimal values shown.

const (
	Ln2Hi = 6.93147180369123816490e-01;	/* 3fe62e42 fee00000 */
	Ln2Lo = 1.90821492927058770002e-10;	/* 3dea39ef 35793c76 */
	Lg1 = 6.666666666666735130e-01;  /* 3FE55555 55555593 */
	Lg2 = 3.999999999940941908e-01;  /* 3FD99999 9997FA04 */
	Lg3 = 2.857142874366239149e-01;  /* 3FD24924 94229359 */
	Lg4 = 2.222219843214978396e-01;  /* 3FCC71C5 1D8E78AF */
	Lg5 = 1.818357216161805012e-01;  /* 3FC74664 96CB03DE */
	Lg6 = 1.531383769920937332e-01;  /* 3FC39A09 D078C69F */
	Lg7 = 1.479819860511658591e-01;  /* 3FC2F112 DF3E5244 */

	Two54 = 1<<54;				// 2^54
	TwoM20 = 1.0/(1<<20);		// 2^-20
	TwoM1022 = 2.2250738585072014e-308;	// 2^-1022
	Sqrt2 = 1.41421356237309504880168872420969808;
)

export func Log(x float64) float64 {
	// special cases
	switch {
	case sys.isNaN(x) || sys.isInf(x, 1):
		return x;
	case x < 0:
		return sys.NaN();
	case x == 0:
		return sys.Inf(-1);
	}

	// reduce
	f1, ki := sys.frexp(x);
	if f1 < Sqrt2/2 {
		f1 *= 2;
		ki--;
	}
	f := f1 - 1;
	k := float64(ki);

	// compute
	s := f/(2+f);
	s2 := s*s;
	s4 := s2*s2;
	t1 := s2*(Lg1 + s4*(Lg3 + s4*(Lg5 + s4*Lg7)));
	t2 := s4*(Lg2 + s4*(Lg4 + s4*Lg6));
	R :=  t1 + t2;
	hfsq := 0.5*f*f;
	return k*Ln2Hi - ((hfsq-(s*(hfsq+R)+k*Ln2Lo)) - f);
}

const
(
	ln10u1	= .4342944819032518276511;
)

export func Log10(arg float64) float64 {
	if arg <= 0 {
		return sys.NaN();
	}
	return Log(arg) * ln10u1;
}


