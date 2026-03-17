package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type EtcdClient interface {
	Incr(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
	Close() error
}

type etcdClient struct {
	client *clientv3.Client
}

func NewEtcdClient(endpoints []string) (EtcdClient, error) {
	// 创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &etcdClient{
		client: cli,
	}, nil
}

// Incr 方法用于原子自增
func (ec *etcdClient) Incr(ctx context.Context, key string) (int64, error) {
	// 创建一个会话用于分布式锁
	session, err := concurrency.NewSession(ec.client)
	if err != nil {
		return 0, err
	}
	defer session.Close()

	// 使用分布式锁确保操作的原子性
	mutex := concurrency.NewMutex(session, "/mutex/"+key)

	// 尝试获取锁
	if err := mutex.Lock(ctx); err != nil {
		return 0, err
	}
	defer func() {
		// 释放锁
		_ = mutex.Unlock(ctx)
	}()

	// 获取当前值
	resp, err := ec.client.Get(context.Background(), key)
	if err != nil {
		return 0, err
	}

	var newValue int64
	if len(resp.Kvs) == 0 {
		// 如果 key 不存在，初始值设为 1
		newValue = 1
	} else {
		// 解析当前值并自增
		currentValue, err := strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
		if err != nil {
			return 0, err
		}
		newValue = currentValue + 1
	}
	// 将新值写回 etcd
	_, err = ec.client.Put(context.Background(), key, strconv.FormatInt(newValue, 10))
	if err != nil {
		return 0, err
	}
	logx.WithContext(ctx).Infof("incr key: %s, value: %d", key, newValue)
	return newValue, nil
}

// Close 关闭客户端的方法
func (ec *etcdClient) Close() error {
	return ec.client.Close()
}

func (ec *etcdClient) Set(ctx context.Context, key string, value any) error {
	_, err := ec.client.Put(ctx, key, fmt.Sprintf("%v", value))
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("set key: %s, value: %+v", key, value)
	return err
}

func (ec *etcdClient) Get(ctx context.Context, key string) (any, error) {
	resp, err := ec.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, code.EtcdKeyNotExist
	}
	logx.WithContext(ctx).Infof("get key: %s, value: %+v", key, resp.Kvs[0].Value)
	return resp.Kvs[0].Value, nil
}
