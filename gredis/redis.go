package gredis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// RedisPool Redis连接结构
type RedisPool struct {
	p        *redis.Pool // redis connection pool
	connInfo string
	dbNum    int
	password string
}

func newRedis() *RedisPool {
	newRedis := new(RedisPool)

	return newRedis
}

func (rc *RedisPool) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// ZAdd redis有序集合添加.
func (rc *RedisPool) ZAdd(keyName string, updateData string) error {

	cur := time.Now()
	timestamp := cur.UnixNano()
	var err error
	if _, err = rc.do("ZADD", keyName, timestamp, updateData); err != nil {
		return err
	}
	return err
}

// ZRange redis有序集合获取
func (rc *RedisPool) ZRange(keyName string) ([]string, error) {

	return redis.Strings(rc.do("ZRANGE", keyName, 0, -1, "WITHSCORES"))
}

// ZScore redis有序集合获取数值.
func (rc *RedisPool) ZScore(keyName string, updateData string) (string, error) {

	return redis.String(rc.do("ZSCORE", keyName, updateData))
}

// ZRem redis有序集合删除
func (rc *RedisPool) ZRem(keyName string, updateData string) error {
	var err error
	if _, err = rc.do("ZREM", keyName, updateData); err != nil {
		return err
	}
	return err
}

// Get redis获取value string.
func (rc *RedisPool) Get(key string) (string, error) {

	return redis.String(rc.do("GET", key))
}

// Keys redis获取符合条件的key.
func (rc *RedisPool) Keys(key string) ([]string, error) {
	keys := key + "*"
	return redis.Strings(rc.do("KEYS", keys))
}

// GetInt64 redis获取value int64.
func (rc *RedisPool) GetInt64(key string) (int64, error) {

	return redis.Int64(rc.do("GET", key))
}

// GetMulti redis获取多组value.
func (rc *RedisPool) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var rv []interface{}
	c := rc.p.Get()
	defer c.Close()
	var err error
	for _, key := range keys {
		err = c.Send("GET", key)
		if err != nil {
			goto ERROR
		}
	}
	if err = c.Flush(); err != nil {
		goto ERROR
	}
	for i := 0; i < size; i++ {
		if v, err := c.Receive(); err == nil {
			rv = append(rv, v.([]byte))
		} else {
			rv = append(rv, err)
		}
	}
	return rv
ERROR:
	rv = rv[0:0]
	for i := 0; i < size; i++ {
		rv = append(rv, nil)
	}

	return rv
}

// Set redis设置value
func (rc *RedisPool) Set(key string, val interface{}) error {
	var err error
	if _, err = rc.do("SET", key, val); err != nil {
		return err
	}
	return err
}

// Delete redis删除key
func (rc *RedisPool) Delete(key string) error {
	var err error
	if _, err = rc.do("DEL", key); err != nil {
		return err
	}

	return err
}

// IsExist 判断是否存在key
func (rc *RedisPool) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr 原子操作添加数值
func (rc *RedisPool) Incr(key string) int64 {
	value, _ := redis.Int64(rc.do("INCRBY", key, int64(1)))
	return value
}

// Decr 原子操作减去数值
func (rc *RedisPool) Decr(key string) int64 {

	value, _ := redis.Int64(rc.do("INCRBY", key, int64(-1)))
	return value
}

func (rc *RedisPool) start(config *RedisConfig) error {

	rc.connInfo = config.Path
	rc.dbNum = 0
	rc.password = config.Password

	rc.connectInit(config)

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

func (rc *RedisPool) connectInit(config *RedisConfig) {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.connInfo)
		if err != nil {
			return nil, err
		}
		if rc.password != "" {
			if _, passerr := c.Do("AUTH", rc.password); passerr != nil {
				c.Close()
				return nil, passerr
			}
		}
		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     config.MaxIdleConns,
		MaxActive:   config.MaxOpenConns,
		IdleTimeout: 300 * time.Second,
		Dial:        dialFunc,
	}
}

//NewRedis redis创建
func NewRedis(config *RedisConfig) (adapter *RedisPool, err error) {
	adapter = newRedis()
	err = adapter.start(config)
	if err != nil {
		adapter = nil
	}
	return
}
