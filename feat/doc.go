/*
Package feat gives an interface for computing feature channels from images.

The two main interfaces are feat.Transform and feat.Real.
The former describes a mapping from an image.Image to real values and the latter describes a mapping from real values to real values.
Both provide an integer downsample rate using Rate().

Real transforms can be chained togather using feat.Compose.

The global factory system is messy.
In order to be deserializable using the factory, a transform should define a Marshaler() method to satisfy the Marshalable interface and be registered using RegisterXxx().
The Register method associates the name to a function which creates a Spec.
Most transforms will be able to use a Spec created by NewXxxSpec, but compound transforms such as Compose must define their own Spec in order to deserialize their abstract (interface) members.
The Marshaler types deserialize a transform using the registered names to obtain a Spec.
*/
package feat
