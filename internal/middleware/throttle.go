package middleware

import (
	"net/http"
	"time"

	"github.com/warrenb95/cloud-native-go/internal/model"
)

type Throttle struct {
	buckets map[string]*bucket
	max     uint
	refill  uint
	d       time.Duration
}

// bucket tracks the request for a given UID.
type bucket struct {
	tokens uint
	time   time.Time
}

func NewThrottle(max uint, refill uint, d time.Duration) *Throttle {
	return &Throttle{
		buckets: make(map[string]*bucket),
		max:     max,
		refill:  refill,
		d:       d,
	}
}

// Throttle will manage incoming request frequency for a given UID.
func (th *Throttle) Throttle(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if th.max < 1 {
			http.Error(w, model.ErrInternalError.Error(), http.StatusInternalServerError)
			return
		}

		UIDCookie, err := r.Cookie("UID")
		if err != nil {
			http.Error(w, model.ErrInvalidArgument.Error(), http.StatusBadRequest)
			return
		}

		b := th.buckets[UIDCookie.Value]

		if b == nil {
			th.buckets[UIDCookie.Value] = &bucket{tokens: th.max - 1, time: time.Now().UTC()}
			next.ServeHTTP(w, r)
			return
		}

		refillInterval := uint(time.Since(b.time) / th.d)
		tokensAdded := th.refill * refillInterval
		currentTokens := b.tokens + tokensAdded

		if currentTokens < 1 {
			http.Error(w, model.ErrTooManyRequests.Error(), http.StatusTooManyRequests)
			return
		}

		if currentTokens > th.max {
			b.time = time.Now().UTC()
			b.tokens = th.max - 1
		} else {
			deltaTokens := currentTokens - b.tokens
			deltaRefills := deltaTokens / th.refill
			deltaTime := time.Duration(deltaRefills) * th.d

			b.time = b.time.Add(deltaTime)
			b.tokens = currentTokens - 1
		}

		next.ServeHTTP(w, r)
	})
}
