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
	for {
		now := time.Now()
		next7Pm := time.Date(now.Year(), now.Month(), now.Day(), 19, 0, 0, 0, now.Location())
		if now.After(next7Pm) {
			next7Pm = next7Pm.Add(24 * time.Hour)
		}
		sleepDuration := next7Pm.Sub(now)
		log.Info().Dur("until_7pm", sleepDuration).Msg("Sleeping until 7pm")
		time.Sleep(sleepDuration)
		var globalMin int64
		globalMin = math.MaxInt64
		for i := 0; i < 90; i++ {
			currentMin := FetchRoutes.Fetch(log)
			log.Debug().Int64("min", currentMin).Msg("FetchRoutes.Fetch")
			if currentMin < globalMin {
				globalMin = currentMin
				FetchRoutes.SendNotification(log, strconv.FormatInt(globalMin/60, 10)+" : "+strconv.FormatInt(globalMin%60, 10))
				log.Debug().Int64("min", globalMin).Msg("FetchRoutes.SendNotification")
			}
			time.Sleep(2 * time.Minute)
		}
	}
}
