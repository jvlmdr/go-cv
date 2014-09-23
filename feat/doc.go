/*
Package feat gives an interface for computing feature channels from images.

The two main interfaces are feat.Transform and feat.Real.
The former describes a mapping from an image.Image to real values and the latter describes a mapping from real values to real values.
Both provide an integer downsample rate using Rate().

Real transforms can be chained togather using feat.Compose.
*/
package feat
