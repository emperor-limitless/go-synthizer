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
	"strings"
	"strconv"
	"bufio"
	"os"
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
	src, err := synthizer.NewSource3D(ctx)
	synthizer.GOCHECK(err)
	defer src.Destroy()
	src.AddGenerator(gen)
	var argv []string
	for true {
		fmt.Println("Enter command...")
		sr := bufio.NewScanner(os.Stdin)
		sr.Scan()
		argv = strings.Split(sr.Text(), " ")
		if argv[0] == "pos" {
			if len(argv) < 4 {
				fmt.Println("Invalid number of arguments.")
				fmt.Println("Usage: pos <x> <y> <z>")
				continue
			}
			x, err := strconv.ParseFloat(argv[1], 32)
			if err != nil{
				fmt.Println(err)
				continue
			}
			y, err := strconv.ParseFloat(argv[2], 32)
			if err != nil{
				fmt.Println(err)
				continue
			}
			z, err := strconv.ParseFloat(argv[3], 32)
			if err != nil{
				fmt.Println(err)
				continue
			}
			err = src.Position.Set(float32(x), float32(y), float32(z))
			if err != nil{
				fmt.Println(err)
				continue
			}
		} else if argv[0] == "pause" {
			gen.Pause()
		} else if argv[0] == "play" {
			gen.Play()
		} else if argv[0] == "gain" {
			val, err := strconv.ParseFloat(argv[1], 32)
			if err != nil {
				fmt.Println("Error: The argument is invalid, Must be a floating point number.")
				continue
			}
			if val > 0.9 {
				fmt.Println("Error, Value out of range.")
				continue
			}
			gen.Gain.Set(float32(val))
		} else if argv[0] == "seek" {
			if len(argv) != 2 {
				fmt.Println("Syntax: seek <seconds>")
				continue
			}
			val, err := strconv.ParseFloat(argv[1], 32)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = gen.PlaybackPosition.Set(float32(val))
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else if argv[0] == "quit" {
			break
		} else {
			fmt.Println("Unknown command.")
			continue
		}
	}
}