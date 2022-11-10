package utils

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RunGrpcServer(ctx context.Context, srv *grpc.Server) error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Fatalf("Failed ot listen: %s", err)
	}

	errCh := make(chan error)
	shutdownCh := make(chan os.Signal, 1)

	signal.Notify(shutdownCh, syscall.SIGTERM, syscall.SIGINT)

	logger.WithField("address", listener.Addr().String()).Info("STARTED GRPC SERVER")

	go func() {
		if err := srv.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	defer func() {
		srv.GracefulStop()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("SERVER ERROR: %s", err)
	case sig := <-shutdownCh:
		logger.WithField("signal", sig.String()).Info("Signal interrupt")
	case <-ctx.Done():
		logger.Info("Context complete")
	}

	return nil
}
