// healthcheck is a package that provides a standart middleware to
// healthcheck http server and its dependencies or sub-component,
// inspired mostly by the upcoming IETF RFC Health Check Response Format for HTTP APIs
// https://inadarei.github.io/rfc-healthcheck/
package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Karzoug/meower-post-service/pkg/ucerr"
	"github.com/rs/zerolog"
)

type HealthChecker interface {
	Healthcheck(ctx context.Context) (key string, err error)
}

type status string

const (
	// pass represents a healthy service "pass"
	pass status = "pass"
	// fail represents an unhealthy service "fail"
	fail status = "fail"
	// warn represents a healthy service with some minor problem "warn"
	warn status = "warn"
)

type responseComponentCheck struct {
	Status status `json:"status"`
	Output string `json:"output,omitempty"`
}

type response struct {
	Status status                            `json:"status"`
	Checks map[string]responseComponentCheck `json:"checks"`
}

func HealthCheck(logger zerolog.Logger, checkers ...HealthChecker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := response{
			Checks: make(map[string]responseComponentCheck),
		}

		var hasProblems bool
		for _, checker := range checkers {
			key, err := checker.Healthcheck(r.Context())
			if err != nil {
				var ch responseComponentCheck
				serr, ok := err.(ucerr.Error)
				if ok {
					ch = responseComponentCheck{
						Status: fail,
						Output: serr.Error(),
					}
				} else {
					ch = responseComponentCheck{
						Status: fail,
					}
				}
				resp.Checks[key] = ch

				hasProblems = true

				continue
			}
			resp.Checks[key] = responseComponentCheck{
				Status: pass,
			}
		}

		if hasProblems {
			w.WriteHeader(500)
			resp.Status = fail
		} else {
			w.WriteHeader(200)
			resp.Status = pass
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error().
				Err(err).
				Msg("failed to encode healthcheck response")
		}
	}
}
