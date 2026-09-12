package ckbplane

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const LoopbackServiceReadySchema = "wingless.loopback-service-ready.v1"

type LoopbackServiceReady struct {
	Schema   string `json:"schema"`
	Endpoint string `json:"endpoint"`
	PID      int    `json:"pid"`
}

// ServeLoopbackService exposes the already accepted Listener on an ephemeral
// numeric IPv4 loopback port. It does not register a worker, consume a queue,
// apply edits, retry work, or claim acceptance. The caller owns process and
// service lifetime.
func ServeLoopbackService(ctx context.Context, session string, engine Intelligence, ready io.Writer) error {
	if ctx == nil {
		return errors.New("WINGLESS_SERVICE_CONTEXT_REQUIRED")
	}
	if ready == nil {
		return errors.New("WINGLESS_SERVICE_READY_WRITER_REQUIRED")
	}
	listener, err := NewListener(session, engine)
	if err != nil {
		return err
	}

	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp4", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("WINGLESS_SERVICE_LISTEN_FAILED: %w", err)
	}

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok || tcpAddr.IP == nil || !tcpAddr.IP.IsLoopback() || tcpAddr.IP.To4() == nil {
		_ = ln.Close()
		return errors.New("WINGLESS_SERVICE_NONLOOPBACK_REFUSED")
	}

	server := &http.Server{
		Handler:           listener.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	serveDone := make(chan error, 1)
	go func() {
		err := server.Serve(ln)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveDone <- err
	}()

	announcement := LoopbackServiceReady{
		Schema:   LoopbackServiceReadySchema,
		Endpoint: "http://" + ln.Addr().String() + TransportPath,
		PID:      os.Getpid(),
	}
	payload, err := json.Marshal(announcement)
	if err != nil {
		_ = server.Close()
		<-serveDone
		return err
	}
	if _, err := ready.Write(append(payload, '\n')); err != nil {
		_ = server.Close()
		<-serveDone
		return fmt.Errorf("WINGLESS_SERVICE_READY_WRITE_FAILED: %w", err)
	}

	select {
	case err := <-serveDone:
		if err == nil {
			return errors.New("WINGLESS_SERVICE_STOPPED_UNEXPECTEDLY")
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
		}
		select {
		case err := <-serveDone:
			if err != nil {
				return err
			}
		case <-time.After(3 * time.Second):
			_ = server.Close()
			return errors.New("WINGLESS_SERVICE_SHUTDOWN_TIMEOUT")
		}
		return ctx.Err()
	}
}
