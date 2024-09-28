package main

import (
	"context"
	"fmt"

	"github.com/jasurxaydarov/book_shop_api_getway/api"
	"github.com/jasurxaydarov/book_shop_api_getway/config"
	"github.com/jasurxaydarov/book_shop_api_getway/pkg/db"
	"github.com/jasurxaydarov/book_shop_api_getway/redis"
	"github.com/jasurxaydarov/book_shop_api_getway/service"
	"github.com/saidamir98/udevs_pkg/logger"
)


func main(){
	cfg := config.Load()

	log:=logger.NewLogger("",logger.LevelDebug)
	service:=service.Service()

	fmt.Println(service)

	redisCli, err := db.ConnRedis(log, context.Background(), cfg.RedisConfig)
	
	if err != nil {

		return
	}

	fmt.Println(redisCli)
	
	cache := redis.NewRedisRepo(redisCli, log)

	engine:=api.Api(api.Options{Service: service,Log:log,Cache: cache})

	engine.Run(":8080")
}