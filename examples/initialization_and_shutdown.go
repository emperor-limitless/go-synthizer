package main

import (
	"github.com/mohamedSulaimanAlmarzooqi/go-synthizer"
	"time"
)

func main() {
	conf := synthizer.NewLibraryConfig(synthizer.LOG_LEVEL_DEBUG, synthizer.LOGGING_BACKEND_STDERR)
	synthizer.InitializeWithConfig(&conf)
	time.Sleep(180 * time.Second)
	synthizer.Shutdown()
}