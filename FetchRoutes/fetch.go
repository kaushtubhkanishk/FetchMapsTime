package FetchRoutes

import (
	"context"
	"crypto/tls"
	"math"
	"os"
	"time"

	"github.com/rs/zerolog"
	routespb "google.golang.org/genproto/googleapis/maps/routing/v2"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	fieldMask  = "*"
	serverAddr = "routes.googleapis.com:443"
)

func Fetch(log zerolog.Logger) int64 {
	apiKey := os.Getenv("MAPS_API_KEY")
	config := tls.Config{}
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(&config)))
	if err != nil {
		log.Fatal().Msgf("did not connect: %v", err)
	}
	defer conn.Close()
	client := routespb.NewRoutesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-Goog-Api-Key", apiKey)
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-Fieldmask", fieldMask)
	defer cancel()

	origin := &routespb.Waypoint{
		LocationType: &routespb.Waypoint_Location{
			Location: &routespb.Location{
				LatLng: &latlng.LatLng{
					Latitude:  28.496639249663314,
					Longitude: 77.09455535620648,
				},
			},
		},
	}

	destination := &routespb.Waypoint{
		LocationType: &routespb.Waypoint_Location{
			Location: &routespb.Location{
				LatLng: &latlng.LatLng{
					Latitude:  28.58145626804548,
					Longitude: 77.18588956970005,
				},
			},
		},
	}

	req := &routespb.ComputeRoutesRequest{
		Origin:                   origin,
		Destination:              destination,
		TravelMode:               routespb.RouteTravelMode_DRIVE,
		RoutingPreference:        routespb.RoutingPreference_TRAFFIC_AWARE,
		ComputeAlternativeRoutes: true,
		Units:                    routespb.Units_METRIC, RouteModifiers: &routespb.RouteModifiers{
			AvoidTolls: false, AvoidHighways: false, AvoidFerries: true,
		},
		PolylineQuality: routespb.PolylineQuality_OVERVIEW}
	resp, err := client.ComputeRoutes(ctx, req)

	if err != nil {
		log.Fatal().Err(err)
	}

	log.Debug().Msgf("No. of routes: %d", len(resp.Routes))
	var minima int64
	minima = math.MaxInt64
	for _, route := range resp.Routes {
		if route.Duration.Seconds < minima {
			minima = route.Duration.Seconds
		}
	}
	log.Debug().Msgf("Duration in seconds: %d", minima)
	log.Info().Msgf("Duration in minutes: %d:%d", minima/60, minima%60)
	return minima
}
