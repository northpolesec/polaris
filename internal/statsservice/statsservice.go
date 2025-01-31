package statsservice

import (
	"context"
	"fmt"
	"log"

	bqStorage "cloud.google.com/go/bigquery/storage/apiv1"
	"cloud.google.com/go/bigquery/storage/apiv1/storagepb"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	apipb "buf.build/gen/go/northpolesec/protos/protocolbuffers/go/stats"
)

type StatsServiceServer struct {
	writeClient *bqStorage.BigQueryWriteClient
	descriptor  *descriptorpb.DescriptorProto
	streamName  string
}

func NewStatsServiceServer(projectId, datasetId, tableId, streamId string) (*StatsServiceServer, error) {
	wc, err := bqStorage.NewBigQueryWriteClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to make BigQuery client: %w", err)
	}

	// Create the descriptor from the proto, so that we don't have to re-calculate it on every
	// submission.
	emptyReq := apipb.SubmitStatsRequest{}
	descriptor, err := adapt.NormalizeDescriptor(emptyReq.ProtoReflect().Descriptor())
	if err != nil {
		return nil, fmt.Errorf("failed to create descriptor proto: %w", err)
	}

	return &StatsServiceServer{
		writeClient: wc,
		descriptor:  descriptor,
		streamName:  fmt.Sprintf("projects/%s/datasets/%s/tables/%s/streams/%s", projectId, datasetId, tableId, streamId),
	}, nil
}

func (s *StatsServiceServer) SubmitStats(ctx context.Context, req *apipb.SubmitStatsRequest) (*apipb.SubmitStatsResponse, error) {
	// Set the submission time for this request.
	req.SetSubmitTime(timestamppb.Now())

	// Get a write stream.
	ws, err := s.writeClient.GetWriteStream(ctx, &storagepb.GetWriteStreamRequest{Name: s.streamName})
	if err != nil {
		log.Printf("Failed to GetWriteStream(): %v", err)
		return nil, err
	}

	// Create a stream to send new rows.
	stream, err := s.writeClient.AppendRows(ctx)
	if err != nil {
		log.Printf("Failed to AppendRows(): %v", err)
		return nil, err
	}
	defer stream.CloseSend()

	// Marshal the row data.
	rows, err := s.appendRowsReqProtoRowsFromSubmitStatsReq(req)
	if err != nil {
		log.Printf("Failed to create AppendRowsRequest: %v", err)
		return nil, err
	}

	// Send our new row.
	if err := stream.Send(&storagepb.AppendRowsRequest{
		WriteStream: ws.Name,
		TraceId:     "polaris",
		Rows:        rows,
	}); err != nil {
		log.Printf("Failed to send row to stream: %v", err)
		return nil, err
	}

	// Ensure the append was successful.
	if _, err := stream.Recv(); err != nil {
		log.Printf("Appending row failed: %v", err)
		return nil, err
	}

	return &apipb.SubmitStatsResponse{}, nil
}

func (s *StatsServiceServer) appendRowsReqProtoRowsFromSubmitStatsReq(req *apipb.SubmitStatsRequest) (*storagepb.AppendRowsRequest_ProtoRows, error) {
	// Marshal the row data.
	rowData, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Return the AppendRowsRequest_ProtoRows message.
	return &storagepb.AppendRowsRequest_ProtoRows{
		ProtoRows: &storagepb.AppendRowsRequest_ProtoData{
			WriterSchema: &storagepb.ProtoSchema{
				ProtoDescriptor: s.descriptor,
			},
			Rows: &storagepb.ProtoRows{
				SerializedRows: [][]byte{rowData},
			},
		},
	}, nil
}
