package main

import (
	"context"
	"flag"
	"runtime"
	"time"

	"github.com/bborbe/argument"
	"github.com/golang/glog"

	"github.com/bborbe/imap-spam-deleter/pkg"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = flag.Set("logtostderr", "true")

	time.Local = time.UTC
	glog.V(2).Infof("set global timezone to UTC")

	var app pkg.Application
	if err := argument.Parse(&app); err != nil {
		glog.Fatalf("parse args failed: %v", err)
	}
	if err := app.Run(context.Background()); err != nil {
		glog.Exitf("failed: %v", err)
	}
	glog.Infof("completed")
}
