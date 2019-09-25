// go:generate protoc -I ../api-common/protocol/proto --go_out=plugins=grpc:../api-common/protocol/rpc ../api-common/protocol/proto/trading.proto ../api-common/protocol/proto/common.proto ../api-common/protocol/proto/account.proto ../api-common/protocol/proto/error.proto ../api-common/protocol/proto/ws.proto ../api-common/protocol/proto/feed.proto

package protocol

import (
	"context"
	"errors"
	"fmt"
	"github.com/pborman/uuid"
	"google.golang.org/grpc/metadata"
	"sync"
	"time"
)

// PoolConfig 连接池相关配置
type GrpcPoolConfig struct {
	// 连接池中拥有的最小连接数
	InitialCap int
	// 连接池中拥有的最大的连接数
	MaxCap int
	// 生成连接的方法
	Factory func() (interface{}, error)
	// 关闭链接的方法
	Close func(interface{}) error
	// 链接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
}

// channelPool 存放链接信息
type GrpcPool struct {
	serverId    string
	mu          sync.Mutex
	conns       chan *idleGrpcConn
	factory     func() (interface{}, error)
	close       func(interface{}) error
	idleTimeout time.Duration
}

type idleGrpcConn struct {
	conn interface{}
	ctx  context.Context
	t    time.Time
}

// NewChannelPool 初始化链接
func NewGrpcPool(grpcPoolConfig *GrpcPoolConfig) (*GrpcPool, error) {
	if grpcPoolConfig.InitialCap < 0 || grpcPoolConfig.MaxCap <= 0 || grpcPoolConfig.InitialCap > grpcPoolConfig.MaxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &GrpcPool{
		serverId:    uuid.NewRandom().String(),
		conns:       make(chan *idleGrpcConn, grpcPoolConfig.MaxCap),
		factory:     grpcPoolConfig.Factory,
		close:       grpcPoolConfig.Close,
		idleTimeout: grpcPoolConfig.IdleTimeout,
	}

	for i := 0; i < grpcPoolConfig.InitialCap; i++ {
		conn, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleGrpcConn{conn: conn, ctx: c.getContext(), t: time.Now()}
	}

	return c, nil
}

func (c *GrpcPool) getContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-server-data", c.serverId, "x-request-data", uuid.NewRandom().String())
	//defer cancel()
	return ctx
}

// getConns 获取所有连接
func (c *GrpcPool) getConns() chan *idleGrpcConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

func (c *GrpcPool) GetServerId() string {
	return c.serverId
}

// Get 从pool中取一个连接
func (c *GrpcPool) Get() (interface{}, context.Context, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, nil, ErrClosed
			}
			// 判断是否超时，超时则丢弃
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					// 丢弃并关闭该链接
					c.Close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, wrapConn.ctx, nil
		default:
			conn, err := c.factory()
			if err != nil {
				return nil, nil, err
			}

			return conn, c.getContext(), nil
		}
	}
}

// Put 将连接放回pool中
func (c *GrpcPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		return c.Close(conn)
	}

	select {
	case c.conns <- &idleGrpcConn{conn: conn, t: time.Now()}:
		return nil
	default:
		// 连接池已满，直接关闭该链接
		return c.Close(conn)
	}
}

// Close 关闭单条连接
func (c *GrpcPool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.close(conn)
}

// Release 释放连接池中所有链接
func (c *GrpcPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

// Len 连接池中已有的连接
func (c *GrpcPool) Len() int {
	return len(c.getConns())
}
