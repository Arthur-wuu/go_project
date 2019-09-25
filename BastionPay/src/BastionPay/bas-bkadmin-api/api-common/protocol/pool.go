package protocol

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrClosed 连接池已经关闭Error
	ErrClosed = errors.New("pool is closed")
)

// Pool 基本方法
type Pool interface {
	Get() (interface{}, interface{}, error)

	Put(interface{}) error

	Close(interface{}) error

	Release()

	Len() int
}

// PoolConfig 连接池相关配置
type PoolConfig struct {
	// 连接池中拥有的最小连接数
	InitialCap int
	// 连接池中拥有的最大的连接数
	MaxCap int
	// 生成连接的方法
	Factory func() (interface{}, interface{}, error)
	// 关闭链接的方法
	Close func(interface{}) error
	// 链接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
}

// channelPool 存放链接信息
type channelPool struct {
	mu          sync.Mutex
	conns       chan *idleConn
	factory     func() (interface{}, interface{}, error)
	close       func(interface{}) error
	idleTimeout time.Duration
}

type idleConn struct {
	conn      interface{}
	channelId interface{}
	t         time.Time
}

// NewChannelPool 初始化链接
func NewChannelPool(poolConfig *PoolConfig) (Pool, error) {
	if poolConfig.InitialCap < 0 || poolConfig.MaxCap <= 0 || poolConfig.InitialCap > poolConfig.MaxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &channelPool{
		conns:       make(chan *idleConn, poolConfig.MaxCap),
		factory:     poolConfig.Factory,
		close:       poolConfig.Close,
		idleTimeout: poolConfig.IdleTimeout,
	}

	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, channelId, err := c.factory()
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleConn{conn: conn, channelId: channelId, t: time.Now()}
	}

	return c, nil
}

// getConns 获取所有连接
func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

// Get 从pool中取一个连接
func (c *channelPool) Get() (interface{}, interface{}, error) {
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
			return wrapConn.conn, wrapConn.channelId, nil
		default:
			conn, channelId, err := c.factory()
			if err != nil {
				return nil, nil, err
			}

			return conn, channelId, nil
		}
	}
}

// Put 将连接放回pool中
func (c *channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		return c.Close(conn)
	}

	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		// 连接池已满，直接关闭该链接
		return c.Close(conn)
	}
}

// Close 关闭单条连接
func (c *channelPool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.close(conn)
}

// Release 释放连接池中所有链接
func (c *channelPool) Release() {
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
func (c *channelPool) Len() int {
	return len(c.getConns())
}
