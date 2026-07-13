package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServerAndGracefullyShutdown(srv *http.Server, timeout time.Duration) {
	// start in a seperate goroutine
	go func() {
		serverStartError := srv.ListenAndServe()
		if serverStartError != nil && serverStartError != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", serverStartError)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		fmt.Fprintf(os.Stderr, "Server forced to shutdown: %v\n", shutdownErr)
	}

	fmt.Println("Done")
}
