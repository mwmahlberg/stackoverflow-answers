package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-retryablehttp"
)

var bind = flag.String("bind", ":8080", "address to bind to")
var target = flag.String("target", "http://localhost:8081", "target address to forward requests to")
var attempts = flag.Uint("attempts", 3, "number of attempts to retry")

func main() {
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// retryablehttp only logs at debug level
		// and we want to see all logs for demonstration purposes
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	// Sanitize the target URL
	targetUrl, err := url.Parse(*target)
	if err != nil {
		logger.Error("failed to parse target URL", "error", err)
		os.Exit(1)
	}

	// Set up the retryable HTTP client
	retryClient := retryablehttp.NewClient()

	// A quick and dirty way to log the requests
	// In production, you probably want to use a more sophisticated logging solutions
	retryClient.Logger = logger.WithGroup("retryablehttp")
	retryClient.RetryMax = int(*attempts)

	// The Backoff heavily depends on your use case.
	// Note however that Backoff is just an alias for
	// func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration
	// and you can easily implement your own backoff function.
	retryClient.Backoff = retryablehttp.LinearJitterBackoff

	// You can customize when to retry.
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// This is probably a very bad idea.
		// It is only here for illustration purposes.
		if resp.StatusCode >= 400 {
			return true, nil
		}
		return false, nil
	}

	// Configure the reverse proxy
	pr := &httputil.ReverseProxy{
		// This is crucial: Since we insert a transport,
		// and a transport cannot deal with a RequestURI,
		// (whch is not a part of the client request, but a part of
		// the request as parsed by the server),
		// we need to set the request URI to empty.
		Rewrite: func(req *httputil.ProxyRequest) {
			req.Out.URL = req.In.URL
			req.Out.URL.Host = targetUrl.Host
			req.Out.URL.Scheme = targetUrl.Scheme
			req.Out.RequestURI = ""
		},
		Transport: retryClient.StandardClient().Transport,
	}

	// With everything set up, we can configure the server
	// to use our reverse proxy
	srv := &http.Server{
		Addr:    *bind,
		Handler: pr,
	}

	// From here, it is pretty much standard HTTP server code
	go func() {
		logger.Info("Starting server", "address", *bind)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			os.Exit(2)
		}
		logger.Info("Server stopped")
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
	logger.Info("Shutting down server")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("Server shutdown error", "error", err)
		os.Exit(3)
	}
	logger.Info("Server gracefully stopped")

}
