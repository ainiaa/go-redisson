package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ainiaa/catutil/v2/rediscat"
	gredis "github.com/go-redis/redis/v8"

	"github.com/ainiaa/go-redission/conf"
)

var redisCluster = &Redis{}
var redisClients = make(map[string]*Redis, 2)

const (
	ConnTypeCluster  = "cluster"
	ConnTypeAlone    = "alone"
	ConnTypeSentinel = "sentinel"
)

type Redis struct {
	gredis.UniversalClient
	*conf.Config
}

func New(ctx context.Context, c *conf.Config) (*Redis, error) {
	switch c.ConnType {
	case ConnTypeCluster:
		return newClusterClient(ctx, c)
	case ConnTypeAlone:
		return newAloneClient(ctx, c)
	case ConnTypeSentinel: // todo
		return newSentinel(ctx, c)
	default:
		panic("unknown connection type")
	}
}

func newClusterClient(ctx context.Context, c *conf.Config) (*Redis, error) {
	client := gredis.NewClusterClient(&gredis.ClusterOptions{
		Addrs:              c.Cluster.Addrs,
		MaxRedirects:       c.Cluster.MaxRedirects,
		ReadOnly:           c.Cluster.ReadOnly,
		RouteByLatency:     c.Cluster.RouteByLatency,
		RouteRandomly:      c.Cluster.RouteRandomly,
		Username:           c.Cluster.Username,
		Password:           c.Cluster.Password,
		MaxRetries:         c.Cluster.MaxRetries,
		MinRetryBackoff:    time.Duration(c.Cluster.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff:    time.Duration(c.Cluster.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:        time.Duration(c.Cluster.DialTimeout) * time.Millisecond,
		ReadTimeout:        time.Duration(c.Cluster.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(c.Cluster.WriteTimeout) * time.Millisecond,
		PoolSize:           c.Cluster.PoolSize,
		MinIdleConns:       c.Cluster.MinIdleConns,
		MaxConnAge:         time.Duration(c.Cluster.MaxConnAge) * time.Millisecond,
		PoolTimeout:        time.Duration(c.Cluster.PoolTimeout) * time.Millisecond,
		IdleTimeout:        time.Duration(c.Cluster.IdleTimeout) * time.Millisecond,
		IdleCheckFrequency: time.Duration(c.Cluster.IdleCheckFrequency) * time.Millisecond,
	})
	pong, err := client.Ping(ctx).Result()
	if pong != "PONG" || err != nil {
		return nil, fmt.Errorf("cluster redis conn error: %s", err)
	}
	client.AddHook(rediscat.RedisTraceHook{})
	return &Redis{UniversalClient: client, Config: c}, nil
}

func newAloneClient(ctx context.Context, c *conf.Config) (*Redis, error) {
	client := gredis.NewClient(&gredis.Options{
		Network:            c.Alone.Network,
		Addr:               c.Alone.Addr,
		Username:           c.Alone.Username,
		Password:           c.Alone.Password,
		DB:                 c.Alone.DB,
		MaxRetries:         c.Alone.MaxRetries,
		MinRetryBackoff:    time.Duration(c.Alone.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff:    time.Duration(c.Alone.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:        time.Duration(c.Alone.DialTimeout) * time.Millisecond,
		ReadTimeout:        time.Duration(c.Alone.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(c.Alone.WriteTimeout) * time.Millisecond,
		PoolSize:           c.Alone.PoolSize,
		MinIdleConns:       c.Alone.MinIdleConns,
		MaxConnAge:         time.Duration(c.Alone.MaxConnAge) * time.Millisecond,
		PoolTimeout:        time.Duration(c.Alone.PoolTimeout) * time.Millisecond,
		IdleTimeout:        time.Duration(c.Alone.IdleTimeout) * time.Millisecond,
		IdleCheckFrequency: time.Duration(c.Alone.IdleCheckFrequency) * time.Millisecond,
	})
	pong, err := client.Ping(ctx).Result()
	if pong != "PONG" || err != nil {
		return nil, fmt.Errorf("alone redis conn error: %s", err)
	}
	client.AddHook(rediscat.RedisTraceHook{})
	return &Redis{UniversalClient: client, Config: c}, nil
}

func newSentinel(ctx context.Context, c *conf.Config) (*Redis, error) {
	client := gredis.NewRing(&gredis.RingOptions{
		Addrs:              c.Sentinel.Addrs,
		HeartbeatFrequency: time.Duration(c.Sentinel.HeartbeatFrequency) * time.Millisecond,
		DB:                 c.Sentinel.DB,
		Password:           c.Sentinel.Password,
		MaxRetries:         c.Sentinel.MaxRetries,
		MinRetryBackoff:    time.Duration(c.Sentinel.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff:    time.Duration(c.Sentinel.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:        time.Duration(c.Sentinel.DialTimeout) * time.Millisecond,
		ReadTimeout:        time.Duration(c.Sentinel.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(c.Sentinel.WriteTimeout) * time.Millisecond,
		PoolSize:           c.Sentinel.PoolSize,
		MinIdleConns:       c.Sentinel.MinIdleConns,
		MaxConnAge:         time.Duration(c.Sentinel.MaxConnAge) * time.Millisecond,
		PoolTimeout:        time.Duration(c.Sentinel.PoolTimeout) * time.Millisecond,
		IdleTimeout:        time.Duration(c.Sentinel.IdleTimeout) * time.Millisecond,
		IdleCheckFrequency: time.Duration(c.Sentinel.IdleCheckFrequency) * time.Millisecond,
	})
	pong, err := client.Ping(ctx).Result()
	if pong != "PONG" || err != nil {
		return nil, fmt.Errorf("sentinel redis conn error: %s", err)
	}
	client.AddHook(rediscat.RedisTraceHook{})
	return &Redis{UniversalClient: client, Config: c}, nil
}


func InitOnceRedis(ctx context.Context, c *conf.Config) (err error) {
	var once sync.Once
	once.Do(func() {
		redisCluster, err = New(ctx, c)
	})
	return
}

func InitNamedRedis(ctx context.Context, name string, c *conf.Config) (err error) {
	redisClients[name], err = New(ctx, c)
	return
}

func GetNamedRedis(name string) *Redis {
	r, ok := redisClients[name]
	if ok {
		return r
	}
	return nil
}

func GetRedis() *Redis {
	if IsRedisReady() {
		panic("redis is not ready,please init it first")
	}
	return redisCluster
}

func IsRedisReady() bool {
	return redisCluster == nil || redisCluster == (&Redis{})
}
