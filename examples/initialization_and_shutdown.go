package main

import "github.com/mohamedSulaimanAlmarzooqi/go-synthizer"

func main() {
	conf := synthizer.NewLibraryConfig(synthizer.LOG_LEVEL_DEBUG, synthizer.LOGGING_BACKEND_STDERR)
	synthizer.InitializeWithConfig(&conf)
	synthizer.Shutdown()
}