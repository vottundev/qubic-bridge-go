package interceptor

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vottundev/vottun-qubic-bridge-go/utils"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

type key int

const requestIDKey key = 0

var (
	defaultRequestID string = ""
)

// Intercepts the request and calculates the total run time from start to finish
func NewElapsedTimeInterceptor() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, err := utils.GenerateRandomString(25)
			if err != nil {
				requestID = &defaultRequestID
			}

			ctx := context.WithValue(r.Context(), requestIDKey, requestID)

			startTime := float64(time.Now().UnixNano()) / float64(time.Millisecond)

			if r == nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Errorf("Missing Request")
				return
			}

			log.Infof("**** %s - New Request Arrived: Request info: [%s %s%s]", *requestID, r.Method, r.Host, r.URL.Path)

			defer func() {
				endTime := float64(time.Now().UnixNano()) / float64(time.Millisecond)
				elapsed := float64((endTime - startTime) / 1000)
				log.Infof("**** %s - Time consumed for query to %s is %.2f seconds", *requestID, r.URL.Path, math.Round(elapsed*100)/100)
			}()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return "UNKNOWN"
}
