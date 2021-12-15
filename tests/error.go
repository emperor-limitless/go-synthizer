package main


import (
	"fmt"
	"github.com/mohamedSulaimanAlmarzooqi/go-synthizer"
	"github.com/mohamedSulaimanAlmarzooqi/go-synthizer/synthizer"
)

func main() {
	synthizer.GOCHECK(synthizer.Initialize())
	fmt.Println("Started the first time here.")
	defer synthizer.GOCHECK(synthizer.Shutdown())
	fmt.Println("Going to start the second time, Should error at this point.")
	synthizer.GOCHECK(synthizer.Initialize()) // Should error here.
}