package prelude

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

// Stuff in this file is relevant for the entire application so it is imported
// using . imports.

const ttl = 24 * time.Hour

var RedisCtx = context.Background()

type UserId = string
type UserLogin = string

type Command struct {
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

// FetchJson runs an http request and unmarshals the response into v.
func FetchJson(url string, v interface{}) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}

// GetRedisKey4 presets a uniform way to builds a redis key in the application to
// avoid programming mistakes. Use GetRedisKey3 if you don't have a provider.
func GetRedisKey4(platform string, uid string, category string, provider string) string {
	return platform + "." + uid + "." + category + "." + provider
}

// GetRedisKey3 presets a uniform way to builds a redis key in the application to
// avoid programming mistakes. Use GetRedisKey4 if you do have a provider.
func GetRedisKey3(platform string, uid string, category string) string {
	return platform + "." + uid + "." + category
}

var errCooldown = errors.New("cooldown not ready")

// RedisCooldown returns nil if the cooldown is ready and an error otherwise.
func RedisCooldown(rdb *redis.Client, key string) error {
	exists, err := rdb.Exists(RedisCtx, key).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		return errCooldown
	}

	_, err = rdb.Set(RedisCtx, key, "", ttl).Result()
	if err != nil {
		return err
	}

	return nil
}
