package db

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Service) SetRedisConn() {
	s.conn = s.redisClient.Conn()
}

func (s *Service) SetRedisClient(client *redis.Client) {
	s.redisClient = client
}

func (s *Service) Set(key, value string) error {
	if s.conn == nil {
		s.SetRedisConn()
	}

	err := s.conn.Set(s.ctx, key, value, time.Hour).Err()
	if err != nil {
		log.WithError(err).Error("Failed to set data")
		return err
	}

	return nil
}

func (s *Service) Get(key string) (string, error) {
	if s.conn == nil {
		s.SetRedisConn()
	}

	val, err := s.conn.Get(s.ctx, key).Result()

	if err != nil {
		if err == redis.TxFailedErr {
			// Tx가 실패한 경우
			return "", err
		}
		if err == redis.Nil {
			// Key 값이 존재하지 않을 때
			return "", nil
		}
		if err.Error() == "redis: client is closed" {
			// Client가 닫혀 있을 때
			s.SetRedisConn()
			s.Get(key)
		}
		return "", nil
	}

	return val, nil
}

func (s *Service) Del(key string) error {
	if s.conn == nil {
		s.SetRedisConn()
	}

	err := s.conn.Del(s.ctx, key).Err()
	if err != nil {
		log.WithError(err).Error("Failed to delete data")
		return err
	}

	return nil
}
