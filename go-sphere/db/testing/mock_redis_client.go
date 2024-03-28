package dbTestHelper

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func NewMockRedis(t *testing.T) (client *redis.Client) {
	t.Helper()

	// Build a mock redis server
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("unexpected error while creating test redis server '%#v'", err)
	}

	// Create a new mock redis client
	client = redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})

	return
}
