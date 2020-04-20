// +build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.melnyk.org/metricer"
	"go.melnyk.org/mlog"
	"go.melnyk.org/mlog/console"
)

func main() {
	logbook := console.NewLogbook(os.Stdout)
	logbook.SetLevel(mlog.Default, mlog.Verbose)

	metrics := metricer.NewHost(&metricer.Config{EnableDebug: true}, logbook)

	logger := logbook.Joiner().Join("example")

	metrics.NewHealthCheck("dummy", "Just dummy health check", func() error {
		// Add here any health check conditions
		logger.Verbose("Dummy health check called")
		return nil
	})

	counter := metrics.NewCounter("counter", "Basic demo counter")
	gauge := metrics.NewGauge("gause", "Basic demo gauge")

	metrics.Start()

	stop := make(chan interface{})
	go func() {
		ticker := time.NewTicker(time.Second)
	forever:
		for {
			select {
			case <-ticker.C:
				counter.Inc(rand.Int63n(20))
				gauge.Update(rand.Int63n(100))
			case <-stop:
				break forever
			}
		}
		ticker.Stop()
	}()

	fmt.Println("Press Ctrl-C...")
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
forever:
	for {
		select {
		case <-sig:
			logger.Warning("Stopping...")
			break forever
		}
	}

	close(stop)
	metrics.Stop()
}
