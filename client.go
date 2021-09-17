package alphaticks

import (
	"context"
	"fmt"
	"gitlab.com/tachikoma.ai/tickstore-go-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	DB_LIVE = "live"
	DB_1S   = "1s"
	DB_1M   = "1m"
	DB_1H   = "1h"
	DB_1D   = "1d"
)

const (
	DATA_CLIENT_LIVE int64 = 0
	DATA_CLIENT_1S   int64 = 1000
	DATA_CLIENT_1M   int64 = DATA_CLIENT_1S * 60
	DATA_CLIENT_1H   int64 = DATA_CLIENT_1M * 60
	DATA_CLIENT_1D   int64 = DATA_CLIENT_1H * 24
)

var ports = map[int64]string{
	DATA_CLIENT_LIVE: "3550",
	DATA_CLIENT_1S:   "3551",
	DATA_CLIENT_1M:   "3552",
	DATA_CLIENT_1H:   "3553",
	DATA_CLIENT_1D:   "3554",
}

var freqTonames = map[int64]string{
	DATA_CLIENT_LIVE: "live",
	DATA_CLIENT_1S:   "1s",
	DATA_CLIENT_1M:   "1m",
	DATA_CLIENT_1H:   "1h",
	DATA_CLIENT_1D:   "1d",
}

var namesToFreq = map[string]int64{
	"live": DATA_CLIENT_LIVE,
	"1s":   DATA_CLIENT_1S,
	"1m":   DATA_CLIENT_1M,
	"1h":   DATA_CLIENT_1H,
	"1d":   DATA_CLIENT_1D,
}

var shardDurations = map[int64]uint64{
	DATA_CLIENT_LIVE: 10000000,
	DATA_CLIENT_1S:   10000000,
	DATA_CLIENT_1M:   1000000000,
	DATA_CLIENT_1H:   10000000000,
	DATA_CLIENT_1D:   100000000000,
}

type DataClient struct {
	store        *tickstore_go_client.RemoteClient
	measurements map[string]string
}

func NewClient(licenseID, licenseKey, db string) (*DataClient, error) {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Append license-id and license-key
		md := metadata.New(map[string]string{"license-id": licenseID, "license-key": licenseKey})
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}))
	opts = append(opts, grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		// Append license-id and license-key
		md := metadata.New(map[string]string{"license-id": licenseID, "license-key": licenseKey})
		ctx = metadata.NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, cc, method, opts...)
	}))

	freq, ok := namesToFreq[db]
	if !ok {
		return nil, fmt.Errorf("unknown db %s", db)
	}
	s := &DataClient{}
	str, err := tickstore_go_client.NewRemoteClient("store.alphaticks.io"+":"+ports[freq], opts...)
	if err != nil {
		return nil, err
	}
	s.store = str

	return s, nil
}

func (c *DataClient) Query(qs *QuerySettings) (tickstore_go_client.TickstoreQuery, error) {
	q, err := c.store.NewQuery(qs.settings)
	if err != nil {
		return nil, fmt.Errorf("error querying store: %v", err)
	}
	return q, nil
}
