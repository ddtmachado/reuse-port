package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

var (
	id            uuid.UUID
	serverAddress = ":9696"
)

func main() {
	id = uuid.New()

	netConfig := net.ListenConfig{
		Control: netControl,
	}

	listener, err := netConfig.Listen(context.Background(), "tcp", serverAddress)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Server %s\n", id)
	})
	mux.HandleFunc("/wait", waitHandler)

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve on %s: %v", serverAddress, err)
		}
	}()

	<-handleSystemSignals(server)
}

// handleSystemSignals capture unix SIGINT and SIGTERM to gracefully shutdown the server process
func handleSystemSignals(s *http.Server) chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs

		log.Printf("%s received, shutting down gracefully.\n", sig)

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to shutdown gracefully: %v", err)
		}

		done <- true
	}()

	return done
}

// netControl is a Control function for net.ListenConfig that enables unix socket reuse for port and address.
func netControl(network, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		if err != nil {
			return
		}

		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if err != nil {
			return
		}
	})
	return err
}

func waitHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(r.URL.String())
	queryParams := u.Query()
	waitTime := queryParams.Get("time")

	duration, err := time.ParseDuration(waitTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	time.Sleep(duration)

	_, _ = fmt.Fprintf(w, "Server %s waited for: %s\n", id, duration.String())
}
