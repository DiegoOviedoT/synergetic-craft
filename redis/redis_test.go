//go:build integration
// +build integration

package redis

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedis_Set_and_Get(t *testing.T) {
	t.Run("should return success when Set and Get return success", func(t *testing.T) {
		ctx := context.Background()

		f := setupRedisFixture()

		err := f.client.Set(ctx, "test_set", "Hi, Set", time.Duration(2)*time.Second)
		assert.NoError(t, err)

		value, err := f.client.Get(ctx, "test_set")

		assert.NoError(t, err)
		assert.Equal(t, "Hi, Set", value)
	})
}

func TestRedis_MSet_and_MGet(t *testing.T) {
	t.Run("should return success when MSet and MGet return success", func(t *testing.T) {
		ctx := context.Background()

		f := setupRedisFixture()

		values := make(map[string]interface{}, 4)
		values["test_mset_1"] = "MSet 1 Redis"
		values["test_mset_2"] = "MSet 2 Redis"
		values["test_mset_3"] = "MSet 3 Redis"
		values["test_mset_4"] = "MSet 4 Redis"

		err := f.client.MSet(ctx, values, time.Duration(3)*time.Second)
		assert.NoError(t, err)

		response, err := f.client.MGet(ctx, "test_mset_1", "test_mset_2", "test_mset_3", "test_mset_4")

		assert.NoError(t, err)
		assert.Equal(t, "MSet 1 Redis", response["test_mset_1"])
		assert.Equal(t, "MSet 2 Redis", response["test_mset_2"])
		assert.Equal(t, "MSet 3 Redis", response["test_mset_3"])
		assert.Equal(t, "MSet 4 Redis", response["test_mset_4"])
	})
}

type redisFixture struct {
	client ClientRedis
}

func setupRedisFixture() *redisFixture {
	conf := Config{
		Address:       "localhost:6379",
		Password:      "passwordtest",
		TimeoutMillis: 3000,
		PoolSize:      3,
		DB:            0,
	}

	return &redisFixture{
		client: NewClient(conf),
	}
}
