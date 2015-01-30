/*
Package slide contains sliding-window operations.

Function names are matched by this regular expression:
	(Cos)?(Conv|Corr)(Multi)?(Bank)?(Naive|FFT|BLAS)
"Multi" means that the input image has multiple channels.
"Bank" means that there is a bank of filters and therefore the output image has multiple channels.
Not all combinations are implemented.
*/
package slide
