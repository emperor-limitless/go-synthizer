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
	"time"
)

func main() {
	fmt.Println("Enter file path")
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
	ctx, err := synthizer.NewContext()
	synthizer.GOCHECK(err)
	defer ctx.Destroy()
	ctx.Default_panner_strategy.Set(synthizer.PANNER_STRATEGY_HRTF)
	buf, err := synthizer.BufferFromFile(flname)
	synthizer.GOCHECK(err)
	defer buf.Destroy()
	gen, err := synthizer.NewBufferGenerator(ctx)
	synthizer.GOCHECK(err)
	defer gen.Destroy()
	gen.Buffer.Set(buf)
	gen.Looping.Set(true)
	src, err := synthizer.NewAngularPannedSource(ctx, synthizer.PANNER_STRATEGY_HRTF)
	synthizer.GOCHECK(err)
	defer src.Destroy()
	src.AddGenerator(gen)
	for angle := 0; angle < 360 * 10; angle ++ {
		src.Azimuth.Set(float32(angle % 360))
		time.Sleep(5 * time.Millisecond)
	}
}