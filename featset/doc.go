/*
Package featset provides a global factory for feature transforms.

An ImageMarshaler is a wrapper around an Image transform interface for the purpose of serialization.
It adds a string which identifies the type of the transform.
The UnmarshalJSON() method of ImageMarshaler will use the global factory to create a new object of the type specified in this string.
In order to be marshalable, an Image transform must provide Marshaler(), returning an ImageMarshaler containing itself and its type identifier string.
It must also define a Transform() method which usually simply returns itself.
In contrast, an ImageMarshaler is an Image transform whose Transform() method returns the transform which it wraps, and whose Marshaler() method returns itself.
This ensures that ImageMarshalers never wrap themselves.
Compound transforms, that is transforms which have another transform as one of their members, should invoke their members' Transform() methods within theirs.
*/
package featset
