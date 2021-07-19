package poller

import (
	"log"
	"time"

	"github.com/memfault/envaws/param_providers"
)

func Poll(provider param_providers.ParamProvider, pollingIntervalSeconds int64, done chan bool) {
	provider.Init()
	log.Print("Polling every ", pollingIntervalSeconds, " seconds")
	tick := time.Tick(time.Duration(pollingIntervalSeconds * 1e9))
	for range tick {
		log.Println("Polling for config change...")
		if provider.Changed() {
			done <- true
		}
	}
}
