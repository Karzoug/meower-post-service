package server

import (
	"encoding/json"
	"net/http"

	"github.com/Karzoug/meower-post-service/pkg/buildinfo"
	"github.com/rs/zerolog"
)

func Buildinfo(logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		if err := json.NewEncoder(w).Encode(buildinfo.Get()); err != nil {
			logger.Error().
				Err(err).
				Msg("failed to encode buildinfo")
		}
	}
}
