package main

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/kaushtubhkanishk/FetchMapsTime/FetchRoutes"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	var globalMin int64
	globalMin = math.MaxInt64
	for {
		currentMin := FetchRoutes.Fetch(log)
		log.Debug().Int64("min", currentMin).Msg("FetchRoutes.Fetch")
		if currentMin < globalMin {
			globalMin = currentMin
			FetchRoutes.SendNotification(log, strconv.FormatInt(globalMin/60, 10)+" : "+strconv.FormatInt(globalMin%60, 10))
			log.Debug().Int64("min", globalMin).Msg("FetchRoutes.SendNotification")
		}
		time.Sleep(30 * time.Second)
	}
}
