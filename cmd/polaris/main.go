package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/northpolesec/polaris/internal/statsservice"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	apipb "buf.build/gen/go/northpolesec/protos/connectrpc/gosimple/stats/statsv1connect"
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
	svc, err := statsservice.NewStatsServiceServer(projectId, datasetId, tableId, streamId)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(apipb.NewStatsServiceHandler(svc))

	if err := http.ListenAndServe(
		net.JoinHostPort("0.0.0.0", port),
		// This wrapping of mux is necessary to make non-TLS HTTP/2 connections work.
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
