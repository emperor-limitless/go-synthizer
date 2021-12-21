//MIT License
//
//Copyright (c) 2021 mohamedSulaimanAlmarzooqi
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.
//
package main

import (
	"github.com/mohamedSulaimanAlmarzooqi/go-synthizer"
	"fmt"
	"io/ioutil"
)

func main() {
	fmt.Println("Enter file path, Example: hello.mp3.")
	var flname string
	fmt.Scanln(&flname)
	if flname == "" {
		fmt.Println("Error: Invalid path.")
		return
	}
	conf := synthizer.NewLibraryConfig(synthizer.LOG_LEVEL_DEBUG, synthizer.LOGGING_BACKEND_STDERR)
	synthizer.InitializeWithConfig(&conf)
	defer synthizer.Shutdown()
	err, ctx := synthizer.NewContext()
	defer ctx.Destroy()
	synthizer.GOCHECK(err)
	dt, err := ioutil.ReadFile(flname)
	synthizer.GOCHECK(err)
	err, sh := synthizer.StreamHandleFromMemory(string(dt))
	synthizer.GOCHECK(err)
	err, gen := synthizer.StreamingGeneratorFromHandle(ctx, sh)
	defer gen.Destroy()
	synthizer.GOCHECK(err)
	err = gen.Looping.Set(false)
	synthizer.GOCHECK(err)
	err, src := synthizer.NewDirectSource(ctx)
	defer src.Destroy()
	synthizer.GOCHECK(err)
	src.AddGenerator(gen.Generator)
	src.Gain.Set(0.6)
	src.Play()
	fmt.Println("Press enter to exit...")
	var empty string
	fmt.Scanln(&empty)
	sh = nil
}