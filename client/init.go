package client

import (
	"fmt"
	"sync"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/common/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn           *grpc.ClientConn
	client         userpb.UserServiceClient
	clientSyncOnce sync.Once
)

func GetUserServiceClient(config *GRpcClientConfig) userpb.UserServiceClient {
	clientSyncOnce.Do(func() {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(1024 * 1024)),
		}
		conn, _ = grpc.NewClient(fmt.Sprintf("%s:%d", config.Host, config.Port), opts...)
		client = userpb.NewUserServiceClient(conn)
	})
	return client
}

func Destroy() {
	if conn != nil {
		conn.Close()
	}
}
