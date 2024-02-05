package core

import (
	"context"
	"database/sql"
	"snowflake/comm"
	"snowflake/config"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

// 获取缓存连接
func OpenRedis(ctx context.Context, module string) *redis.Client {

	endpoint := config.GetViper().GetString("redis.endpoint")
	if len(endpoint) <= 0 {
		comm.Error(module, "failed to get redis config")
		return nil
	}

	options, err := redis.ParseURL(endpoint)
	if err != nil {
		comm.Error(module, "failed to get redis options: %s", err.Error())
		return nil
	}

	options.PoolSize = 1000
	if !strings.Contains(endpoint, "min_idle_conns") {
		options.MinIdleConns = 8
	}

	client := redis.NewClient(options)
	if _, err := client.Ping(ctx).Result(); err != nil {
		comm.Error(module, "failed to connect redis: %s", err.Error())
		return nil
	}

	return client
}

// 获取数据库连接
func OpenMysql(ctx context.Context, module string) *sql.DB {

	endpoint := config.GetViper().GetString("mysql.endpoint")
	if len(endpoint) <= 0 {
		comm.Error(module, "failed to get mysql config")
		return nil
	}

	db, err := sql.Open("mysql", endpoint)
	if err != nil {
		comm.Error(module, "failed to connect mysql: %s -> %s", endpoint, err.Error())
		return nil
	}
	return db
}
