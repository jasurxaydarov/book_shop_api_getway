package db

import (
	"context"

	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jasurxaydarov/book_shop_api_getway/config"
	"github.com/saidamir98/udevs_pkg/logger"
)

func RedisAdr(host string, port int) string {

	return host + ":" + strconv.Itoa(port)
}

func ConnRedis(log logger.LoggerI, ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {

	redisCli := redis.NewClient(&redis.Options{
		Addr: RedisAdr(cfg.Host, cfg.Port),
		DB:   cfg.DBIndex,
	})

	_, err := redisCli.Ping(ctx).Result()
	if err != nil {

		log.Error("err on Con to Redis", logger.Error(err))
		return nil, err
	}

	log.Debug("sucesfully conected with redis")
	return redisCli, nil
}
