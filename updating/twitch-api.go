package updating

import (
	"errors"
	"log"
	"os"
	. "self/prelude"
	"time"

	"github.com/nicklaw5/helix"

	"github.com/go-redis/redis/v8"
)

var (
	helixClient         *helix.Client
	errAppTokenNotReady = errors.New("app access token not ready")
)

func init() {
	if os.Getenv("TWITCH_CLIENT_ID") == "" {
		log.Println("warning: twitch client id not set")
	}
	if os.Getenv("TWITCH_CLIENT_SECRET") == "" {
		log.Println("warning: twitch client secret not set")
	}

	c, err := helix.NewClient(&helix.Options{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
	})

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			resp, err := c.RequestAppAccessToken([]string{})
			if err != nil {
				log.Println("couldn't get app access token")

				time.Sleep(time.Minute * 10)
				continue
			}

			helixClient.SetAppAccessToken(resp.Data.AccessToken)
			time.Sleep(time.Hour * 24)
		}
	}()

	helixClient = c
}

// getLogin returns the login of the user with the given user id.
func getLogin(rdb *redis.Client, uid UserId) (string, error) {
	key := GetRedisKey3("twitch", uid, "login")

	// cache
	login, err := rdb.Get(RedisCtx, key).Result()
	if err == nil {
		return login, nil
	}
	if err != redis.Nil {
		return "", nil
	}

	// get from twitch
	if helixClient.GetAppAccessToken() == "" {
		return "", errAppTokenNotReady
	}

	resp, err := helixClient.GetUsers(&helix.UsersParams{
		IDs: []string{uid},
	})

	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}
	if len(resp.Data.Users) == 0 {
		return "", errors.New("no user in list found")
	}

	rdb.Set(RedisCtx, key, resp.Data.Users[0].Login, time.Hour*100)

	return resp.Data.Users[0].Login, nil
}
