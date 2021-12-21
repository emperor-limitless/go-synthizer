# go-synthizer
A golang binding to the [synthizer](https://github.com/synthizer/synthizer) audio library.
## Status
This binding is not production ready yet, And there's a big possibility major buggs can be found, Since I'm new to go, So if you found any, Please submit an issue or a pull request if you would like to fix it you're self.
## What can we do now?
 - Initialize, With or without libraryConfig.
 - Creating context.
 - Creating generators, StreamingGenerator and bufferGenerator.
 - Creating Sources, Direct, ScalarPanned, AngularPanned, And Source3D
 - Creating Buffers.
 - Creating streams, From file and memory, Custom is not implemented because I'm not really sure of a safe way to do it, And to be honest, I don't understand how it works my self, So if you do, Please submit a pull request.
 - Editing properties.