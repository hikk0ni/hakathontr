package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func SaveSession(RfTokenUUID string, RtExists int64) {
	err := RDB.Set(ctx, RfTokenUUID, "", time.Duration(RtExists)).Err()
	if err != nil {
		panic(err)
	}
}

func CheckSession(RfTokenUUID string) error {
	count, err := RDB.Exists(ctx, RfTokenUUID).Result()
	if count > 0 {
		return nil
	} else {
		return err
	}
}

func DeleteSession(RfTokenUUID string) {
	err := RDB.Del(ctx, RfTokenUUID)
	if err != nil {
		panic(err)
	}
}
