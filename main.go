package main

import (
	"context"
	"fmt"
	mgrpc "github.com/mirshahriar/multiplexing-simple/grpc"
	mhttp "github.com/mirshahriar/multiplexing-simple/http"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	grpcServer := mgrpc.NewGRPCServer()
	httpServer := mhttp.NewHTTPServer()

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	combinedRouter := getCombinedRouter(grpcServer, httpServer.Handler)
	server := &http.Server{Handler: h2c.NewHandler(combinedRouter, &http2.Server{})}

	fmt.Println("Running server on port :8080")
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		_ = httpServer.Close()
		_ = lis.Close()
	}()

	_ = server.Serve(lis)
}

func getCombinedRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 {
			fmt.Println("gRPC request")
			grpcServer.ServeHTTP(w, r)
		} else {
			fmt.Println("HTTP request")
			httpHandler.ServeHTTP(w, r)
		}
	})
}
