package redis

import (
	"fmt"
	"github.com/astaxie/beego"
	redigo "github.com/gomodule/redigo/redis"
	"strconv"
)

var redisPool *redigo.Pool

func init() {
	var (
		con redigo.Conn
		err error
	)
	host := beego.AppConfig.String("redis::host")
	port, _ := beego.AppConfig.Int("redis::port")
	password := beego.AppConfig.String("redis::password")
	maxIdle, _ := beego.AppConfig.Int("redis::max_idle")
	maxActive, _ := beego.AppConfig.Int("redis::max_active")
	db := 0
	redisPool = &redigo.Pool{
		MaxIdle:   maxIdle,   //最大空闲数
		MaxActive: maxActive, // 最大连接数
		Wait:      true,
		Dial: func() (redigo.Conn, error) {
			if password != "" {
				con, err = redigo.Dial("tcp", host+":"+strconv.Itoa(port),
					redigo.DialPassword(password),
					redigo.DialDatabase(db))
			} else {
				con, err = redigo.Dial("tcp", host+":"+strconv.Itoa(port),
					redigo.DialDatabase(db))
			}
			//redis.DialConnectTimeout(timeout*time.Second),
			//redis.DialReadTimeout(timeout*time.Second),
			//redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, err
		},
	}
}

//获取连接
func GetCon(db uint) redigo.Conn {
	con := redisPool.Get()
	con.Do("SELECT", db)
	return con
}

