package main

import (
	"log"
	"net/http"

	limiter "github.com/ulule/limiter/v3"
	stdlib "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

// LimiterMiddleware provides limiter middleware pointer
var LimiterMiddleware *stdlib.Middleware

// initialize Limiter middleware pointer
func initLimiter(period string) {
	log.Printf("limiter rate='%s'", period)
	// create rate limiter with 5 req/second
	rate, err := limiter.NewRateFromFormatted(period)
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate, limiter.WithSkipList(Config.LimiterSkipList))
	if Config.LimiterHeader != "" {
		instance = limiter.New(
			store,
			rate,
			limiter.WithClientIPHeader(Config.LimiterHeader),
			limiter.WithSkipList(Config.LimiterSkipList))
	}
	LimiterMiddleware = stdlib.NewMiddleware(instance)
}

// helper to auth/authz incoming requests to the server
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// perform authentication
		status := CMSAuth.CheckAuthnAuthz(r.Header)
		if !status {
			log.Printf("ERROR: fail to authenticate, HTTP headers %+v\n", r.Header)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if Config.Verbose > 2 {
			log.Printf("Auth layer status: %v headers: %+v\n", status, r.Header)
		}

		// check if user has proper roles to DBS (non GET) APIs
		if r.Method != "GET" && Config.CMSRole != "" && Config.CMSGroup != "" {
			status = CMSAuth.CheckCMSAuthz(r.Header, Config.CMSRole, Config.CMSGroup, "")
			if !status {
				log.Printf("ERROR: fail to authorize used with role=%v and group=%v, HTTP headers %+v\n", Config.CMSRole, Config.CMSGroup, r.Header)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// limit middleware limits incoming requests
func limitMiddleware(next http.Handler) http.Handler {
	return LimiterMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}))
}
