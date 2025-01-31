package main

import (
	"log"
	"net"
	"os"

	"github.com/northpolesec/polaris/internal/statsservice"
	"google.golang.org/grpc"

	apipb "buf.build/gen/go/northpolesec/protos/grpc/go/stats/statsv1grpc"
)

const (
	projectId = "polaris-449516"
	datasetId = "polaris"
	tableId   = "stats"
	streamId  = "_default"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create listener.
	lis, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	// Create server and register StatsService.
	s := grpc.NewServer()
	svc, err := statsservice.NewStatsServiceServer(projectId, datasetId, tableId, streamId)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	apipb.RegisterStatsServiceServer(s, svc)

	// Serve forever.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
