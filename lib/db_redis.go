package lib

/*
http://godoc.org/github.com/garyburd/redigo/redis

封装常用的Redis命令
REDIS pool的使用
*/

import (
	"errors"
	"fmt"
	"log"
	"time"
)

import (
	"github.com/garyburd/redigo/redis"

	. "github.com/ubrabbit/go-server/common"
)

var (
	g_RedisDB   redis.Conn  = nil
	g_RedisPool *redis.Pool = nil
)

const (
	MAX_REDIS_POOL_ACTIVE = 3000
)

type PoolDial struct {
	Conn redis.Conn
	Err  error
}

func newRedisPool(host string) *redis.Pool {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			ch := make(chan PoolDial)
			defer close(ch)
			go func() {
				wait := 0
				for {
					if wait >= 60 {
						err := errors.New("fatal: redis pool want to connect db, but wait too long")
						ch <- PoolDial{nil, err}
						break
					}
					//在短期高并发导致端口用尽时，会报 cannot assign requested address 错误
					//所以需要用chan等待连接释放
					c, err := redis.Dial("tcp", host)
					if err != nil {
						fmt.Println(err)
						time.Sleep(1 * time.Second)
						wait++
						continue
					}
					ch <- PoolDial{c, err}
					break
				}
			}()

			rlt := <-ch
			return rlt.Conn, rlt.Err
		},
		MaxIdle:     10,
		MaxActive:   MAX_REDIS_POOL_ACTIVE, // 最大连接数量，如果不设置这个值默认就是无限，当短时间高并发时报：too many open files
		Wait:        true,                  // 当达到最大连接数量时，阻塞， 如果不加这个参数，会报: connection pool exhausted
		IdleTimeout: 360 * time.Second,
		//获取连接对象前检查下连接是否还活着
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}

func checkRedisConn(conn redis.Conn) {
	if conn == nil {
		log.Fatalf("DB RedisConn %v is not inited!!!!!", conn)
	}
}

func InitRedis(host string, port int) redis.Conn {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := redis.Dial("tcp", address)
	CheckFatal(err)

	fmt.Printf("Connect Redis %s Succ\n", address)
	g_RedisDB = conn
	g_RedisPool = newRedisPool(address)
	return conn
}

func CloseRedis() {
	if g_RedisPool != nil {
		g_RedisPool.Close()
	}
	if g_RedisDB != nil {
		g_RedisDB.Close()
	}
}

func GetRedisConn() redis.Conn {
	conn := g_RedisPool.Get()
	if conn != nil {
		return conn
	}
	checkRedisConn(g_RedisDB)
	return g_RedisDB
}

func RedisConnAlive(conn redis.Conn) bool {
	_, err := conn.Do("PING")
	if err != nil {
		return false
	}
	return true
}

func RedisConnExec(conn redis.Conn, cmd string, arg ...interface{}) (int, error) {
	result, err := redis.Int(conn.Do(cmd, arg...))
	if err != nil {
		fmt.Println("Redis Exec error: ", cmd, arg, result, err)
		return result, err
	}

	//fmt.Println("result is ", result)
	return result, nil
}

func RedisExec(cmd string, arg ...interface{}) (int, error) {
	conn := GetRedisConn()
	//不加这行语句会导致死锁
	//比如同一个函数执行了两次 RedisExec，但获取的是不同的conn的情况
	defer conn.Close()

	return RedisConnExec(conn, cmd, arg...)
}

func RedisConnGetString(conn redis.Conn, cmd string, arg ...interface{}) interface{} {
	value, err := conn.Do(cmd, arg...)
	CheckFatal(err)
	if value == nil {
		return nil
	}
	value, err = redis.String(value, err)
	if err != nil {
		fmt.Println("RedisGetString error: ", cmd, arg, err)
		return nil
	}
	//fmt.Println("value is ", value)
	return value
}

func RedisGetString(cmd string, arg ...interface{}) interface{} {
	conn := GetRedisConn()
	defer conn.Close()

	return RedisConnGetString(conn, cmd, arg...)
}

func RedisConnGetInt(conn redis.Conn, cmd string, arg ...interface{}) interface{} {
	value, err := conn.Do(cmd, arg...)
	CheckFatal(err)
	if value == nil {
		return nil
	}
	value, err = redis.Int64(value, err)
	if err != nil {
		fmt.Println("RedisGetInt error: ", cmd, arg, err)
		return nil
	}
	//fmt.Println("value is ", value)
	return value
}

func RedisGetInt(cmd string, arg ...interface{}) interface{} {
	conn := GetRedisConn()
	defer conn.Close()

	return RedisConnGetInt(conn, cmd, arg...)
}

func RedisConnGetList(conn redis.Conn, cmd string, arg ...interface{}) []string {
	value_list, err := redis.Values(conn.Do(cmd, arg...))
	if err != nil {
		fmt.Println("RedisGetList error: ", cmd, arg, err)
		return nil
	}
	result := make([]string, 0)
	for _, value := range value_list {
		result = append(result, string(value.([]byte)))
	}
	fmt.Println("result is ", result)
	return result
}

func RedisGetList(cmd string, arg ...interface{}) []string {
	conn := GetRedisConn()
	defer conn.Close()

	return RedisConnGetList(conn, cmd, arg...)
}

func RedisConnGetMap(conn redis.Conn, cmd string, arg ...interface{}) map[string]string {
	value, err := redis.StringMap(conn.Do(cmd, arg...))
	if err != nil {
		fmt.Println("RedisGetMap error: ", cmd, arg, err)
		return nil
	}
	//fmt.Println("value is ", value)
	return value
}

func RedisGetMap(cmd string, arg ...interface{}) map[string]string {
	conn := GetRedisConn()
	defer conn.Close()

	return RedisConnGetMap(conn, cmd, arg...)
}
