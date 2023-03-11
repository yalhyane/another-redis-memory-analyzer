package utils

import (
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

var redisClient *redis.Client
var redisConfig RedisConfig

type RedisConfig struct {
	host     string
	port     int
	password string
}

func InitConfig(host string, port int, password string) {
	redisConfig = RedisConfig{host, port, password}
}

func ConnectionToDb(db int) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:        redisConfig.host + ":" + strconv.Itoa(redisConfig.port),
		Password:    redisConfig.password,
		DB:          db,
		ReadTimeout: 1 * time.Minute,
	})

	pong, err := client.Ping().Result()

	if err != nil || pong == "" {
		log.Fatal("Cannot connect to redis", err)
	}

	return client
}

func Connection(host string, port int, password string, db int) *redis.Client {

	InitConfig(host, port, password)

	client := redis.NewClient(&redis.Options{
		Addr:        host + ":" + strconv.Itoa(port),
		Password:    password, // no password set
		DB:          db,       // use default DB
		ReadTimeout: 1 * time.Minute,
	})

	pong, err := client.Ping().Result()

	if err != nil || pong == "" {
		log.Fatal("Cannot connect to redis", err)
	}

	return client
}

func GetDatabases(client *redis.Client) map[uint64]int64 {

	var databases = make(map[uint64]int64)
	reply, _ := client.Info("keyspace").Result()
	keyspace := strings.Trim(reply[12:], "\n")
	keyspaces := strings.Split(keyspace, "\r")
	//log.Println(keyspace)
	for _, db := range keyspaces {
		dbKeysParse := strings.Split(db, ",")
		//fmt.Println(dbKeysParse[0])
		if dbKeysParse[0] == "" {
			continue
		}

		dbKeysParsed := strings.Split(dbKeysParse[0], ":")
		//fmt.Println("DB", strings.Trim(dbKeysParsed[0], "\n")[2:])
		dbNo, _ := strconv.ParseUint(strings.Trim(dbKeysParsed[0], "\n")[2:], 10, 64)
		// fmt.Println("DB", dbNo)
		dbKeySize, _ := strconv.ParseInt(dbKeysParsed[1][5:], 10, 64)
		if dbKeySize == 0 {
			continue
		}
		databases[dbNo] = dbKeySize
	}
	// log.Println(databases)

	return databases
}

func Close(client *redis.Client) {
	if client != nil {
		_ = client.Close()
	}
}
