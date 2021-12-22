# go-synthizer
A golang binding to the [synthizer](https://github.com/synthizer/synthizer) audio library.
## Status
While many functionality is missing, But this should work if you try to use it for basic audio usage, Stuff like custom streams, FromFloatArrays, Events and filters are not here yet, Although should be soon.
## What can we do now?
 - Initialize, With or without libraryConfig.
 - Creating context.
 - Creating generators, StreamingGenerator, NoiseGenerator and bufferGenerator.
 - Creating Sources, Direct, ScalarPanned, AngularPanned, And Source3D
 - Creating Buffers.
 - Creating streams, From file and memory, Custom is not implemented because I'm not really sure of a safe way to do it, And to be honest, I don't understand how it works my self, So if you do, Please submit a pull request.
 - Editing properties.
 - Creating reverb with GlobalFdnReverb.