package redis

import (
	"context"
	rds "github.com/redis/go-redis/v9"
	"time"
)

type ClientRedis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key ...string) error
	MSet(ctx context.Context, values map[string]interface{}, ttl time.Duration) error
	MGet(ctx context.Context, keys ...string) (map[string]string, error)
}

type Config struct {
	Address       string
	Password      string
	TimeoutMillis int64
	PoolSize      int
	DB            int
}

type redis struct {
	client rds.Cmdable
}

func NewClient(conf Config) ClientRedis {
	return &redis{
		client: initClient(conf),
	}
}

func initClient(conf Config) *rds.Client {
	timeout := time.Duration(conf.TimeoutMillis) * time.Millisecond

	opts := &rds.Options{
		Addr:         conf.Address,
		Password:     conf.Password,
		PoolSize:     conf.PoolSize,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		PoolTimeout:  timeout,
		DB:           conf.DB,
	}

	client := rds.NewClient(opts)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return client
}

func (r *redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redis) Del(ctx context.Context, key ...string) error {
	return r.client.Del(ctx, key...).Err()
}

func (r *redis) MSet(ctx context.Context, values map[string]interface{}, ttl time.Duration) error {
	var args []interface{}
	for key, value := range values {
		args = append(args, key, value)
	}

	tx := r.client.TxPipeline()

	err := r.client.MSet(ctx, args...).Err()
	if err != nil {
		return err
	}

	for key := range values {
		err = r.client.Expire(ctx, key, ttl).Err()
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

func (r *redis) MGet(ctx context.Context, keys ...string) (map[string]string, error) {
	result, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]string)
	for i, key := range keys {
		resultMap[key] = result[i].(string)
	}

	return resultMap, nil
}
