package grpc_client_pool

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"user/pkg/log"
	"user/pkg/zerror"
)

/*
通过sync.pool 实现grpc连接池，存储协程溢出的问题
*/
type clientPool struct {
	pool sync.Pool
}

func NewPool(target string, opts ...grpc.DialOption) (ClientPool, error) {
	return &clientPool{
		pool: sync.Pool{
			New: func() any {
				conn, err := grpc.Dial(target, opts...)
				if err != nil {
					log.Error(zerror.NewByErr(err))
					return nil
				}
				return conn
			},
		},
	}, nil
}

func (c *clientPool) Get() *grpc.ClientConn {
	conn := c.pool.Get().(*grpc.ClientConn)
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = c.pool.New().(*grpc.ClientConn)
	}
	return conn
}
func (c *clientPool) Put(conn *grpc.ClientConn) {
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		return
	}
	c.pool.Put(conn)
}
