//go:generate protoc proto/rpc.proto --go_out=plugins=grpc:.
package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
)

func main() {
    // Create new server
    srv := grpc.NewServer()

    // Register protobuf service with server
    // pb.RegisterXXXServer(srv, &pb.Server{})

    // Now listening on: http://{{ .Host }}:{{ .Port }}
    lis, err := net.Listen("tcp", net.JoinHostPort("{{ .Host }}", "{{ .Port }}"))
    if err != nil {
        log.Fatal(err)
    }

    // Application started. Press CTRL+C to shut down.
    srv.Serve(lis)
}
