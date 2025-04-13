package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gomodule/redigo/redis"
)

func connectToRedis(password string) (redis.Conn, error) {

	var opts = make([]redis.DialOption, 0)
	if password != "" {
		opts = append(opts, redis.DialPassword(password))
	}

	conn, err := redis.Dial("tcp", "redis:6379",
		opts...,
	)

	if err != nil {
		// Wrap the error instead of loggin in a helper function
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	_, err = redis.String(conn.Do("PING"))
	if err != nil {
		return nil, fmt.Errorf("error pinging Redis: %w", err)
	}

	return conn, nil
}

func main() {
	password, hasPassword := os.LookupEnv("REDIS_PASSWORD")
	if !hasPassword {
		slog.Warn("No password provided for Redis connection")
	}

	conn, err := connectToRedis(password)
	if err != nil {
		slog.Error("Error connecting to Redis", "error", err)
		return
	}
	defer conn.Close()
	slog.Info("Connected to Redis successfully")

}
