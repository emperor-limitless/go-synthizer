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
	"bufio"
	"os"
)

func main() {
	fmt.Println("Enter file path, Example: hello.mp3.")
	var flname string
	sr := bufio.NewScanner(os.Stdin)
	sr.Scan()
	flname = sr.Text()
	if flname == "" {
		fmt.Println("Error: Invalid path.")
		return
	}
	conf := synthizer.NewLibraryConfig(synthizer.LOG_LEVEL_DEBUG, synthizer.LOGGING_BACKEND_STDERR)
	synthizer.InitializeWithConfig(&conf)
	defer synthizer.Shutdown()
	err, ctx := synthizer.NewContext()
	synthizer.GOCHECK(err)
	defer ctx.Destroy()
	err, buf := synthizer.BufferFromFile(flname)
	synthizer.GOCHECK(err)
	defer buf.Destroy()
	err, gen := synthizer.NewBufferGenerator(ctx)
	synthizer.GOCHECK(err)
	defer gen.Destroy()
	gen.Buffer.Set(buf.ObjectBase)
	err, src := synthizer.NewSource3D(ctx)
	synthizer.GOCHECK(err)
	defer src.Destroy()
	err, reverb := synthizer.NewGlobalFdnReverb(ctx)
	synthizer.GOCHECK(err)
	reverb.T60.Set(10.0)
	err = ctx.ConfigRoute(src.ObjectBase, reverb.ObjectBase)
	synthizer.GOCHECK(err)
	src.AddGenerator(gen.Generator)
	fmt.Println("Press enter to exit.")
	fmt.Scanln(&flname)
}