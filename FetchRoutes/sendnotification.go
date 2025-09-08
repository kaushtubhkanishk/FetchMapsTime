package FetchRoutes

import (
	"net/http"
	"net/url"
	"os"

	"github.com/rs/zerolog"
)

var (
	token = os.Getenv("APP_TOKEN")
	user  = os.Getenv("USER_TOKEN")
)

// SendNotification is a function that sends a notification with a given string
func SendNotification(log zerolog.Logger, minDuration string) {
	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", url.Values{
		"token":   {token},
		"user":    {user},
		"message": {minDuration},
	})
	defer resp.Body.Close()
	if err != nil {
		log.Fatal().Msg("Failed to send notification with : " + err.Error())
	}
	log.Info().Msg("Successfully sent notification with : " + resp.Status)
}
