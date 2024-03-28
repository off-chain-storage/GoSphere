package db

import (
	"context"
	"testing"

	dbTestHelper "github.com/off-chain-storage/GoSphere/go-sphere/db/testing"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testKey   = "test"
	testValue = "1"
)

func setupTest(t *testing.T) (*Service, *redis.Client) {
	// Create a new mock redis client
	client := dbTestHelper.NewMockRedis(t)

	// Create a new redis client
	svc, err := NewRedisClient(context.Background(), &Config{})
	require.NoError(t, err, "creating new redis client should not error")

	// Set the mock redis client
	svc.SetRedisClient(client)
	return svc, client
}

func TestSet(t *testing.T) {
	svc, client := setupTest(t)

	// Set value with my function(Set)
	err := svc.Set(testKey, testValue)
	require.NoError(t, err, "setting value should not error")

	// Get value with mock redis client
	actual, err := client.Get(context.Background(), testKey).Result()
	require.NoError(t, err, "getting value should not error")

	// Compare expected value and actual value
	assert.Equal(t, testValue, actual, "expected value to match")
}

func TestGet(t *testing.T) {
	svc, client := setupTest(t)

	// Set value with mock redis client
	err := client.Set(context.Background(), testKey, testValue, 0).Err()
	require.NoError(t, err, "setting value with mock client should not error")

	// Get value with my function(Get)
	actual, err := svc.Get(testKey)
	require.NoError(t, err, "getting value should not error")

	// Compare expected value and actual value
	assert.Equal(t, testValue, actual, "expected value to match")
}

func TestDel(t *testing.T) {
	svc, client := setupTest(t)

	// Set value with mock redis client
	err := client.Set(context.Background(), testKey, testValue, 0).Err()
	require.NoError(t, err, "setting value with mock client should not error")

	// Delete value with my function(Del)
	err = svc.Del(testKey)
	require.NoError(t, err, "deleting value should not error")

	// Get value with mock redis client
	_, err = client.Get(context.Background(), testKey).Result()
	require.Error(t, err, "expected an error when getting a deleted key")
	assert.Equal(t, err.Error(), "redis: nil", "expected 'redis: nil' error for deleted key")
}
