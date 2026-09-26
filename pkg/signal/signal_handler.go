package signal

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const stopTimeout = 5 * time.Second

// graceful shutdown
func NewSignalHandler(log *slog.Logger, srv *http.Server) <-chan struct{} {
	done := make(chan struct{})

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-termSignal
		log.Warn("stop trigger", "signal", sig.String())

		if srv != nil {
			log.Warn("stopping http server")
			ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Error("can't stop http server gracefully", "error", err)
			}
		}

		close(done)
	}()

	return done
}
