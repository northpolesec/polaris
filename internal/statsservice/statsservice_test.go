package statsservice_test

import (
	"context"
	"flag"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-replayers/grpcreplay"
	"github.com/northpolesec/polaris/internal/statsservice"
	"github.com/shoenig/test/must"
	"google.golang.org/api/option"

	apipb "buf.build/gen/go/northpolesec/protos/protocolbuffers/go/stats"
)

const (
	replayFilename = "testdata/bigquery.replay"

	projectId = "polaris-449516"
	datasetId = "testset"
	tableId   = "testtable"
	streamId  = "_default"
)

var (
	update = flag.Bool("update", false, "update golden file(s)")
)

func TestSubmitStats_SimpleSubmission(t *testing.T) {
	ss, err := statsservice.NewStatsServiceServer(projectId, datasetId, tableId, streamId, dialOptionsForTest(t)...)
	must.NoError(t, err)

	// Make a single SubmitStats request, verify that it doesn't report an error.
	_, err = ss.SubmitStats(context.Background(), connect.NewRequest(&apipb.SubmitStatsRequest{
		SantaVersion:  "2025.1",
		MachineIdHash: "c9bcf04f8d69279ad881b6d9467785ea72d99c561976892d2e39f63b4c1df4b4",
		MacosVersion:  "15.2",
		MacosBuild:    "24C101",
		MacModel:      "Mac15,7",
	}))
	must.NoError(t, err)
}

// dialOptionsForTest returns DialOptions that can be passed to the BigQuery API when creating a client.
// Depending on the value for the `-update` flag this will attach either a recorder or replayer to this test.
func dialOptionsForTest(t *testing.T) []option.ClientOption {
	t.Helper()

	// If the user passed the -update flag, use a recorder to update the golden file.
	if *update {
		recorder, err := grpcreplay.NewRecorder(replayFilename, &grpcreplay.RecorderOptions{Text: true})
		must.NoError(t, err)
		t.Cleanup(func() {
			must.NoError(t, recorder.Close())
		})

		opts := make([]option.ClientOption, len(recorder.DialOptions()))
		for i, v := range recorder.DialOptions() {
			opts[i] = option.WithGRPCDialOption(v)
		}
		return opts
	}

	// User didn't pass -update, so replay the existing golden.
	replayer, err := grpcreplay.NewReplayer(replayFilename, nil)
	must.NoError(t, err)
	t.Cleanup(func() {
		must.NoError(t, replayer.Close())
	})
	conn, err := replayer.Connection()
	must.NoError(t, err)
	return []option.ClientOption{option.WithGRPCConn(conn)}
}
