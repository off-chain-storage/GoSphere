package iface

type ReadOnlyRedisDB interface {
	Get(key string) (string, error)
}

type AccessRedisDB interface {
	ReadOnlyRedisDB

	Set(key, value string) error
	Del(key string) error
}

type RedisDB interface {
	AccessRedisDB

	SetRedisConn()
}