// Key是否存在
func KeyExists(db uint, key string) bool {
	var (
		err      error
		flag     int
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return flag > 0
}

//匹配的keys列表
func KeysList(db uint, pattern string) []string {
	var (
		err      error
		list     []string
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	list, err = redigo.Strings(redisCon.Do("Keys", pattern))
	if err != nil {
		return nil
	}
	return list
}

//设置Key
func SetString(db uint, key string, value interface{}, expireSeconds int) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	if expireSeconds > 0 {
		_, err = redisCon.Do("SET", key, value, "EX", expireSeconds)
	} else {
		_, err = redisCon.Do("SET", key, value)
	}
	return err == nil
}

//获取Key
func GetString(db uint, key string) (string, error) {
	var (
		value    interface{}
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	//value, err = redisCon.Do("GET", key)
	value, err = redigo.String(redisCon.Do("GET", key))
	return value.(string), err
}

//设置过期时间
func ExpireKey(db uint, key string, expireSeconds int) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	_, err = redigo.Int(redisCon.Do("expire", key, expireSeconds))
	if err != nil {
		return false
	}
	return true
}

//获取Key剩余有效秒数
func GetTTL(db uint, key string) int {
	var (
		value    int
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	value, err = redigo.Int(redisCon.Do("TTL", key))
	if err != nil || value == -2 {
		return 0
	}
	return value
}

//删除Key
func DelKey(db uint, key string) bool {
	var (
		err  error
		flag int
	)
	redisCon := GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("DEL", key))
	if err != nil {
		return false
	}
	return flag > 0
}

//Hash Key是否存在
func HashExists(db uint, key string) bool {
	var (
		err      error
		flag     int
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("HEXISTS", key))
	if err != nil {
		return false
	}
	return flag > 0
}

//删除Key
func DelHash(db uint, key string, fields ...string) bool {
	var (
		err      error
		flag     int
		redisCon redigo.Conn
		args     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	args = []interface{}{key}
	for _, v := range fields {
		args = append(args, v)
	}
	flag, err = redigo.Int(redisCon.Do("HDEL", args...))
	if err != nil {
		return false
	}
	return flag > 0
}

//设置Hash
func SetHash(db uint, key string, kvs ...interface{}) bool {
	var (
		err error
		//flag     interface{}
		redisCon redigo.Conn
		args     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	args = []interface{}{key}
	for _, v := range kvs {
		args = append(args, v)
	}
	_, err = redisCon.Do("HMSET", args...)
	return err == nil
}

//获取Hash
func GetHash(db uint, key string, fields ...string) []interface{} {
	var (
		err      error
		data     []interface{}
		redisCon redigo.Conn
		args     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	args = []interface{}{key}
	for _, v := range fields {
		args = append(args, v)
	}
	data, err = redigo.Values(redisCon.Do("HMGET", args...))
	if err != nil {
		return nil
	}
	return data
}

//设置List
func SetList(db uint, key string, value string, pre bool) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	if pre {
		_, err = redisCon.Do("LPUSH", key, value)
	} else {
		_, err = redisCon.Do("RPUSH", key, value)
	}
	return err == nil
}

//获取List
func GetList(db uint, key string, begin int64, end int64) []interface{} {
	var (
		err      error
		redisCon redigo.Conn
		data     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	data, err = redigo.Values(redisCon.Do("LRANGE", key, begin, end))
	if err != nil {
		return nil
	}
	return data
}

//删除list所有value值的项
func DelList(db uint, key string, value string) bool {
	var (
		err      error
		redisCon redigo.Conn
		flag     int
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("LREM", key, 0, value))
	if err != nil {
		return false
	}
	return flag > 0
}

//设置有序集合
func SetSortSet(db uint, key string, score int64, value string) bool {
	var (
		err      error
		redisCon redigo.Conn
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	_, err = redisCon.Do("ZADD", key, score, value)
	if err != nil {
		fmt.Println("redis hash:" + err.Error())
	}
	return err == nil
}

//获取有序集合
func GetSortSet(db uint, key string, begin int64, end int64, withScore bool) []interface{} {
	var (
		err      error
		redisCon redigo.Conn
		data     []interface{}
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	if withScore {
		data, err = redigo.Values(redisCon.Do("ZRANGE", key, begin, end, "WITHSCORES"))
	} else {
		data, err = redigo.Values(redisCon.Do("ZRANGE", key, begin, end))
	}
	if err != nil {
		return nil
	}
	return data
}

// 获取有序集合总数
func GetZCard(db uint, key string) int64 {
	var (
		err      error
		redisCon redigo.Conn
		flag     int64
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int64(redisCon.Do("Zcard", key))
	if err != nil {
		return 0
	}
	return flag
}

//删除SortSet所有value值的项
func DelSortSet(db uint, key string, value string) bool {
	var (
		err      error
		redisCon redigo.Conn
		flag     int
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	flag, err = redigo.Int(redisCon.Do("ZREM", key, value))
	if err != nil {
		return false
	}
	return flag > 0
}

//限制访问次数
func LimitVisit(key string, LimitCount int64) bool {
	value, err := GetString(0, key)
	if err != nil {
		return false
	}
	currentCount, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}
	if currentCount < LimitCount {
		return false
	}

	return true
}

//设置key过期时间，如果还有效不改变原有过期时间
//判断是否存在该Key
//不存在：设置K-V 超时时间
//存在：更新K-V 保留原有的超时时间
//设值
func SetKeyRemainExpire(db uint, key string, value interface{}, expireSeconds int) error {
	var (
		err      error
		redisCon redigo.Conn
		ttl      int64
	)
	redisCon = GetCon(db)
	defer redisCon.Close()
	ttl, err = redigo.Int64(redisCon.Do("TTL", key))
	if err != nil || ttl <= 0 {
		_, err = redisCon.Do("SET", key, value, "EX", expireSeconds)
	} else {
		value, err = redigo.Int64(redisCon.Do("SET", key, value, "EX", ttl))
	}
	if err != nil {
		return err
	}
	return nil
}

//初始过期时间内限制次数增加
func LimitCountIncrease(db uint, key string, expireSeconds int) bool {
	var (
		err          error
		currentCount int64
		value        interface{}
	)
	currentCount = 0
	value, err = GetString(db, key)
	if err == nil && value != nil {
		valueStr := value.(string)
		currentCount, err = strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return false
		}
	}
	err = SetKeyRemainExpire(db, key, currentCount+1, expireSeconds)
	if err != nil {
		return false
	}

	return true
}
