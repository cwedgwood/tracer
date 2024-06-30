// contrived example using tracer
//
// I would blame an AI for this contrived example, but it's suggestion was actually better than this...

/*
   Try some of:

   curl "http://localhost:28080" ; echo
   curl --cookie "AUTH=eoj" "http://localhost:28080?username=joe" ; echo
   curl --cookie "AUTH=bad" "http://localhost:28080?username=joe" ; echo

*/

package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cwedgwood/tracer"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

func runServer() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zl := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	baseLogger := zerologr.New(&zl)

	server := http.Server{
		Addr: "localhost:28080",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = tracer.ContextLoggerWithTraceId(ctx, baseLogger, "", "example.server")
		doRequest(ctx, w, r)
	})
	server.ListenAndServe()
}

func doRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("request received")

	// parse request arguments from URL
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "anonymous"
	}
	logger = logger.WithValues("username", username)
	ctx = logr.NewContext(ctx, logger)
	// get auth cooking
	var authValue string
	authCookie, err := r.Cookie("AUTH")
	if err != nil {
		logger.Error(err, "failed to get auth cookie, assuming anonymous")
	} else {
		authValue = authCookie.Value
		logger = logger.WithValues("auth", authValue)
		ctx = logr.NewContext(ctx, logger)
	}
	if authCookie != nil {
		logger = logger.WithValues("auth", authValue)
		ctx = logr.NewContext(ctx, logger)
	}

	userallowed := authenticate(ctx, username, authValue)
	// log if the user is allowed
	if userallowed {
		logger.Info("user allowed")
	} else {
		logger.Info("user not allowed")
	}
	if !userallowed {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	authorize(ctx, username)
	accounting(ctx, username)
	response(ctx, username, w)
}

func authenticate(ctx context.Context, username string, auth string) bool {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("authenticate")
	// anonymous is always allowed
	if username == "anonymous" {
		return true
	}
	// create an array of runes from the username
	reversedRunes := []rune{}
	userRune := []rune(username)
	// iterate backwards through the userRune
	for i := len(userRune) - 1; i >= 0; i-- {
		reversedRunes = append(reversedRunes, userRune[i])
	}
	revUsername := string(reversedRunes)
	// if the reversed username is the same as the auth cookie, allow
	return revUsername == auth
}

func authorize(ctx context.Context, username string) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("authorize")
}

func accounting(ctx context.Context, username string) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("accounting")
}

func response(ctx context.Context, username string, w http.ResponseWriter) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("response")
	// write a response saying hello to whomever they are
	w.Write([]byte("Hello, " + username))
}

func main() {
	runServer()
}
