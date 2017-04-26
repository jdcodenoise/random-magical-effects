package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// Opens a connection to redis, and loads magical effects from TSV into Redis
func loadEffects() error {
	conn := pool.Get()
	defer conn.Close()

	file, err := os.Open("random_magical_effects.tsv")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), "\t")

		exists, err := redis.Int(conn.Do("SISMEMBER", "effects", data[0]))
		if err != nil {
			return err
		}

		if exists != 1 {
			log.Printf("adding effect to redis %v: %v", data[0], data[1])
			if _, err := conn.Do("SADD", "effects", data[0]); err != nil {
				return err
			}
			if _, err := conn.Do("SET", data[0], data[1]); err != nil {
				return err
			}
		}
	}
	return nil
}

func randomEffect() (Effect, error) {
	conn := pool.Get()
	defer conn.Close()

	var e Effect

	effect_key, err := redis.String(conn.Do("SRANDMEMBER", "effects"))
	if err != nil {
		return e, err
	}

	effect, err := redis.String(conn.Do("GET", effect_key))
	if err != nil {
		return e, err
	}

	e = Effect{Key: effect_key, Text: effect}

	return e, nil
}
