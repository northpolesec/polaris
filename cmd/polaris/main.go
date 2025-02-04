package main

import (
	"log"
	"net"
	"os"

	"github.com/northpolesec/polaris/internal/statsservice"
	"google.golang.org/grpc"

	apipb "buf.build/gen/go/northpolesec/protos/grpc/go/stats/statsv1grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	projectId := os.Getenv("POLARIS_PROJECT_ID")
	datasetId := os.Getenv("POLARIS_DATASET_ID")
	tableId := os.Getenv("POLARIS_TABLE_ID")
	streamId := os.Getenv("POLARIS_STREAM_ID")

	// Create server and register StatsService.
	s := grpc.NewServer()
	svc, err := statsservice.NewStatsServiceServer(projectId, datasetId, tableId, streamId)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	apipb.RegisterStatsServiceServer(s, svc)

	// Create listener.
	lis, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	// Serve forever.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
