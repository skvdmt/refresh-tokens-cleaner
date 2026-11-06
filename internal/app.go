package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/skvdmt/refresh-tokens-cleaner/internal/model"
)

const (
	REDIS_PASSWORD_ENV = "REDIS_PASSWORD"
)

// RemoveExpiredRefreshTokens Удаление просроченых токенов обновления.
func RemoveExpiredRefreshTokens() error {
	pwd, ok := os.LookupEnv(REDIS_PASSWORD_ENV)
	if !ok {
		return fmt.Errorf("env %s not set", REDIS_PASSWORD_ENV)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			model.Conf.Redis.Host, model.Conf.Redis.Port),
		Password: pwd,
		DB:       model.Conf.Redis.Db,
	})
	ctx, cance := context.WithTimeout(context.Background(), 10*time.Second)
	defer cance()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	var cnt uint64
	var usrs int
	var tdel int64
	exp := time.Now().Unix() - int64((time.Duration(model.Conf.TokenTTLMinutes) * time.Minute).Seconds())
	for {
		var keys []string
		var err error
		keys, cnt, err = rdb.Scan(ctx, cnt, "*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			pipe := rdb.Pipeline()
			cmds := make([]*redis.IntCmd, len(keys))
			for i, u := range keys {
				cmds[i] = pipe.ZRemRangeByScore(ctx, u, "-inf", fmt.Sprintf("%d", exp))
			}
			if _, err := pipe.Exec(ctx); err != nil {
				model.Log.Error.Error(fmt.Sprintf("cleaning error: %v", err))
			}
			var d int64
			for _, c := range cmds {
				if r, err := c.Result(); err == nil {
					d += r
				}
			}
			tdel += d
			usrs += len(keys)
		}
		if cnt == 0 {
			break
		}
	}
	model.Log.Info.Info(fmt.Sprintf("users cleaned: %d", usrs))
	model.Log.Info.Info(fmt.Sprintf("token removed: %d", tdel))
	return nil
}
